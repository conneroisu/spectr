## 1. Setup & Dependencies
- [ ] 1.1 Add goldmark dependency to go.mod (`go get github.com/yuin/goldmark@v1.7.x`)
- [ ] 1.2 Verify goldmark builds successfully and passes go mod tidy
- [ ] 1.3 Create `internal/parsers/ast_helpers.go` for shared AST utilities
- [ ] 1.4 Implement helper functions: `walkHeadings()`, `extractTextContent()`, `getSourcePosition()`
- [ ] 1.5 Write unit tests for AST helper functions

## 2. Simple Parser Migration
- [ ] 2.1 Create goldmark-based implementation of `ExtractTitle()` in parsers.go
- [ ] 2.2 Run existing tests for `ExtractTitle()` to verify parity
- [ ] 2.3 Add edge case tests (code blocks with `# title`, escaped hashes, unicode)
- [ ] 2.4 Create goldmark-based implementation of `CountTasks()` in parsers.go
- [ ] 2.5 Run existing tests for `CountTasks()` to verify parity
- [ ] 2.6 Add edge case tests (nested lists, checkboxes in code blocks)
- [ ] 2.7 Benchmark `ExtractTitle()` and `CountTasks()` performance vs baseline

## 3. Requirement Parser Migration
- [ ] 3.1 Create goldmark-based implementation of `ParseRequirements()` in requirement_parser.go
- [ ] 3.2 Implement AST walker to extract H3 headings with "Requirement:" prefix
- [ ] 3.3 Extract full requirement block content including scenarios
- [ ] 3.4 Preserve line/column position information in RequirementBlock struct
- [ ] 3.5 Run existing tests for `ParseRequirements()` to verify parity
- [ ] 3.6 Add edge case tests (requirements in code blocks, escaped characters, unicode names)
- [ ] 3.7 Create goldmark-based implementation of `ParseScenarios()` in requirement_parser.go
- [ ] 3.8 Implement AST walker to extract H4 headings with "Scenario:" prefix
- [ ] 3.9 Run existing tests for `ParseScenarios()` to verify parity
- [ ] 3.10 Benchmark requirement parsing performance

## 4. Delta Parser Migration
- [ ] 4.1 Create goldmark-based implementation of `ParseDeltaSpec()` in delta_parser.go
- [ ] 4.2 Implement section extraction for H2 headers (ADDED/MODIFIED/REMOVED/RENAMED)
- [ ] 4.3 Implement `parseDeltaSection()` using AST walker instead of regex
- [ ] 4.4 Implement `parseRemovedSection()` using AST traversal
- [ ] 4.5 Implement `parseRenamedSection()` using AST traversal
- [ ] 4.6 Run existing tests for delta parsing to verify parity
- [ ] 4.7 Add edge case tests (delta sections in code blocks, complex nesting)
- [ ] 4.8 Update DeltaPlan struct to include position information if needed
- [ ] 4.9 Benchmark delta parsing performance

## 5. Validation Package Updates
- [ ] 5.1 Review `internal/validation/parser.go` for regex usage
- [ ] 5.2 Update section extraction logic to use goldmark if needed
- [ ] 5.3 Enhance error messages to include line/column information from AST
- [ ] 5.4 Update `internal/validation/change_rules.go` to leverage position data
- [ ] 5.5 Update error formatters to display position-aware messages
- [ ] 5.6 Add tests for enhanced error messages with position info
- [ ] 5.7 Run full validation test suite to verify no regressions

## 6. Archive Package Updates
- [ ] 6.1 Review `internal/archive/spec_merger.go` for parser usage
- [ ] 6.2 Verify spec merger works correctly with goldmark-based parsers
- [ ] 6.3 Update requirement matching logic if needed for AST compatibility
- [ ] 6.4 Run archive integration tests to verify merge operations
- [ ] 6.5 Test archiving on actual change proposals from `spectr/changes/`

## 7. Integration Testing
- [ ] 7.1 Run full test suite: `go test ./...`
- [ ] 7.2 Fix any failing tests and investigate behavior differences
- [ ] 7.3 Test against all existing spec files in `spectr/specs/`
- [ ] 7.4 Test against archived changes in `spectr/changes/archive/`
- [ ] 7.5 Validate active changes in `spectr/changes/` parse correctly
- [ ] 7.6 Create integration test comparing regex vs goldmark on corpus
- [ ] 7.7 Document any intentional behavior changes with rationale

## 8. Performance Validation
- [ ] 8.1 Run benchmark suite on parser functions
- [ ] 8.2 Benchmark `spectr validate` on typical spec files (5-50KB)
- [ ] 8.3 Benchmark `spectr list` on directory with 10+ changes
- [ ] 8.4 Benchmark `spectr show` on individual specs
- [ ] 8.5 Verify <5% degradation on typical files
- [ ] 8.6 Profile and optimize if performance targets not met
- [ ] 8.7 Document performance results in proposal

## 9. Code Cleanup
- [ ] 9.1 Remove old regex-based parsing code from parsers package
- [ ] 9.2 Remove unused regex imports from validation/archive packages
- [ ] 9.3 Update package documentation to mention goldmark
- [ ] 9.4 Run golangci-lint and fix any new warnings
- [ ] 9.5 Update godoc comments for public functions
- [ ] 9.6 Verify no dead code remains with code coverage analysis

## 10. Documentation & Finalization
- [ ] 10.1 Update README.md to mention goldmark dependency
- [ ] 10.2 Update project.md Tech Stack section with goldmark
- [ ] 10.3 Add migration notes to CHANGELOG if applicable
- [ ] 10.4 Update proposal.md with final performance metrics
- [ ] 10.5 Run `spectr validate` on this change proposal
- [ ] 10.6 Mark all tasks complete in this tasks.md
- [ ] 10.7 Request code review and proposal approval
