## ADDED Requirements

### Requirement: Position-Aware Error Messages
The validation system SHALL provide error messages with precise source file locations (line and column numbers) for parsing and validation issues.

#### Scenario: Error message includes line number
- **WHEN** validation detects a requirement missing scenarios
- **THEN** the error message SHALL include the file path with line number (e.g., "spec.md:42")
- **AND** the message SHALL show the problematic source line
- **AND** the error SHALL be easier to locate and fix than generic messages

#### Scenario: Error message includes column position
- **WHEN** validation detects incorrect heading syntax
- **THEN** the error message SHALL include both line and column position (e.g., "spec.md:42:4")
- **AND** the message SHALL use visual indicators (^~~~) to highlight the issue
- **AND** developers SHALL quickly identify the exact location of the problem

#### Scenario: Multi-line context in error messages
- **WHEN** validation detects structural issues spanning multiple lines
- **THEN** the error message SHALL include surrounding context lines
- **AND** the problematic section SHALL be clearly marked
- **AND** developers SHALL understand the issue without opening the file

## MODIFIED Requirements

### Requirement: Spec File Validation
The validation system SHALL validate spec files for structural correctness and adherence to Spectr conventions using AST-based markdown parsing for robust edge case handling.

#### Scenario: Valid spec with all required sections
- **WHEN** a spec file contains Purpose and Requirements sections with properly formatted requirements and scenarios
- **THEN** validation SHALL pass with no errors
- **AND** the validation report SHALL indicate valid=true
- **AND** parsing SHALL handle edge cases like code blocks containing markdown syntax

#### Scenario: Missing Purpose section
- **WHEN** a spec file lacks a "## Purpose" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL indicate which section is missing with line number
- **AND** the error message SHALL include remediation guidance showing correct format

#### Scenario: Missing Requirements section
- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL provide example of correct structure with source position

#### Scenario: Requirement without scenarios
- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example scenario format with line reference

#### Scenario: Requirement missing SHALL or MUST
- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL suggest using normative language with position info

#### Scenario: Incorrect scenario format
- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR
- **AND** the message SHALL show the correct "#### Scenario:" header format
- **AND** the error SHALL include precise line and column position

#### Scenario: Code blocks containing markdown-like syntax
- **WHEN** a spec file contains code blocks with markdown syntax (e.g., ```# Example```)
- **THEN** validation SHALL NOT parse code block content as markdown structure
- **AND** validation SHALL correctly distinguish code from actual headings
- **AND** this SHALL prevent false positives from regex-based parsing

### Requirement: Change Delta Validation
The validation system SHALL validate change delta specs for structural correctness and delta operation validity using AST-based parsing.

#### Scenario: Valid change with deltas
- **WHEN** a change directory contains specs with proper ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** validation SHALL pass with no errors
- **AND** each delta requirement SHALL be counted toward the total
- **AND** parsing SHALL correctly handle CommonMark edge cases

#### Scenario: Change with no deltas
- **WHEN** a change directory has no specs/ subdirectory or no delta sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL explain that at least one delta is required
- **AND** remediation guidance SHALL explain the delta header format

#### Scenario: Delta sections present but empty
- **WHEN** delta sections exist (## ADDED Requirements) but contain no requirement entries
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which sections are empty with line numbers
- **AND** guidance SHALL explain requirement block format

#### Scenario: ADDED requirement without scenario
- **WHEN** an ADDED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which requirement is missing scenarios with position

#### Scenario: MODIFIED requirement without scenario
- **WHEN** a MODIFIED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL require at least one scenario for MODIFIED requirements with location

#### Scenario: Duplicate requirement in same section
- **WHEN** two requirements with the same normalized name appear in the same delta section
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL identify the duplicate requirement name with line numbers for both

#### Scenario: Cross-section conflicts
- **WHEN** a requirement appears in both ADDED and MODIFIED sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate the conflicting requirement and sections with positions

#### Scenario: RENAMED requirement validation
- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name
- **AND** error messages SHALL include line numbers for rename declarations
