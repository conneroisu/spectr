# Implementation Tasks

## 1. Create Provider Package Foundation
- [ ] 1.1 Create `internal/providers/provider.go` with Provider interface
- [ ] 1.2 Implement provider registry (map[string]Provider)
- [ ] 1.3 Add Register() and GetProvider() functions
- [ ] 1.4 Add tests for registry (register, lookup, duplicate handling)

## 2. Extract Base Slash Command Provider
- [ ] 2.1 Create `internal/providers/base_slash.go`
- [ ] 2.2 Move SlashCommandConfigurator struct and methods
- [ ] 2.3 Update to implement Provider interface
- [ ] 2.4 Add tests for base slash provider

## 3. Extract Config-Based Providers (6 files)
- [ ] 3.1 Create `internal/providers/claude.go` with ClaudeCodeConfigurator
- [ ] 3.2 Create `internal/providers/cline.go` with ClineConfigurator
- [ ] 3.3 Create `internal/providers/qwen.go` with QwenConfigurator
- [ ] 3.4 Create `internal/providers/qoder.go` with QoderConfigurator
- [ ] 3.5 Create `internal/providers/codebuddy.go` with CodeBuddyConfigurator
- [ ] 3.6 Create `internal/providers/costrict.go` with CostrictConfigurator
- [ ] 3.7 Add init() functions to register each provider
- [ ] 3.8 Add tests for each config-based provider

## 4. Extract Slash Command Providers (15 files)
- [ ] 4.1 Create `internal/providers/claude_slash.go` with factory
- [ ] 4.2 Create `internal/providers/cline_slash.go` with factory
- [ ] 4.3 Create `internal/providers/kilocode_slash.go` with factory
- [ ] 4.4 Create `internal/providers/qoder_slash.go` with factory
- [ ] 4.5 Create `internal/providers/cursor_slash.go` with factory
- [ ] 4.6 Create `internal/providers/aider_slash.go` with factory
- [ ] 4.7 Create `internal/providers/continue_slash.go` with factory
- [ ] 4.8 Create `internal/providers/copilot_slash.go` with factory
- [ ] 4.9 Create `internal/providers/mentat_slash.go` with factory
- [ ] 4.10 Create `internal/providers/tabnine_slash.go` with factory
- [ ] 4.11 Create `internal/providers/smol_slash.go` with factory
- [ ] 4.12 Create `internal/providers/costrict_slash.go` with factory
- [ ] 4.13 Create `internal/providers/windsurf_slash.go` with factory
- [ ] 4.14 Create `internal/providers/codebuddy_slash.go` with factory
- [ ] 4.15 Create `internal/providers/qwen_slash.go` with factory
- [ ] 4.16 Add init() functions to register all slash providers
- [ ] 4.17 Add tests for slash provider factories

## 5. Refactor Internal/Init Package
- [ ] 5.1 Update `configurator.go` - remove all concrete provider types
- [ ] 5.2 Keep Configurator interface in `configurator.go`
- [ ] 5.3 Keep marker utilities (UpdateFileWithMarkers, etc.) in `configurator.go`
- [ ] 5.4 Update `executor.go` - replace getConfigurator() switch with registry lookup
- [ ] 5.5 Update `registry.go` - add provider registration support
- [ ] 5.6 Update imports in `wizard.go` if needed
- [ ] 5.7 Update all test files with new import paths

## 6. Validation and Testing
- [ ] 6.1 Run all existing unit tests - verify 100% pass
- [ ] 6.2 Test `spectr init` interactive mode with each tool selection
- [ ] 6.3 Test config-based tools create correct markdown files
- [ ] 6.4 Test slash-command tools create correct .claude/commands/ files
- [ ] 6.5 Verify marker-based file updates still work correctly
- [ ] 6.6 Run `go build` and verify no compilation errors
- [ ] 6.7 Run golangci-lint and fix any issues

## 7. Documentation and Cleanup
- [ ] 7.1 Add package documentation to `providers/provider.go`
- [ ] 7.2 Add godoc comments for all exported types
- [ ] 7.3 Update any internal documentation referencing old structure
- [ ] 7.4 Verify code coverage maintained or improved
