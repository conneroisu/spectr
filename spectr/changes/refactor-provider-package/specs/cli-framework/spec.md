# CLI Framework Specification Deltas

## ADDED Requirements

### Requirement: Provider Package Architecture

The system SHALL organize provider implementations in a separate `internal/providers` package with one file per provider, using a registry-based discovery mechanism instead of hardcoded switch statements.

#### Scenario: Providers organized in dedicated package
- **WHEN** the codebase is inspected
- **THEN** all provider implementations exist in `internal/providers/` package
- **AND** each provider (claude, cline, qwen, etc.) has its own .go file
- **AND** the package contains 19+ provider files plus registry infrastructure
- **AND** `internal/init/configurator.go` does not contain concrete provider implementations

#### Scenario: Provider files are focused and maintainable
- **WHEN** a developer needs to modify a specific provider
- **THEN** they can locate it in `internal/providers/{name}.go`
- **AND** the file contains only that provider's implementation
- **AND** the file is less than 100 lines (typical case)
- **AND** no navigation through large monolithic files is required

#### Scenario: New providers added without modifying existing code
- **WHEN** a new AI tool provider is added
- **THEN** a new file is created in `internal/providers/`
- **AND** the provider registers itself via `init()` function
- **AND** NO modifications to executor.go switch statements are required
- **AND** NO modifications to configurator.go are required
- **AND** the provider appears in `spectr init` automatically

### Requirement: Provider Interface

The system SHALL define a Provider interface (aliased as Configurator for backward compatibility) inside `internal/providerkit` so both orchestrator and provider packages share the same contract without import cycles.

#### Scenario: Provider interface definition
- **WHEN** the provider kit package is examined
- **THEN** it defines a Provider interface with methods: Configure, IsConfigured, GetName
- **AND** Configure accepts projectPath and spectrDir parameters
- **AND** Configure returns an error if configuration fails
- **AND** IsConfigured checks if the provider's files already exist
- **AND** GetName returns a human-readable provider name

#### Scenario: Both provider types implement same interface
- **WHEN** a config-based provider (e.g., ClaudeCodeConfigurator) is examined
- **THEN** it implements the Provider interface
- **AND** WHEN a slash-command provider is examined
- **THEN** it also implements the Provider interface
- **AND** both types are usable interchangeably by the executor

#### Scenario: Backward compatible Configurator alias
- **WHEN** existing code uses the Configurator interface
- **THEN** it continues to work without modification
- **AND** Configurator is a type alias for Provider
- **AND** no import changes are required in wizard.go or other consumers

### Requirement: Provider Registry

The system SHALL provide a global provider registry that maps provider IDs to provider implementations, enabling runtime lookup without hardcoded switch statements.

#### Scenario: Registry initialization
- **WHEN** the application starts
- **THEN** the provider registry is empty initially
- **AND** providers register themselves during package init phase
- **AND** all 19+ providers are registered before main() executes
- **AND** the registry is thread-safe (uses sync.RWMutex)

#### Scenario: Provider registration
- **WHEN** a provider file is loaded
- **THEN** its init() function calls Register(id, provider)
- **AND** the provider is stored in the global registry map
- **AND** duplicate registrations return an error
- **AND** empty IDs return an error
- **AND** nil providers return an error

#### Scenario: Provider lookup by ID
- **WHEN** executor needs to get a provider by ID (e.g., "claude-code")
- **THEN** it calls providers.GetProvider("claude-code")
- **AND** the registry returns the registered provider instance
- **AND** if the ID is not found, returns nil and an error
- **AND** lookup is thread-safe

#### Scenario: Registry introspection
- **WHEN** code needs to list all registered providers
- **THEN** it calls providers.ListProviders()
- **AND** receives structured metadata (ID, name, type, priority, output files, auto-install relationships)
- **AND** the list is sorted alphabetically by ID as a convenience API
- **AND** useful for debugging, validation, and UI rendering

#### Scenario: Registry supplies wizard metadata
- **WHEN** the init wizard requests tool options
- **THEN** it can read friendly names, help text, and priorities from the registry definitions
- **AND** the wizard never maintains a parallel hardcoded list
- **AND** adding a new provider entry automatically propagates to the UI

### Requirement: Provider Metadata

The registry SHALL store the metadata necessary for CLI UX (names, priorities), execution (file paths, slash auto-install), and dependency handling alongside each provider.

#### Scenario: Metadata describes files produced
- **WHEN** a provider registers, it specifies the absolute/relative paths it writes
- **AND** executor can report created vs updated files without poking concrete types
- **AND** slash providers declare their three command files explicitly

#### Scenario: Config-to-slash mapping lives in metadata
- **WHEN** a config provider depends on slash commands
- **THEN** its metadata declares which slash provider(s) should auto-install
- **AND** executor reads this relationship from the registry instead of a map
- **AND** removing or adding a mapping occurs entirely within provider registration

#### Scenario: Wizard ordering derives from metadata
- **WHEN** the wizard sorts the checkbox list
- **THEN** it uses metadata.Priority values
- **AND** no additional constants or manual indices exist outside the registry

### Requirement: Config-Based Providers

