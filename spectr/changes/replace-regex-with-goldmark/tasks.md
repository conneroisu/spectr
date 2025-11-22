# Implementation Tasks

## 1. Foundation
- [ ] 1.1 Add goldmark dependency to go.mod
- [ ] 1.2 Create internal/parsers/ast_walker.go with utilities for walking AST
- [ ] 1.3 Create internal/parsers/line_tracker.go for segment-to-line conversion
- [ ] 1.4 Write tests for walker and line tracker utilities

## 2. Simple Parsers
- [ ] 2.1 Migrate ExtractTitle() to use goldmark AST (find first H1 node)
- [ ] 2.2 Migrate CountRequirements(), CountDeltas(), CountTasks() to AST traversal
- [ ] 2.3 Update parsers_test.go to test goldmark implementations
- [ ] 2.4 Run parallel tests comparing regex vs goldmark output

## 3. Requirements Parsing
- [ ] 3.1 Update RequirementBlock struct with ASTNode and LineNumber fields
- [ ] 3.2 Migrate ParseRequirements() to walk AST for H3 "Requirement:" nodes
- [ ] 3.3 Migrate ParseScenarios() to walk AST for H4 "Scenario:" nodes
- [ ] 3.4 Update requirement_parser_test.go for goldmark
- [ ] 3.5 Test requirement extraction with edge cases (code blocks, nested lists)

## 4. Delta Section Parsing
- [ ] 4.1 Implement AST walker to find H2 delta section headers (ADDED, MODIFIED, etc.)
- [ ] 4.2 Parse ADDED and MODIFIED sections using AST requirement collection
- [ ] 4.3 Parse REMOVED section extracting requirement names
- [ ] 4.4 Parse RENAMED section with list post-processing for FROM:/TO: format
- [ ] 4.5 Update DeltaPlan struct with source location metadata
- [ ] 4.6 Update delta_parser_test.go for goldmark

## 5. Integration and Testing
- [ ] 5.1 Update internal/archive/spec_merger.go to use new parser structures
- [ ] 5.2 Update internal/validation/change_rules.go to use goldmark parsers
- [ ] 5.3 Update internal/validation/delta_validators.go for new structures
- [ ] 5.4 Update internal/list/lister.go and internal/view/dashboard.go
- [ ] 5.5 Run full test suite (go test ./...)
- [ ] 5.6 Performance benchmark goldmark vs regex on real spec files
- [ ] 5.7 Remove old regex parsing code from parsers.go, requirement_parser.go, delta_parser.go
- [ ] 5.8 Update any documentation referencing old parsing approach
