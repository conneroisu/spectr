# Design: Provider Package Architecture

## Context

The initialization system currently has all 19+ provider implementations in a single 831-line `configurator.go` file. The system supports two types of providers:
- **Config-based providers**: Create single markdown files (CLAUDE.md, CLINE.md, etc.)
- **Slash command providers**: Create 3 files each in .claude/commands/ directory

Current pain points:
- Adding new providers requires modifying 3 locations (configurator.go, executor.go switch, registry.go)
- Finding specific provider logic difficult in monolithic file
- No clear extension point for external provider plugins
- Switch statement in executor couples implementation to discovery

## Goals

- Separate provider implementations into individual files (one per provider)
- Create registry-based provider lookup to eliminate hardcoded switch
- Unify config-based and slash-command providers under common interface
- Maintain 100% backward compatibility with existing CLI behavior
- Enable future plugin architecture (providers could load from external packages)
- Reduce cognitive load when adding/modifying providers

## Non-Goals

- Plugin system implementation (just prepare architecture)
- Changes to user-facing CLI commands or flags
- Changes to file formats (CLAUDE.md, .claude/commands/, etc.)
- Performance optimization (current implementation is fast enough)
- Changes to template rendering system

## Decisions

### Decision 1: Package Location - internal/providers

**Chosen**: `internal/providers` as sibling to `internal/init`

**Rationale**:
- Providers are domain-specific logic, not general initialization utilities
- Could be reused in future by other packages (e.g., auto-update mechanism)
- Clear separation of concerns: `init` orchestrates, `providers` implements
- Follows Go convention of domain-based package organization

**Alternatives considered**:
- `internal/init/providers` - Would tightly couple to init, harder to reuse
- `internal/configurators` - Less clear name, "configurator" is implementation detail

### Decision 2: One File Per Provider

**Chosen**: Separate file for each of 19+ providers

**Rationale**:
- Easy to locate specific provider (e.g., claude.go vs searching 831-line file)
- Smaller diffs when modifying providers (clearer PR reviews)
- Natural code organization by responsibility
- Enables future per-provider testing and versioning
- Follows pattern seen in Go standard library (e.g., database/sql/driver)

**Alternatives considered**:
- Group by type (config.go, slash.go) - Still would have 2 large files
- Keep monolithic file - Current state, harder to maintain

### Decision 3: Unified Provider Interface

**Chosen**: Single `Provider` interface for both config-based and slash-command providers

**Rationale**:
- Existing `Configurator` interface already works for both types
- Simpler executor logic (no type switching needed)
- Registry can store both types uniformly
- Future providers can mix behaviors without new interfaces

**Alternatives considered**:
- Separate interfaces (ConfigProvider, SlashProvider) - More complex, no clear benefit
- No interface (use concrete types) - Loses abstraction, harder to test/mock

### Decision 4: Registry-Based Lookup

**Chosen**: Global provider registry with `Register()` and `GetProvider()` functions

**Rationale**:
- Eliminates hardcoded switch statement in executor.go
- Providers self-register via `init()` - no central modification needed
- Enables future dynamic provider loading (plugins)
- Standard pattern used in many Go libraries (flag, database/sql)
- Thread-safe with sync.RWMutex (though init is single-threaded)

**Alternatives considered**:
- Keep switch statement - Current state, requires modification for each provider
- Factory function per provider - More boilerplate, no central registry
- Reflection-based discovery - Too magical, hard to debug

### Decision 5: Provider Registration in init()

**Chosen**: Each provider file has `init()` function calling `Register(id, provider)`

```go
// providers/claude.go
func init() {
    Register("claude-code", &ClaudeProvider{})
}
```

**Rationale**:
- Automatic registration when package imported
- No manual registration list to maintain
- Standard Go idiom (used by image formats, database drivers)
- Clear ownership (registration lives with implementation)

**Alternatives considered**:
- Manual registration list - Requires remembering to register new providers
- Constructor functions - More boilerplate, delayed registration

### Decision 6: Shared ProviderKit Package

**Chosen**: Create `internal/providerkit` package that owns the interface alias plus all shared utilities (marker helpers, template manager, filesystem helpers used by providers)

**Rationale**:
- Eliminates the import cycle: providers import `providerkit`, init imports registry, neither needs to import each other
- Keeps all shared helpers versioned together so providers and orchestrator stay in sync
- Allows TemplateManager and marker utilities to be consumed by providers without reaching into `internal/init`
- Keeps provider-specific package (`internal/providers`) lean; utilities live in a dedicated home

**Alternatives considered**:
- Keep helpers in `internal/init` - creates unavoidable cycle once executor imports providers
- Copy helpers into providers - duplication risk and inconsistent behavior
- Move helpers into providers - muddies package responsibilities and still leaves interface alias stranded

### Decision 7: Provider Metadata Owned by Registry

**Chosen**: Provider registry stores both the provider factory and user-facing metadata (name, priority, config/slash file paths, auto-install relationships)