The system SHALL support config-based providers that create single markdown configuration files in the project root.

#### Scenario: Config provider file structure
- **WHEN** a config-based provider (e.g., claude.go) is examined
- **THEN** it defines a struct (e.g., ClaudeProvider) implementing Provider interface
- **AND** the file contains Configure, IsConfigured, and GetName methods
- **AND** the file contains an init() function registering the provider
- **AND** the provider uses marker-based file updates from internal/init

#### Scenario: Config provider creates markdown file
- **WHEN** Configure() is called on a config-based provider
- **THEN** it creates or updates a markdown file (e.g., CLAUDE.md)
- **AND** uses UpdateFileWithMarkers() for safe content injection
- **AND** the file contains project-specific instructions for the AI tool
- **AND** marker boundaries (SpectrStartMarker, SpectrEndMarker) wrap the content

#### Scenario: Config provider checks existing configuration
- **WHEN** IsConfigured() is called on a config-based provider
- **THEN** it checks if the target markdown file exists
- **AND** verifies the file contains Spectr markers
- **AND** returns true if properly configured
- **AND** returns false if file missing or markers missing

### Requirement: Slash Command Providers

The system SHALL support slash command providers that create multiple files in the .claude/commands/ directory.

#### Scenario: Slash provider file structure
- **WHEN** a slash command provider (e.g., claude_slash.go) is examined
- **THEN** it defines a factory function (e.g., NewClaudeSlashConfigurator)
- **AND** returns a base slash provider configured with specific parameters
- **AND** the file contains an init() function registering the provider
- **AND** base implementation is in base_slash.go

#### Scenario: Slash provider creates command files
- **WHEN** Configure() is called on a slash command provider
- **THEN** it creates .claude/commands/ directory if needed
- **AND** creates 3 files: proposal.md, apply.md, archive.md
- **AND** uses template rendering system for file contents
- **AND** templates are parameterized with provider-specific paths

#### Scenario: Slash provider checks existing configuration
- **WHEN** IsConfigured() is called on a slash command provider
- **THEN** it checks if all 3 command files exist
- **AND** verifies each file contains required content
- **AND** returns true if all files properly configured
- **AND** returns false if any file missing or incomplete

### Requirement: Base Slash Provider

The system SHALL provide a reusable base implementation for slash command providers to eliminate code duplication across 15+ slash providers.

#### Scenario: Base slash provider structure
- **WHEN** base_slash.go is examined
- **THEN** it defines a SlashCommandProvider struct
- **AND** the struct has fields for toolName, configPath, commandsPath
- **AND** it implements the Provider interface
- **AND** all slash providers use this base implementation

#### Scenario: Factory pattern for slash providers
- **WHEN** a specific slash provider is needed (e.g., Claude slash)
- **THEN** the factory function NewClaudeSlashConfigurator() is called
- **AND** it returns &SlashCommandProvider{toolName: "claude-code", ...}
- **AND** the returned instance is ready to use
- **AND** no subclassing or embedding required

#### Scenario: Base slash provider Configure implementation
- **WHEN** Configure() is called on a SlashCommandProvider
- **THEN** it renders templates for proposal, apply, and archive commands
- **AND** creates .claude/commands/ directory structure
- **AND** writes all 3 markdown files
- **AND** returns error if any step fails
- **AND** implementation is shared across all slash providers

### Requirement: Executor Registry Integration

The system SHALL update the executor to use registry-based provider lookup instead of hardcoded switch statements.

#### Scenario: Executor uses registry lookup
- **WHEN** executor.getConfigurator(toolID) is called
- **THEN** it calls providers.GetProvider(toolID)
- **AND** returns the provider from the registry
- **AND** returns error if provider not found
- **AND** no switch statement exists in the code

#### Scenario: No hardcoded provider references
- **WHEN** internal/init/executor.go is examined
- **THEN** it does not contain switch statements on tool IDs
- **AND** it does not import concrete provider types
- **AND** it only uses the Provider interface
- **AND** all provider-specific logic is in providers package

#### Scenario: Backward compatible behavior
- **WHEN** spectr init is run with any tool selection
- **THEN** behavior is identical to before refactoring
- **AND** same files are created with same contents
- **AND** same error messages are shown
- **AND** wizard interaction is unchanged

### Requirement: Shared ProviderKit Utilities

The system SHALL centralize shared provider helpers (marker utilities, template manager, slash base implementation) in a dedicated `internal/providerkit` package to prevent import cycles.

#### Scenario: Marker utilities location
- **WHEN** providerkit is examined
- **THEN** it exposes UpdateFileWithMarkers(), SpectrStartMarker, SpectrEndMarker, and helper validation functions
- **AND** these helpers are the single source of truth used by every provider implementation

#### Scenario: Template manager location
- **WHEN** providerkit is inspected
- **THEN** the embedded template manager used by providers and executor lives there
- **AND** providers access templates without referencing the init package
- **AND** init orchestrator uses the same manager so template parsing logic stays consistent

#### Scenario: Slash base implementation
- **WHEN** base slash provider logic is needed
- **THEN** providerkit exposes it so slash providers can reuse the implementation
- **AND** init executor does not need to import concrete slash types directly
