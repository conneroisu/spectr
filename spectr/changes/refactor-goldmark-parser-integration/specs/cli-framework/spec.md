## MODIFIED Requirements

### Requirement: Data Reuse from Discovery and Parsers
The view command SHALL reuse existing discovery and parsing infrastructure to avoid code duplication, using goldmark-based AST parsing for robust markdown extraction.

#### Scenario: Discover changes and specs
- **WHEN** building dashboard data
- **THEN** use `internal/discovery` package functions to find changes
- **AND** use `internal/discovery` package functions to find specs
- **AND** exclude archived changes from active/completed lists

#### Scenario: Parse titles and counts
- **WHEN** extracting metadata from markdown files
- **THEN** use `internal/parsers` package to parse proposal.md for titles using goldmark AST
- **AND** use `internal/parsers` package to parse spec.md for titles and requirement counts using goldmark AST
- **AND** use `internal/parsers` package to parse tasks.md for task counts
- **AND** goldmark-based parsing SHALL correctly handle edge cases like headers in code blocks

#### Scenario: AST-based title extraction
- **WHEN** extracting titles from markdown files
- **THEN** the system SHALL walk the goldmark AST to find heading nodes
- **AND** SHALL extract text content from heading nodes
- **AND** SHALL ignore headers that appear in code blocks or quoted sections