**Rationale**:
- Wizard can render tool labels, help text, and ordering directly from the registry instead of duplicated hardcoded slices
- Adding a provider becomes a single operation: register implementation + metadata in one call
- Auto-install relationships (config → slash) live with provider definitions, so executor doesn’t need a separate mapping table
- Ensures CLI/UX requirements (“appears automatically in spectr init”) are actually achievable

**Alternatives considered**:
- Keep metadata in `internal/init/registry.go` - violates “no edits outside provider file” goal
- Rely on provider implementations to expose metadata ad-hoc - requires reflection/introspection and increases coupling
- Hardcode metadata in wizard/executor - same maintenance burden we’re trying to remove

### Decision 8: Backward-Compatible Configurator Interface

**Chosen**: Keep existing `Configurator` interface name and methods unchanged

```go
type Configurator interface {
    Configure(projectPath, spectrDir string) error
    IsConfigured(projectPath string) bool
    GetName() string
}
```

**Rationale**:
- Zero breaking changes for existing code
- wizard.go and other consumers work without modification
- Interface already well-designed for both provider types
- Renaming to `Provider` can be alias/type synonym

**Alternatives considered**:
- Rename to Provider interface - Requires updating all consumers
- Add new methods - No clear need, can extend later if needed

## Risks / Trade-offs

### Risk: ProviderKit Drift

**Risk**: Extracting utilities into `internal/providerkit` introduces another shared package that could diverge from init/providers expectations.

**Mitigation**:
- Keep ProviderKit surface area small (interface alias, marker utilities, template manager, slash base helpers)
- Add unit tests in both packages that assert shared behavior (e.g., marker updates, template rendering)
- Document ProviderKit contracts so contributor checklists include verifying changes across both call sites

### Risk: Test Complexity

**Risk**: Tests may become fragmented across 19+ files

**Mitigation**:
- Keep integration tests in `internal/init/` that test full flow
- Unit tests in `internal/providers/` test individual providers
- Table-driven tests can iterate over registered providers

### Risk: Provider Registration Order

**Risk**: `init()` execution order is undefined, potential race conditions

**Mitigation**:
- Registry uses sync.RWMutex for thread safety
- init is single-threaded, so no actual race in practice
- Duplicate registration returns error (detected in tests)

### Trade-off: More Files vs Navigability

**Trade-off**: 19+ files increases file count significantly

**Accepted because**:
- Modern editors handle many files well (fuzzy search, tree views)
- Benefit of focused, small files outweighs navigation cost
- Package organization (internal/providers/) groups them logically
- Alternative (1-2 large files) has worse maintainability

## Migration Plan

### Phase 1: Stand Up Shared Infrastructure
1. Create `internal/providerkit` containing the Configurator alias, marker utilities, template manager, and slash base implementation (no functional change yet)
2. Introduce provider registry in `internal/providers/registry.go` that stores metadata + factories but still backed by existing configurator types
3. Add tests that cover registry behavior, metadata validation, and ProviderKit helpers

### Phase 2: Pilot Extraction
4. Move Claude config + slash providers into `internal/providers` so they use ProviderKit and register metadata/factories
5. Wire executor/wizard to use the new registry path for Claude only, keeping the old switch for remaining tools as fallback
6. Run the full test suite to ensure hybrid path works before large migration

### Phase 3: Extract Remaining Providers
7. Move the rest of the config-based providers into their own files, ensuring each registers metadata + factories
8. Move slash providers and factories, deleting the duplicated constructors in `internal/init/configurator.go`
9. Keep running targeted tests after each move to catch regressions early

### Phase 4: Flip Execution Path
10. Remove the old `getConfigurator` switch in executor and rely solely on the registry response for both tool lists and configurators
11. Update wizard + registry consumers to rely on provider metadata (names, priorities, file paths) rather than the legacy `ToolRegistry`
12. Ensure auto-install relationships are modeled via metadata (`DependsOn`/`Installs`) instead of `configToSlashMapping`

### Phase 5: Cleanup
13. Delete the legacy configurator implementations from `internal/init/configurator.go`, leaving only ProviderKit imports if needed
14. Remove the obsolete ToolRegistry implementation and associated tests
15. Verify golangci-lint + go test + manual init smoke tests all pass

### Rollback Plan
- Keep git history clean with small, focused commits
- Each phase can be reverted independently
- If issues found post-merge, revert entire refactor (single change ID)
- ProviderKit maintains marker/template helpers, so providers can be moved back without touching orchestrator logic

## Open Questions

1. **Should we add Provider versioning?**
   - Not needed for v1, all providers bundled with CLI
   - Future: could add Version() string to interface for plugin system

2. **Should providers be able to declare dependencies?**
   - Not needed currently, all providers are independent
   - Future: could add DependsOn() []string if providers need ordering

3. **How to handle provider-specific configuration?**
   - Current design: providers are stateless, configuration passed via methods
   - Future: could add Config(opts map[string]any) if needed
