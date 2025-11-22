# Design: Goldmark Parser Integration

## Context

Spectr's current parsing layer uses regex patterns to extract markdown structure from spec files. This works but creates maintenance burden (28 patterns across 6 files) and is fragile for edge cases. We need a robust, maintainable solution that handles all CommonMark syntax correctly while preserving backward compatibility.

**Constraints:**
- Must maintain existing parser function signatures (backward compatible)
- Must pass all existing tests without behavior changes
- Must not degrade performance for typical spec files (<50KB)
- Must provide better error messages than regex approach

**Stakeholders:**
- Users: No visible changes; more robust parsing of edge cases
- Developers: Simpler codebase, easier to add new parsing features
- Maintainers: Less regex debugging, better error diagnostics

## Goals / Non-Goals

### Goals
- Replace all regex-based markdown parsing with goldmark AST traversal
- Maintain 100% backward compatibility with existing parser APIs
- Improve error messages with line/column information from AST
- Reduce code complexity by consolidating 28 regex patterns into AST walkers
- Achieve full CommonMark compliance for spec parsing

### Non-Goals
- Changing parser function signatures or return types
- Adding new markdown features (tables, footnotes) in this change
- Replacing validation logic (only the parsing layer)
- Supporting non-CommonMark markdown extensions initially

## Decisions

### Decision 1: Use Goldmark as Markdown Parser
**Choice:** `github.com/yuin/goldmark` v1.7.x

**Rationale:**
- De facto standard Go markdown parser (used by Hugo, Gitea, GitHub CLI)
- Full CommonMark 0.31.2 compliance
- Extensible architecture with parser/renderer/AST separation
- Active maintenance (5+ years, 400+ contributors)
- Excellent performance (benchmarked against other Go parsers)

**Alternatives Considered:**
1. **blackfriday** - Older, less maintained, not fully CommonMark compliant
2. **gomarkdown/markdown** - Less popular, smaller ecosystem
3. **Custom AST parser** - Too much implementation burden vs benefit

### Decision 2: Maintain Parser Function Signatures
**Choice:** Keep all public functions unchanged; only modify internal implementation

**Rationale:**
- Zero breaking changes for consumers (validation, archive, list packages)
- Gradual migration path: can test both implementations in parallel
- Easier to roll back if issues discovered

**Alternatives Considered:**
1. **New API with AST types** - Would break all consumers; too disruptive
2. **Parallel package** - Creates duplication and migration burden

### Decision 3: AST Traversal Strategy
**Choice:** Use goldmark's AST Walker for targeted extraction

**Approach:**
```go
// Example: Extract H3 Requirement headings
func ParseRequirements(filePath string) ([]RequirementBlock, error) {
    source, _ := os.ReadFile(filePath)
    doc := goldmark.New().Parser().Parse(text.NewReader(source))

    var requirements []RequirementBlock
    ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
        if !entering || n.Kind() != ast.KindHeading {
            return ast.WalkContinue, nil
        }
        heading := n.(*ast.Heading)
        if heading.Level != 3 {
            return ast.WalkContinue, nil
        }
        // Extract "Requirement: X" text...
        return ast.WalkContinue, nil
    })
    return requirements, nil
}
```

**Rationale:**
- Single-pass traversal for all extraction needs
- AST provides structural guarantees (no false matches in code blocks)
- Clean separation: parsing (goldmark) vs extraction logic (our code)

**Alternatives Considered:**
1. **Regex on rendered HTML** - Loses source line numbers, adds complexity
2. **Custom AST builder** - Reinvents goldmark, defeats purpose

### Decision 4: Incremental Migration Strategy
**Choice:** Implement feature-by-feature with parallel testing

**Phases:**
1. **Phase 1**: Add goldmark dependency, create AST helper utilities
2. **Phase 2**: Reimplement simple parsers (ExtractTitle, CountTasks)
3. **Phase 3**: Reimplement complex parsers (ParseRequirements, ParseDeltaSpec)
4. **Phase 4**: Update validation and archive packages
5. **Phase 5**: Remove old regex code after full validation

**Rationale:**
- Lower risk than big-bang rewrite
- Can validate each phase against existing tests
- Easy to roll back individual phases if issues arise

