# Design: Goldmark Parser Integration

## Context

The current markdown parsing implementation uses regex-based pattern matching (`regexp.MustCompile`, `FindStringSubmatch`, etc.) to extract structure from specification files. While this works for simple cases, it has limitations:

1. **Fragility**: Regex patterns are brittle and fail on edge cases (code blocks containing `##`, escaped characters, nested structures)
2. **Maintainability**: Complex regex patterns are hard to read and modify
3. **Correctness**: Not CommonMark compliant, may misparse valid markdown
4. **Limited Context**: Regex doesn't understand markdown structure (can't distinguish between headers in code blocks vs real headers)

The previous design decision (archived in 2025-11-17-add-validate-command) chose regex for simplicity, stating "custom regex-based extraction is tractable" and rejecting goldmark as "overkill." However, as the parsing logic has grown (3 parser files, 13 consuming files), the complexity justifies a proper markdown parser.

### Constraints
- Must maintain backward-compatible API for existing consumers
- Must not break existing tests (or update them appropriately)
- Must handle all current parsing scenarios (titles, sections, requirements, scenarios, deltas, tasks)
- Must integrate with existing error handling patterns
- Binary size increase should be justified by benefits

### Stakeholders
- Developers using Spectr (benefit from more robust parsing)
- CI/CD systems (benefit from better error messages with locations)
- Future contributors (benefit from cleaner, more maintainable code)

## Goals / Non-Goals

### Goals
- Replace all regex-based markdown parsing with goldmark AST traversal
- Maintain 100% API compatibility for existing functions
- Improve error messages with AST node location information
- Add test coverage for edge cases that regex parsing missed
- Reduce code complexity in parser modules
- Provide foundation for future extensions (frontmatter, custom blocks)

