# Change: Replace Regex-Based Markdown Parsing with Goldmark AST Parser

## Why

Spectr currently uses regex-based parsing in `internal/parsers/` to extract markdown structure (headings, requirements, scenarios, tasks). While functional, this approach has several limitations:

1. **Fragility**: Regex patterns are brittle and fail on edge cases like nested structures, code blocks containing markdown-like syntax, escaped characters, and non-standard whitespace
2. **Maintainability**: 28+ regex patterns scattered across 6 files create high maintenance burden and cognitive overhead
3. **Standards Compliance**: No guarantee of CommonMark compliance; custom regex may diverge from markdown spec
4. **Error Handling**: Regex provides poor error messages when parsing fails, making debugging difficult
5. **Performance**: Regex scanning entire files multiple times is inefficient compared to single-pass AST traversal

Goldmark is the de facto standard markdown parser in Go (used by Hugo, Gitea), provides full CommonMark compliance, and offers a robust AST-based API that eliminates these issues.

## What Changes

This change replaces regex-based parsing with goldmark AST traversal while maintaining backward compatibility:

### Code Changes
- **internal/parsers/**: Rewrite parsers to use goldmark AST instead of regex patterns
  - `ExtractTitle()`: Walk AST for first H1 heading node
  - `ParseRequirements()`: Extract H3 heading nodes matching "Requirement:" prefix
  - `ParseScenarios()`: Extract H4 heading nodes matching "Scenario:" prefix
  - `ParseDeltaSpec()`: Extract sections by H2 nodes (ADDED/MODIFIED/REMOVED/RENAMED)
  - `CountTasks()`: Parse list items with checkbox markers
- **internal/validation/**: Update validation rules to leverage AST structure
- **internal/archive/**: Update spec merger to use new parser API

### Dependency Changes
- **Add**: `github.com/yuin/goldmark` (stable, well-maintained, ~200KB)
- **Remove**: None (reduced regex usage is internal)

### API Changes
- Parser function signatures remain unchanged (backward compatible)
- Internal implementation uses goldmark AST instead of regex
- Return types and error handling stay consistent

### Testing Changes
- All existing parser tests must pass with goldmark implementation
- Add edge case tests (code blocks, nested lists, escaping, unicode)
- Performance benchmarks to verify no regression

## Impact

### Affected Specs
- **validation**: Core parsing logic changes affect validation behavior
- **cli-framework**: Parser package is part of CLI framework architecture

### Affected Code
- `internal/parsers/` (3 files): Complete rewrite using goldmark
- `internal/validation/parser.go`: Update regex-based section extraction
- `internal/archive/spec_merger.go`: May need updates for AST-based matching
- `cmd/` commands: No changes (use parser package interface)

### Breaking Changes
**None** - Parser interfaces remain backward compatible. Internal implementation change only.

### Migration Path
1. Implement goldmark parsers alongside existing regex parsers
2. Run both implementations against test suite to verify parity
3. Switch to goldmark implementation
4. Remove old regex code after validation

### Risks
- **Parsing behavior changes**: Goldmark may interpret edge cases differently than regex
  - *Mitigation*: Comprehensive test suite comparing both implementations
- **Performance**: AST parsing might be slower than targeted regex for small files
  - *Mitigation*: Benchmark critical paths; optimize if needed
- **Dependency**: Adding external dependency increases maintenance burden
  - *Mitigation*: Goldmark is stable (>5 years), widely used, actively maintained

### Benefits
- **Robustness**: Handles all CommonMark edge cases correctly
- **Maintainability**: Single AST traversal replaces 28 regex patterns
- **Extensibility**: Easy to add new parsing features (tables, footnotes, etc.)
- **Error Messages**: AST provides line/column info for better error reporting
- **Performance**: Single-pass parsing vs multiple regex scans
- **Standards**: Full CommonMark compliance out of the box
