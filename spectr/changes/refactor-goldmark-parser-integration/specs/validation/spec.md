## ADDED Requirements

### Requirement: Enhanced Error Location Information
The validation system SHALL provide precise line and column information for validation errors using AST node position data.

#### Scenario: Error with line number
- **WHEN** validation detects an error in a spec file
- **THEN** the error report SHALL include the line number where the issue occurs
- **AND** the line number SHALL be extracted from the goldmark AST node segment information

#### Scenario: Error with column information
- **WHEN** validation detects an error in a specific part of a line
- **THEN** the error report MAY include column information for precise location
- **AND** this SHALL help users quickly locate and fix issues

#### Scenario: Multi-line validation context
- **WHEN** an error involves multiple lines (e.g., missing scenario in requirement block)
- **THEN** the error SHALL report the starting line of the problematic requirement
- **AND** MAY include ending line for context

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
- **AND** the error SHALL include line number information from AST

#### Scenario: Missing Requirements section
- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL provide example of correct structure
- **AND** the error SHALL include line number information from AST

#### Scenario: Requirement without scenarios
- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example scenario format
- **AND** the warning SHALL include the requirement line number

#### Scenario: Requirement missing SHALL or MUST
- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL suggest using normative language

#### Scenario: Incorrect scenario format
- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR
- **AND** the message SHALL show the correct "#### Scenario:" header format

#### Scenario: Headers in code blocks ignored
- **WHEN** a spec file contains headers inside code blocks or inline code
- **THEN** validation SHALL correctly ignore these headers using AST structure
- **AND** SHALL only validate actual markdown heading nodes