### Non-Goals
- Changing the spec format or conventions
- Auto-migration of existing specs
- Adding new parsing capabilities beyond current functionality
- Rendering markdown to HTML (goldmark supports this, but we don't need it)
- Supporting non-CommonMark markdown dialects

## Decisions

### Decision 1: Use Goldmark for All Markdown Parsing
**Choice**: Replace all regex-based parsing with goldmark AST traversal

**Rationale**:
- **CommonMark Compliance**: Goldmark is fully compliant with CommonMark spec, handles edge cases correctly
- **AST-Based**: Proper tree structure makes it easy to distinguish headers in code blocks vs real headers
- **Performance**: Goldmark performance is comparable to cmark (C reference implementation)
- **Extensibility**: Interface-based AST allows future extensions without regex changes
- **Memory Efficiency**: Uses text segments instead of copying strings
- **Battle-Tested**: Widely used in Go ecosystem (Hugo, Gitea, etc.)

**Implementation approach**:
1. Parse markdown file with goldmark parser
2. Walk AST to find heading nodes at specific levels
3. Extract text content from AST nodes
4. Maintain same function signatures for backward compatibility

**Alternatives considered**:
- Continue with regex: Rejected due to growing complexity and fragility
- Other parsers (gomarkdown, blackfriday): Rejected due to inferior extensibility or maintenance

### Decision 2: Maintain Backward-Compatible API
**Choice**: Keep all existing function signatures, only change internals

**Rationale**:
- **Zero Breaking Changes**: No need to update 13 consuming files
- **Incremental Migration**: Can verify correctness by comparing outputs
- **Testing Strategy**: Existing tests should pass with new implementation
- **Risk Reduction**: Limits scope of change to parser internals

**API to maintain**:
```go
// parsers/parsers.go
func ExtractTitle(filePath string) (string, error)
func CountTasks(filePath string) (TaskStatus, error)
func CountDeltas(changeDir string) (int, error)
func CountRequirements(specPath string) (int, error)

// parsers/requirement_parser.go
func ParseRequirements(filePath string) ([]RequirementBlock, error)
func ParseScenarios(requirementContent string) []string
func NormalizeRequirementName(name string) string

// parsers/delta_parser.go
func ParseDeltaSpec(filePath string) (*DeltaPlan, error)

// validation/parser.go
func ExtractSections(content string) map[string]string
func ExtractRequirements(content string) []Requirement
func ExtractScenarios(requirementBlock string) []string
```

**Alternatives considered**:
- New API with breaking changes: Rejected to limit scope and risk
- Separate goldmark package: Rejected as unnecessary abstraction

### Decision 3: AST Walking Strategy
**Choice**: Use goldmark's ast.Walker with custom NodeVisitor for each parsing function

**Rationale**:
- **Idiomatic Goldmark**: Walker pattern is the recommended approach
- **Stateful Extraction**: NodeVisitor can maintain state during traversal
- **Type Safety**: Can handle different node types (Heading, List, etc.) with type assertions
- **Performance**: Single-pass AST walk is efficient

**Implementation pattern**:
```go
type headingExtractor struct {
    level int      // Target heading level (1 for #, 2 for ##, etc.)
    text  []string // Collected heading text
}

func (h *headingExtractor) Visit(node ast.Node, entering bool) (ast.WalkStatus, error) {
    if !entering {
        return ast.WalkContinue, nil
    }

    if heading, ok := node.(*ast.Heading); ok && heading.Level == h.level {
        // Extract text from heading node
        text := extractText(heading, source)
        h.text = append(h.text, text)
    }

    return ast.WalkContinue, nil
}
```

**Alternatives considered**:
- Manual tree traversal: More code, less idiomatic
- Recursive descent: Harder to manage state across levels

### Decision 4: Text Extraction from AST Nodes
**Choice**: Use ast.Node.Text() with source byte slice to extract content

**Rationale**:
- **Goldmark Design**: Nodes store segments, not text directly
- **Memory Efficient**: No string copying during parsing
- **Accurate**: Preserves original source text exactly

**Implementation**:
```go
func extractText(node ast.Node, source []byte) string {
    var buf bytes.Buffer
    for child := node.FirstChild(); child != nil; child = child.NextSibling() {
        segment := child.Text(source)
        buf.Write(segment)
    }
    return buf.String()
}
```

**Alternatives considered**:
- Render to string: Overkill, would include markdown formatting
- Manual segment concatenation: Error-prone, reinventing wheel

### Decision 5: Error Handling and Location Information
**Choice**: Enhance error messages with line numbers from AST nodes

**Rationale**:
- **Better UX**: Users can quickly find problematic sections
- **Goldmark Capability**: AST nodes have segment information with line/column
- **Validation Enhancement**: Validation errors can point to exact locations

**Implementation**:
```go
type ValidationIssue struct {
    Level    ValidationLevel
    Path     string
    Line     int    // NEW: Line number from AST
    Column   int    // NEW: Column number from AST
    Message  string
}
```

**Alternatives considered**:
- No location info: Missed opportunity to improve error messages
- Full source context: Too verbose, harder to implement

### Decision 6: Testing Strategy
**Choice**: Keep all existing tests, add new edge case tests

**Rationale**:
- **Regression Prevention**: Existing tests verify backward compatibility
- **Coverage Improvement**: Add tests for cases that regex couldn't handle
- **Confidence**: Comprehensive test suite ensures correctness

**New test cases**:
- Headers in code blocks (should be ignored)
- Headers in blockquotes
- Escaped hash characters
- Headers with inline code
- Unicode in headers and content
- Deeply nested structures
- Malformed markdown (missing closing backticks, etc.)

**Test approach**:
1. Run existing tests with new implementation
2. Add new edge case tests
3. Compare outputs with regex implementation on real spec files
4. Benchmark parsing performance

### Decision 7: Dependency Management
**Choice**: Add `github.com/yuin/goldmark` to go.mod, no version constraints

**Rationale**:
- **Stability**: Goldmark is mature and stable (v1.7+ is current)
- **Zero Breaking Changes**: Goldmark maintains API compatibility
- **Size Justification**: ~200KB increase in binary is acceptable for benefits
- **Transitive Dependencies**: Goldmark has zero dependencies itself

**go.mod addition**:
```
require (
    github.com/yuin/goldmark v1.7.8
)
```

**Alternatives considered**:
- Vendor goldmark: Unnecessary complexity
- Pin specific version: No known stability issues to warrant this

## Risks / Trade-offs

### Risk: Subtle Behavior Changes
**Impact**: Medium
**Likelihood**: Medium

Goldmark may parse some edge cases differently than regex, potentially breaking specs that rely on regex quirks.

**Mitigation**:
1. Comprehensive testing with existing spec files
2. Side-by-side comparison tool to detect differences
3. Gradual rollout with feature flag if needed
4. Document any intentional behavior changes

### Risk: Performance Regression
**Impact**: Low
**Likelihood**: Low

Goldmark might be slower than regex for simple cases.

**Mitigation**:
1. Benchmark existing vs new implementation
2. Goldmark is highly optimized, unlikely to be slower
3. Can cache parsed AST if needed (future optimization)

### Risk: Binary Size Increase
**Impact**: Low
**Likelihood**: High

Goldmark will increase binary size by ~200KB.

**Trade-off Analysis**:
- **Cost**: 200KB increase (from ~8MB to ~8.2MB, 2.5% increase)
- **Benefit**: More robust parsing, better error messages, maintainable code
- **Verdict**: Acceptable trade-off for a CLI tool

### Risk: Migration Complexity
**Impact**: Medium
**Likelihood**: Low

Refactoring all parsing code could introduce subtle bugs.

**Mitigation**:
1. Incremental migration (one parser file at a time)
2. Keep regex implementation in git history for reference
3. Extensive testing before merge
4. Code review focused on correctness

## Migration Plan

### Phase 1: Foundation (parsers/parsers.go)
1. Add goldmark dependency to go.mod
2. Create internal helper functions for AST walking (extractText, findHeadings, etc.)
3. Refactor `ExtractTitle` to use goldmark
4. Refactor `CountTasks`, `CountDeltas`, `CountRequirements` to use goldmark
5. Run existing tests, add edge case tests
6. Verify no regressions

### Phase 2: Requirement Parsing (parsers/requirement_parser.go)
1. Refactor `ParseRequirements` to walk AST for `### Requirement:` headings
2. Refactor `ParseScenarios` to walk AST for `#### Scenario:` headings
3. Run existing tests, add edge case tests
4. Verify no regressions

### Phase 3: Delta Parsing (parsers/delta_parser.go)
1. Refactor `ParseDeltaSpec` to walk AST for delta sections
2. Update `parseDeltaSection`, `parseRemovedSection`, `parseRenamedSection`
3. Run existing tests, add edge case tests
4. Verify no regressions

### Phase 4: Validation Parsing (validation/parser.go)
1. Refactor `ExtractSections` to use goldmark
2. Refactor `ExtractRequirements` and `ExtractScenarios` to use goldmark
3. Run existing tests, add edge case tests
4. Verify no regressions

### Phase 5: Integration and Testing
1. Run full test suite across all packages
2. Test with real spec files in spectr/ directory
3. Benchmark performance comparison
4. Update documentation with goldmark usage patterns
5. Clean up any remaining regex patterns

### Phase 6: Enhancement (Optional)
1. Add line/column information to validation errors
2. Consider future extensions (frontmatter, custom blocks)
3. Document patterns for future parser additions

### Rollback Plan
If goldmark integration proves problematic:
1. Revert commits in phases (newest to oldest)
2. Each phase is atomic and can be reverted independently
3. Regex implementation remains in git history
4. No external API changes, so rollback is safe

## Open Questions

### Question 1: Should we cache parsed AST?
**Context**: If the same file is parsed multiple times in one command, caching could improve performance.

**Options**:
- **A**: No caching (simple, correct, minimal benefit for current usage)
- **B**: In-memory LRU cache with file mtime check (complex, potential benefit for bulk operations)
- **C**: Cache at call site (let callers decide if they want to cache)

**Recommendation**: Start with Option A. Profile during bulk validation to see if caching is needed.

### Question 2: Should we expose AST to consumers?
**Context**: Some consumers might benefit from direct AST access for custom operations.

**Options**:
- **A**: Keep AST internal, only expose parsed data structures (current API)
- **B**: Add optional `ParseAST(path string) (ast.Node, error)` function
- **C**: Change all APIs to return both data and AST

**Recommendation**: Option A for now. Can add Option B later if use cases emerge.

### Question 3: Should we add position information to all parsed structures?
**Context**: Goldmark provides line/column info, we could expose this.

**Options**:
- **A**: No position info (maintain current API)
- **B**: Add optional position fields to structs (backward compatible)
- **C**: Return separate position map

**Recommendation**: Option B for new error messages, Option A for backward compatibility.

### Question 4: Should we validate markdown correctness?
**Context**: Goldmark can detect malformed markdown (unclosed code blocks, etc.)

**Options**:
- **A**: Parse as-is, ignore markdown errors (lenient)
- **B**: Warn on markdown errors (helpful)
- **C**: Error on markdown errors (strict)

**Recommendation**: Option A initially (match current regex behavior), consider Option B in validation command.

## Success Criteria

1. All existing tests pass with goldmark implementation
2. No performance regression (parsing within 150% of regex performance)
3. Binary size increase under 5%
4. Zero API breaking changes
5. Improved error messages with location information
6. Code complexity reduced (fewer regex patterns, clearer logic)
7. New test coverage for edge cases (headers in code blocks, etc.)

## References

- [Goldmark GitHub](https://github.com/yuin/goldmark)
- [Goldmark Docs](https://pkg.go.dev/github.com/yuin/goldmark)
- [CommonMark Spec](https://spec.commonmark.org/)
- Previous design decision: `spectr/changes/archive/2025-11-17-add-validate-command/design.md`
