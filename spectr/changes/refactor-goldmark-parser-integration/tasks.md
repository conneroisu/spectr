## 1. Dependencies and Foundation

- [ ] 1.1 Add `github.com/yuin/goldmark v1.7.8` to go.mod
- [ ] 1.2 Run `go get github.com/yuin/goldmark` and verify dependency installation
- [ ] 1.3 Create internal helper functions for goldmark AST operations in parsers package
- [ ] 1.4 Implement `extractText(node ast.Node, source []byte) string` helper
- [ ] 1.5 Implement `findHeadings(source []byte, level int) []string` helper
- [ ] 1.6 Implement `walkAST(node ast.Node, visitor NodeVisitor) error` pattern helpers
- [ ] 1.7 Write unit tests for helper functions

## 2. Refactor parsers/parsers.go

- [ ] 2.1 Refactor `ExtractTitle` to use goldmark AST for finding H1 headings
- [ ] 2.2 Ensure `ExtractTitle` maintains exact same API and behavior
- [ ] 2.3 Refactor `CountTasks` to use goldmark for parsing task list items (if applicable, or keep regex for simple pattern)
- [ ] 2.4 Refactor `CountDeltas` to use goldmark for finding H2 delta section headers
- [ ] 2.5 Refactor `CountRequirements` to use goldmark for finding H3 requirement headers
- [ ] 2.6 Run existing tests in parsers_test.go and verify all pass
- [ ] 2.7 Add new edge case tests (headers in code blocks, escaped markdown, unicode)

## 3. Refactor parsers/requirement_parser.go

- [ ] 3.1 Refactor `ParseRequirements` to walk AST for H3 "### Requirement:" headings
- [ ] 3.2 Extract full requirement block content from AST (including scenarios)
- [ ] 3.3 Refactor `ParseScenarios` to walk AST for H4 "#### Scenario:" headings
- [ ] 3.4 Ensure RequirementBlock struct population matches existing behavior
- [ ] 3.5 Maintain `NormalizeRequirementName` as-is (no changes needed)
- [ ] 3.6 Run existing tests in requirement_parser_test.go and verify all pass
- [ ] 3.7 Add edge case tests for nested structures, code blocks with scenario examples

## 4. Refactor parsers/delta_parser.go

- [ ] 4.1 Refactor `ParseDeltaSpec` to use goldmark for section extraction
- [ ] 4.2 Update `extractSectionContent` to use AST walking instead of regex
- [ ] 4.3 Update `parseRequirementsFromSection` to use goldmark AST
- [ ] 4.4 Update `parseRemovedSection` to use goldmark AST
- [ ] 4.5 Update `parseRenamedSection` to use goldmark AST (parse list items with FROM/TO)
- [ ] 4.6 Ensure DeltaPlan struct population matches existing behavior
- [ ] 4.7 Run existing tests in delta_parser_test.go and verify all pass
- [ ] 4.8 Add edge case tests for complex delta specs

## 5. Refactor validation/parser.go

- [ ] 5.1 Refactor `ExtractSections` to use goldmark for H2 section extraction
- [ ] 5.2 Refactor `ExtractRequirements` to use goldmark AST
- [ ] 5.3 Update Requirement struct to optionally include line number (for enhanced errors)
- [ ] 5.4 Refactor `ExtractScenarios` to use goldmark AST
- [ ] 5.5 Maintain `ContainsShallOrMust` as-is (regex is fine for text search)
- [ ] 5.6 Maintain `NormalizeRequirementName` as-is
- [ ] 5.7 Run existing validation tests and verify all pass
- [ ] 5.8 Update validation error reporting to include line numbers from AST

## 6. Update Archive Package (archive/spec_merger.go)

- [ ] 6.1 Verify spec_merger.go works correctly with updated parsers package
- [ ] 6.2 Run archive package tests and verify all pass
- [ ] 6.3 Test end-to-end archiving workflow with goldmark-based parsing

## 7. Integration Testing

- [ ] 7.1 Run full test suite: `go test ./...`
- [ ] 7.2 Fix any failing tests due to subtle behavior changes
- [ ] 7.3 Test with real spec files in spectr/ directory
- [ ] 7.4 Validate all changes and specs: `spectr validate --all`
- [ ] 7.5 Test list command: `spectr list` and `spectr list --specs`
- [ ] 7.6 Test view command: `spectr view`
- [ ] 7.7 Test archive command on a test change

## 8. Performance and Quality

- [ ] 8.1 Benchmark parsing performance: compare old regex vs new goldmark
- [ ] 8.2 Verify binary size increase is acceptable (<5% increase)
- [ ] 8.3 Run linter: `golangci-lint run`
- [ ] 8.4 Fix any linting issues
- [ ] 8.5 Update code documentation and comments where goldmark is used

## 9. Documentation and Cleanup

- [ ] 9.1 Document goldmark integration patterns in code comments
- [ ] 9.2 Remove unused regex patterns and imports
- [ ] 9.3 Update project.md if needed to mention goldmark dependency
- [ ] 9.4 Add inline comments explaining AST walking patterns
- [ ] 9.5 Ensure all exported functions have proper godoc comments

## 10. Final Validation

- [ ] 10.1 Run `spectr validate refactor-goldmark-parser-integration --strict`
- [ ] 10.2 Verify all delta specs are valid
- [ ] 10.3 Run full test suite one final time
- [ ] 10.4 Build binary and test manually with real operations
- [ ] 10.5 Prepare for code review and approval
