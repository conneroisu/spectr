## MODIFIED Requirements

### Requirement: Spec File Validation
The validation system SHALL validate spec files for structural correctness and adherence to Spectr conventions using AST-based markdown parsing.

#### Scenario: Valid spec with all required sections
- **WHEN** a spec file contains Purpose and Requirements sections with properly formatted requirements and scenarios
- **THEN** validation SHALL pass with no errors
- **AND** the validation report SHALL indicate valid=true

#### Scenario: Missing Purpose section
- **WHEN** a spec file lacks a "## Purpose" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL indicate which section is missing
- **AND** the error message SHALL include remediation guidance showing correct format

#### Scenario: Missing Requirements section
- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL provide example of correct structure

#### Scenario: Requirement without scenarios
- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example scenario format

#### Scenario: Requirement missing SHALL or MUST
- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL suggest using normative language

#### Scenario: Incorrect scenario format
- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR
- **AND** the message SHALL show the correct "#### Scenario:" header format

#### Scenario: Requirements in code blocks ignored
- **WHEN** text matching "### Requirement:" appears inside a fenced code block
- **THEN** validation SHALL NOT treat it as a requirement
- **AND** it SHALL be ignored during parsing and validation

#### Scenario: Malformed markdown structure
- **WHEN** a spec file contains invalid markdown structure (e.g., unclosed code blocks, malformed headers)
- **THEN** validation SHALL detect the structural issues during AST parsing
- **AND** SHALL report an ERROR with the specific structural problem
- **AND** the error message SHALL include line number context from the AST

### Requirement: Change Delta Validation
The validation system SHALL validate change delta specs for structural correctness and delta operation validity using AST-based markdown parsing.

#### Scenario: Valid change with deltas
- **WHEN** a change directory contains specs with proper ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** validation SHALL pass with no errors
- **AND** each delta requirement SHALL be counted toward the total

#### Scenario: Change with no deltas
- **WHEN** a change directory has no specs/ subdirectory or no delta sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL explain that at least one delta is required
- **AND** remediation guidance SHALL explain the delta header format

#### Scenario: Delta sections present but empty
- **WHEN** delta sections exist (## ADDED Requirements) but contain no requirement entries
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which sections are empty
- **AND** guidance SHALL explain requirement block format

#### Scenario: ADDED requirement without scenario
- **WHEN** an ADDED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which requirement is missing scenarios

#### Scenario: MODIFIED requirement without scenario
- **WHEN** a MODIFIED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL require at least one scenario for MODIFIED requirements

#### Scenario: Duplicate requirement in same section
- **WHEN** two requirements with the same normalized name appear in the same delta section
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL identify the duplicate requirement name

#### Scenario: Cross-section conflicts
- **WHEN** a requirement appears in both ADDED and MODIFIED sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate the conflicting requirement and sections

#### Scenario: RENAMED requirement validation
- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Delta sections in code blocks ignored
- **WHEN** text matching "## ADDED Requirements" or similar delta headers appears inside a fenced code block
- **THEN** validation SHALL NOT treat it as a delta section
- **AND** it SHALL be ignored during parsing and validation

#### Scenario: Nested delta structure detection
- **WHEN** a delta spec contains nested lists or complex markdown structure
- **THEN** the AST parser SHALL correctly identify requirement boundaries
- **AND** SHALL validate that each requirement is properly structured
- **AND** SHALL NOT confuse nested list items with separate requirements

### Requirement: Helpful Error Messages
The validation system SHALL provide actionable error messages with remediation guidance and AST context information.

#### Scenario: Error with remediation steps
- **WHEN** validation fails due to missing required content
- **THEN** the error message SHALL explain what is wrong
- **AND** SHALL provide "Next steps" section with concrete actions
- **AND** SHALL include format examples when applicable

#### Scenario: Ambiguous item name
- **WHEN** user provides an item name that matches both a change and a spec
- **THEN** validation SHALL report the ambiguity
- **AND** SHALL suggest using --type flag to disambiguate
- **AND** SHALL show available type options (change, spec)

#### Scenario: Item not found with suggestions
- **WHEN** user provides an item name that does not exist
- **THEN** validation SHALL report item not found
- **AND** SHALL provide nearest match suggestions based on string similarity
- **AND** SHALL limit suggestions to 5 most similar items

#### Scenario: Error with AST line number context
- **WHEN** validation encounters a structural error in markdown
- **THEN** the error message SHALL include the line number from the AST node
- **AND** SHALL include surrounding context when available
- **AND** SHALL help users locate the exact position of the problem in their file

#### Scenario: Parse error with markdown structure details
- **WHEN** the AST parser encounters malformed markdown during validation
- **THEN** the error message SHALL describe the specific markdown structure issue
- **AND** SHALL include the type of AST node that failed to parse correctly
- **AND** SHALL provide guidance on correct markdown formatting
