# Implementation Tasks

## 1. Remove ConfigPath from ToolDefinition

- [ ] 1.1 Remove `ConfigPath string` field from ToolDefinition struct in `internal/init/models.go`
- [ ] 1.2 Remove ConfigPath assignments from all tool registrations in `internal/init/registry.go` (6 tools)
- [ ] 1.3 Update `getToolFileInfo()` in `internal/init/executor.go` to return actual configurator file paths
- [ ] 1.4 Remove ConfigPath assertions from `internal/init/registry_test.go`
- [ ] 1.5 Run tests to verify no regressions: `go test ./internal/init/...`
- [ ] 1.6 Run full test suite: `go test ./...`
- [ ] 1.7 Verify the build compiles: `go build`
