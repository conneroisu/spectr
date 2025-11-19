## MODIFIED Requirements

### Requirement: Flat Tool List in Initialization Wizard

The initialization wizard SHALL present all AI tool options in a single unified flat list without visual grouping by tool type. Each tool in the list is now a consolidated "primary tool" that automatically installs all integration variants (both config files and slash commands).

#### Scenario: Display all tools in single list

- **WHEN** user runs `spectr init` and reaches the tool selection screen
- **THEN** all 6 primary AI tools are displayed in a single flat list
- **AND** tools are sorted by priority (1-6)
- **AND** no section headers (e.g., "Config-Based Tools", "Slash Command Tools") are shown
- **AND** each tool appears as a single checkbox item with its name
- **AND** each tool listed: Claude Code, Cline, CoStrict, Qoder, CodeBuddy, Qwen

#### Scenario: Keyboard navigation across all tools

- **WHEN** user navigates with arrow keys (↑/↓)
- **THEN** the cursor moves through all 6 tools sequentially
- **AND** navigation is continuous without group boundaries
- **AND** the first tool is selected by default on screen load

#### Scenario: Tool selection works uniformly

- **WHEN** user presses space to toggle any tool
- **THEN** the checkbox state changes (checked/unchecked)
- **AND** all 6 tools behave identically
- **AND** selection state is preserved when navigating

#### Scenario: Bulk selection operations

- **WHEN** user presses 'a' to select all
- **THEN** all 6 tools are checked
- **AND** WHEN user presses 'n' to select none
- **THEN** all 6 tools are unchecked
- **AND** operations work across all tools

#### Scenario: Tool completion installs all variants

- **WHEN** user selects "Claude Code" and completes initialization
- **THEN** the following files are created:
  - `/ (project root)
    └── CLAUDE.md` - Project config with spectr markers
  - `.claude/commands/spectr/proposal.md` - Slash command for proposals
  - `.claude/commands/spectr/apply.md` - Slash command for applying changes
  - `.claude/commands/spectr/archive.md` - Slash command for archiving
- **AND** user receives one selection, but all integration variants are installed
- **AND** completion message reflects all created files

#### Scenario: Help text clarity

- **WHEN** the tool selection screen is displayed
- **THEN** the help text shows keyboard controls (↑/↓, space, a, n, enter, q)
- **AND** the help text indicates that each tool "includes all integration variants"
- **AND** the screen title clearly indicates "Select AI Tools to Configure (6 tools available)"

#### Scenario: Backward compatibility note

- **WHEN** user reviews initialization documentation
- **THEN** they understand that separate slash command tool selection is no longer available
- **AND** to get full Claude integration, select "Claude Code" instead of selecting both "Claude Code" and "Claude"
- **AND** existing installations are unaffected by this change (re-running init updates existing files)