**Alternatives Considered:**
1. **Big-bang rewrite** - High risk; hard to debug failures
2. **Parallel packages** - Code duplication, unclear migration path

### Decision 5: Error Handling Enhancement
**Choice:** Leverage AST position information for detailed error messages

**Before (regex):**
```
Error: Requirement missing scenarios in spec.md
```

**After (goldmark AST):**
```
Error: Requirement missing scenarios in spec.md:42:4
    ### Requirement: User Authentication
    ^~~~ Missing "#### Scenario:" block
```

**Rationale:**
- AST nodes carry source position (line, column)
- Dramatically improves debugging experience
- Low implementation cost (AST provides info for free)

**Alternatives Considered:**
1. **Keep generic messages** - Wastes opportunity for better UX
2. **Add position tracking to regex** - Complex, fragile

## Architecture

### Before (Regex-Based)
```
┌─────────────────┐
│   spec.md       │
└────────┬────────┘
         │ os.ReadFile + bufio.Scanner
         ▼
┌─────────────────────────────────────┐
│  Regex Patterns (28 total)          │
│  • H1: `^# (.+)$`                   │
│  • H3: `^###\s+Requirement:\s*(.+)` │
│  • H4: `^####\s+Scenario:\s*(.+)`   │
│  • Tasks: `^\s*-\s*\[([xX ])\]`     │
│  • Deltas: `^##\s+(ADDED|...)`     │
└────────┬────────────────────────────┘
         │ Pattern matching
         ▼
┌─────────────────┐
│ Parsed Structs  │
│ • RequirementBlock
│ • DeltaPlan
│ • TaskStatus
└─────────────────┘
```

### After (Goldmark AST)
```
┌─────────────────┐
│   spec.md       │
└────────┬────────┘
         │ os.ReadFile
         ▼
┌─────────────────┐
│ Goldmark Parser │
│ (CommonMark AST)│
└────────┬────────┘
         │ Single parse
         ▼
┌─────────────────────────────────────┐
│  AST Tree                            │
│  Document                            │
│   ├─ Heading(level=1) "Title"       │
│   ├─ Heading(level=2) "Purpose"     │
│   ├─ Heading(level=3) "Req: X"      │
│   │   └─ Heading(level=4) "Scenario"│
│   └─ List                            │
│       └─ ListItem [x] Task           │
└────────┬────────────────────────────┘
         │ AST Walker
         ▼
