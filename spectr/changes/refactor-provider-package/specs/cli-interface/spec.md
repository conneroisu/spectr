# CLI Interface Specification Deltas

## MODIFIED Requirements

### Requirement: Flat Tool List in Initialization Wizard

The initialization wizard SHALL present all AI tool options in a single unified flat list without visual grouping by tool type, with providers loaded from the registry.

#### Scenario: Display all tools in single list
- **WHEN** user runs `spectr init` and reaches the tool selection screen
- **THEN** all 17+ AI tools are displayed in a single flat list
- **AND** tools are sorted by the priority value defined in provider metadata (1-17+)
- **AND** no section headers (e.g., "Config-Based Tools", "Slash Command Tools") are shown
- **AND** each tool appears as a single checkbox item with its name
- **AND** tools are discovered from the provider registry, not hardcoded

#### Scenario: Keyboard navigation across all tools
- **WHEN** user navigates with arrow keys (↑/↓)
- **THEN** the cursor moves through all 17+ tools sequentially
- **AND** navigation is continuous without group boundaries
- **AND** the first tool is selected by default on screen load

#### Scenario: Tool selection works uniformly
- **WHEN** user presses space to toggle any tool
- **THEN** the checkbox state changes (checked/unchecked)
- **AND** both config-based and slash command tools behave identically
- **AND** selection state is preserved when navigating
- **AND** provider implementations are loaded from the registry transparently

#### Scenario: Bulk selection operations
- **WHEN** user presses 'a' to select all
- **THEN** all 17+ tools are checked
- **AND** WHEN user presses 'n' to select none
- **THEN** all 17+ tools are unchecked
- **AND** operations work across all tools regardless of type

#### Scenario: Help text clarity
- **WHEN** the tool selection screen is displayed
- **THEN** the help text shows keyboard controls (↑/↓, space, a, n, enter, q)
- **AND** the help text does NOT reference tool groupings or categories
- **AND** the screen title clearly indicates "Select AI Tools to Configure"

#### Scenario: Newly registered tools appear automatically
- **WHEN** a provider registers new metadata (ID, name, priority)
- **THEN** the wizard list reflects the change without code edits in the UI layer
- **AND** the new provider can be selected, configured, and auto-installs slash commands if specified
