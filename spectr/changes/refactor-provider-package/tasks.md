# Implementation Tasks

## 1. Stand Up ProviderKit (Shared Utilities)
- [ ] 1.1 Create `internal/providerkit` package with Provider interface alias, Configurator type, and documentation
- [ ] 1.2 Move marker constants + `UpdateFileWithMarkers`, helper validation functions, and filesystem helpers consumed by providers into ProviderKit
- [ ] 1.3 Relocate `TemplateManager` and slash base implementation into ProviderKit so providers no longer import `internal/init`
- [ ] 1.4 Add focused unit tests that cover marker updates, template rendering, and slash base behaviors within ProviderKit

## 2. Build Provider Registry & Metadata
- [ ] 2.1 Create `internal/providers/registry.go` that exposes registration APIs, metadata structs (name, type, priority, file outputs, auto-install relationships), and thread-safe lookup
- [ ] 2.2 Add validation tests covering duplicate IDs, empty metadata, dependency wiring, and introspection helpers (ListProviders, ListDefinitions)
- [ ] 2.3 Create helper constructors for common metadata shapes (config provider with markdown path, slash provider with proposal/apply/archive paths)

## 3. Extract Base Slash Command Provider
- [ ] 3.1 Port the SlashCommandConfigurator struct + helpers into ProviderKit with new typed configuration struct
- [ ] 3.2 Update implementation to consume ProviderKit filesystem + template helpers
- [ ] 3.3 Add unit tests for configure/update flows, marker enforcement, and error handling

## 4. Extract Config-Based Providers
- [ ] 4.1 Move Claude, Cline, Qwen, Qoder, CodeBuddy, and Costrict configurators into `internal/providers/{name}.go`
- [ ] 4.2 Ensure each provider registers itself with metadata: human name, priority, config paths, slash auto-installs
- [ ] 4.3 Write unit tests validating metadata wiring and Configure/IsConfigured behavior for each provider

## 5. Extract Slash Providers
- [ ] 5.1 Create individual files for each slash provider factory (Claude through Qwen)
- [ ] 5.2 Register each slash provider with metadata describing command files and display names
- [ ] 5.3 Cover factories with unit tests ensuring each returns a ProviderKit base with expected config

## 6. Refactor Init/Wizard Integration
- [ ] 6.1 Update executor to replace `getConfigurator` switch with registry lookups and metadata-driven file tracking
- [ ] 6.2 Remove `configToSlashMapping`; use metadata dependencies for auto-installation
- [ ] 6.3 Update wizard to source tool lists from provider registry metadata (IDs, names, priorities) instead of `ToolRegistry`
- [ ] 6.4 Delete the legacy `ToolRegistry` implementation and tests once wizard/executor rely solely on provider metadata
- [ ] 6.5 Adjust imports in init package to consume ProviderKit utilities where needed

## 7. Validation & Documentation
- [ ] 7.1 Run `go test ./...`, `go build`, and golangci-lint to ensure regressions are caught
- [ ] 7.2 Perform manual `spectr init` smoke tests covering config-based providers and slash auto-installs
- [ ] 7.3 Document ProviderKit contracts + registry usage within package comments and README/design updates
