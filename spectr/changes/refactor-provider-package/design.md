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

### Decision 6: Marker Utilities Stay in init/configurator.go

**Chosen**: Keep `UpdateFileWithMarkers()` and related utilities in `internal/init/configurator.go`

**Rationale**:
- Used by multiple providers (not provider-specific logic)
- Part of the initialization infrastructure
- Moving would require providers to import complex helpers
- Clear separation: providers use utilities, don't implement them

**Alternatives considered**:
- Move to providers package - Would pollute provider package with utilities
- Move to separate util package - Over-engineering for ~150 lines
- Duplicate in providers - Violates DRY principle

### Decision 7: Backward-Compatible Configurator Interface

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

### Risk: Import Cycle

**Risk**: `internal/init` imports `internal/providers`, providers might need init utilities

**Mitigation**:
- Keep shared utilities (markers) in `internal/init/configurator.go`
- Providers import `internal/init` for utilities - one-way dependency
- If cycle occurs, extract shared utilities to `internal/init/shared` or similar

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

### Phase 1: Create New Package (No Breaking Changes)
1. Create `internal/providers/provider.go` with interface and registry
2. Extract one example provider (e.g., claude.go) and verify it works
3. All existing code still uses `internal/init/configurator.go`

### Phase 2: Extract All Providers
4. Extract remaining 18+ providers following claude.go pattern
5. Each provider self-registers in init()
6. Run tests after each extraction to catch issues early

### Phase 3: Update Executor
7. Replace executor.getConfigurator() switch with `providers.GetProvider()`
8. Update imports across init package
9. Run full test suite

### Phase 4: Cleanup
10. Remove extracted providers from configurator.go (now empty except utilities)
11. Update documentation and comments
12. Verify golangci-lint passes

### Rollback Plan
- Keep git history clean with small, focused commits
- Each phase can be reverted independently
- If issues found post-merge, revert entire refactor (single change ID)
- Marker utilities remaining in init/ means providers can be moved back easily

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
