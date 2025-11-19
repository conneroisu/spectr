# Change: Refactor Provider Logic to Separate Package

## Why

The `internal/init/configurator.go` file has grown to 831 lines containing all 19+ provider implementations (Claude, Cline, Qwen, etc.). This creates maintainability challenges: adding new providers requires modifying multiple hardcoded locations (configurator.go, executor.go switch statement, registry.go), and finding specific provider logic is difficult in the monolithic file. Extracting providers to a separate package with registry-based lookup will improve code organization and make adding new AI tools trivial.

## What Changes

- Create new `internal/providers` package with provider interface and registry
- Add `internal/providerkit` package that hosts the Configurator alias, shared marker/template utilities, and base slash implementation to avoid import cycles
- Move tool metadata (friendly name, priority, config/slash files, auto-install relationships) into the provider registry so the wizard/executor discover providers dynamically
- Extract all 19+ provider implementations into individual files (one per provider)
- Unify config-based providers (CLAUDE.md creators) and slash-command providers under common interface
- Replace hardcoded `executor.getConfigurator()` switch statement with registry lookup
- Reduce `configurator.go` from 831 lines to ~150 lines (keeping only interface, markers, utilities)
- Add self-registration mechanism so providers register themselves via `init()` functions
- Update `registry.go` to support dynamic provider registration
- Maintain 100% backward compatibility - no CLI behavior changes

## Impact

- Affected specs: `cli-interface` (initialization wizard), `cli-framework` (new provider architecture)
- Affected code:
  - `internal/providerkit/` (new shared utilities + interface)
  - `internal/providers/` (19+ provider files plus registry metadata)
  - `internal/init/configurator.go` (major reduction, now consumes ProviderKit)
  - `internal/init/executor.go` (replace switch with registry lookup)
  - `internal/init/wizard.go` (read provider metadata instead of hardcoded registry)
- `internal/init/templates.go` and filesystem helpers move into ProviderKit but retain behavior
- All existing tests must pass with updated imports
