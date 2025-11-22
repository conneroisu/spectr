## MODIFIED Requirements

### Requirement: Data Reuse from Discovery and Parsers
The view command SHALL reuse existing discovery and parsing infrastructure based on goldmark AST traversal to avoid code duplication and ensure robust markdown parsing.

#### Scenario: Discover changes and specs
- **WHEN** the view command needs to list available items
- **THEN** it SHALL use the discovery package functions
- **AND** SHALL NOT reimplement file walking or filtering logic

#### Scenario: Parse requirement counts
- **WHEN** displaying spec details
- **THEN** it SHALL use the parsers package to count requirements
- **AND** parsing SHALL use goldmark AST for accurate CommonMark-compliant extraction

#### Scenario: Extract delta information
- **WHEN** displaying change details
- **THEN** it SHALL use the parsers package to extract delta operations
- **AND** SHALL NOT duplicate parsing logic in the view package
- **AND** AST-based parsing SHALL handle edge cases correctly (code blocks, escaping, nested structures)

#### Scenario: Handle markdown edge cases
- **WHEN** spec files contain code blocks with markdown-like syntax
- **THEN** the parser SHALL correctly distinguish code from actual markdown structure
- **AND** the view command SHALL display accurate counts and information
- **AND** this SHALL prevent false positives from pattern-based parsing
