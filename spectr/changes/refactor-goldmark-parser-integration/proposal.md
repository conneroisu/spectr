# Change: Replace regex parsing with goldmark AST-based parsing

## Why

The current markdown parsing implementation in `internal/parsers/` uses regex-based pattern matching to extract structure from spec files. While functional, this approach has limitations in robustness, maintainability, and extensibility. Goldmark provides a proper CommonMark-compliant AST that handles edge cases correctly, reduces custom regex logic, and provides a cleaner foundation for future parsing enhancements.

## What Changes

- Replace regex-based parsing with goldmark AST traversal in `internal/parsers/`
- Update `internal/validation/parser.go` to use goldmark AST
- Update `internal/archive/spec_merger.go` to use new parsing interface
- Add `github.com/yuin/goldmark` dependency to `go.mod`
- Maintain backward-compatible API for all existing consumers
- Add comprehensive test coverage for edge cases that regex parsing missed
- Document goldmark integration patterns in code

## Impact

**Affected specs:**
- validation
- cli-framework
- archive-workflow

**Affected code:**
- `internal/parsers/parsers.go` - Core parsing utilities (ExtractTitle, CountTasks, CountDeltas, CountRequirements)
- `internal/parsers/requirement_parser.go` - Requirement and scenario extraction
- `internal/parsers/delta_parser.go` - Delta operation parsing (ADDED/MODIFIED/REMOVED/RENAMED)
- `internal/validation/parser.go` - Validation-specific parsing (ExtractSections, ExtractRequirements, ExtractScenarios)
- `internal/archive/spec_merger.go` - Uses parsers package for merging
- All test files for the above modules
- 13 files total that call parsing functions

**Benefits:**
- More robust parsing with proper markdown AST
- Better handling of edge cases (nested lists, code blocks in requirements, etc.)
- Reduced maintenance burden (less custom regex)
- CommonMark compliance
- Foundation for future enhancements (frontmatter, custom extensions)
- Improved error messages with AST node location information

**Risks:**
- Dependency addition increases binary size
- Potential subtle behavior changes in edge cases
- Migration requires careful testing to ensure API compatibility
