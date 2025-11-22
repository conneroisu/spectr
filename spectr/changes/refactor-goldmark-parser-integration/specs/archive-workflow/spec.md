## MODIFIED Requirements

### Requirement: Delta Operation Parsing
The system SHALL parse delta operations from change spec files using goldmark AST-based parsing for robust and correct markdown structure extraction.

#### Scenario: Parse ADDED requirements
- **WHEN** a delta spec contains `## ADDED Requirements` section
- **THEN** the system extracts all requirement blocks with headers and scenarios using goldmark AST traversal
- **AND** SHALL correctly identify requirement heading nodes at level 3
- **AND** SHALL correctly identify scenario heading nodes at level 4

#### Scenario: Parse MODIFIED requirements
- **WHEN** a delta spec contains `## MODIFIED Requirements` section
- **THEN** the system extracts complete modified requirement blocks using goldmark AST
- **AND** SHALL preserve all content including scenarios and body text

#### Scenario: Parse REMOVED requirements
- **WHEN** a delta spec contains `## REMOVED Requirements` section
- **THEN** the system extracts requirement names to be removed using goldmark AST

#### Scenario: Parse RENAMED requirements
- **WHEN** a delta spec contains `## RENAMED Requirements` section with FROM/TO pairs
- **THEN** the system extracts the old and new requirement names using goldmark AST

#### Scenario: Require at least one delta operation
- **WHEN** a delta spec has no ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** the system returns an error indicating no delta operations were found

#### Scenario: Ignore headers in code blocks
- **WHEN** parsing delta specs containing code examples with markdown headers
- **THEN** the system SHALL use goldmark AST to distinguish between real headers and code block content
- **AND** SHALL only extract actual heading nodes
- **AND** SHALL not misparse code examples as requirements

#### Scenario: Handle escaped markdown
- **WHEN** parsing requirements containing escaped markdown characters
- **THEN** goldmark AST parsing SHALL correctly handle escaped characters
- **AND** SHALL preserve the semantic structure of the document