┌─────────────────┐
│ Parsed Structs  │
│ • RequirementBlock
│ • DeltaPlan
│ • TaskStatus
└─────────────────┘
```

### Key Components

**1. AST Helper Utilities** (`internal/parsers/ast_helpers.go`)
- `walkHeadings(doc, level, filter)`: Generic heading extractor
- `extractTextContent(node)`: Convert AST text to string
- `getSourcePosition(node)`: Get line/column for errors

**2. Updated Parsers** (`internal/parsers/*.go`)
- `parsers.go`: ExtractTitle, CountTasks, CountDeltas
- `requirement_parser.go`: ParseRequirements, ParseScenarios
- `delta_parser.go`: ParseDeltaSpec, extract delta sections

**3. Test Compatibility Layer** (`internal/parsers/*_test.go`)
- All existing tests run against goldmark implementation
- Add edge case tests (code blocks, escaping, unicode)
- Benchmark suite to verify performance

## Risks / Trade-offs

### Risk 1: Parsing Behavior Divergence
**Description:** Goldmark may interpret edge cases differently than regex

**Impact:** Existing specs might parse differently, causing validation changes

**Mitigation:**
- Run dual implementation on entire test suite
- Test against all existing spec files in `spectr/specs/` and `spectr/changes/archive/`
- Document any intentional behavior changes with rationale

**Likelihood:** Medium | **Severity:** Medium

### Risk 2: Performance Regression
**Description:** AST parsing might be slower than targeted regex for small files

**Impact:** CLI commands feel sluggish on common operations

**Mitigation:**
- Benchmark `spectr validate`, `spectr list`, `spectr show` before/after
- Target: <5% performance degradation on typical specs (<50KB)
- Profile and optimize hot paths if needed

**Likelihood:** Low | **Severity:** Low

### Risk 3: Dependency Maintenance Burden
**Description:** Adding goldmark creates ongoing maintenance obligation

**Impact:** Need to track security issues, version updates, breaking changes

**Mitigation:**
- Goldmark is stable (v1.x since 2019, semantic versioning)
- Widely used (Hugo, Gitea) means community will catch issues
- Pin to minor version; test before major upgrades

**Likelihood:** Low | **Severity:** Low

### Trade-off: Complexity vs Robustness
**Trade-off:** Goldmark adds dependency complexity but dramatically improves robustness

**Analysis:**
- **Pro**: Handles all CommonMark edge cases correctly (code blocks, escaping, nesting)
- **Pro**: Reduces our regex maintenance from 28 patterns to 0
- **Con**: External dependency to monitor and update
- **Con**: Learning curve for goldmark AST API

**Decision:** Accept complexity for robustness gains; benefits outweigh costs

### Trade-off: Migration Effort vs Long-term Benefit
**Trade-off:** Significant upfront work to migrate, but long-term maintenance savings

**Analysis:**
- **Cost**: ~3-5 days to reimplement parsers and validate
- **Benefit**: Eliminates regex debugging burden (estimated 1-2 days/year)
- **Benefit**: Easier to add new features (e.g., spec templates, auto-formatting)
- **ROI**: Positive after ~2 years; sooner if we add new parsing features

**Decision:** Upfront investment justified by long-term maintainability

## Migration Plan

### Phase 1: Setup & Foundation (Day 1)
1. Add goldmark dependency: `go get github.com/yuin/goldmark@v1.7.x`
2. Create `internal/parsers/ast_helpers.go` with utilities
3. Add integration test comparing regex vs goldmark on sample specs

### Phase 2: Simple Parsers (Day 2)
1. Reimplement `ExtractTitle()` using goldmark
2. Reimplement `CountTasks()` using goldmark
3. Validate against existing tests
4. Benchmark performance

### Phase 3: Requirement Parsing (Day 3)
1. Reimplement `ParseRequirements()` using AST walker
2. Reimplement `ParseScenarios()` using AST walker
3. Update `internal/validation/parser.go` if needed
4. Validate against all spec files

### Phase 4: Delta Parsing (Day 4)
1. Reimplement `ParseDeltaSpec()` using goldmark
2. Reimplement delta section extractors (ADDED/MODIFIED/etc.)
3. Update `internal/archive/spec_merger.go` for AST compatibility
4. Validate against change deltas

### Phase 5: Integration & Cleanup (Day 5)
1. Run full test suite against goldmark implementation
2. Performance benchmarks on all CLI commands
3. Update error messages to include position info
4. Remove old regex code
5. Update documentation

### Rollback Plan
If critical issues discovered after deployment:
1. Keep regex implementation in separate branch
2. Revert to regex via feature flag or direct rollback
3. File issues for goldmark integration problems
4. Re-implement with fixes before next attempt

## Performance Targets

### Baseline (Regex)
- `spectr validate`: ~50ms for typical spec file (5KB)
- `spectr list`: ~100ms for 10 changes with deltas
- `spectr show`: ~30ms for single spec display

### Target (Goldmark)
- No more than 5% degradation on typical files
- No more than 10% degradation on large files (>100KB)
- Improved performance on multi-section files (less re-scanning)

### Monitoring
- Add benchmarks to CI pipeline
- Track p50, p95, p99 latencies on standard corpus
- Alert on >10% regression

## Open Questions

### Q1: Should we enable goldmark extensions?
**Context:** Goldmark supports tables, strikethrough, task lists as extensions

**Decision Needed:** Enable now or wait for user demand?

**Recommendation:** Wait. Keep scope focused on CommonMark core. Add extensions in future change if needed.

---

### Q2: How do we handle malformed markdown?
**Context:** Goldmark is permissive (per CommonMark spec); regex is strict

**Decision Needed:** Should we add extra validation on top of goldmark?

**Recommendation:** Trust goldmark's parsing. If specific validation needed, add it in validation package, not parser.

---

### Q3: Should we expose AST to other packages?
**Context:** Could let validation package work directly with AST

**Decision Needed:** Expose ast.Node types or keep parsers encapsulated?

**Recommendation:** Keep encapsulated for now. Parser package returns structs, hides AST. Can expose later if clear benefit emerges.
