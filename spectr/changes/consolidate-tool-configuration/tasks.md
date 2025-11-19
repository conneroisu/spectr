# Implementation Tasks: Consolidate Tool Configuration

## Phase 1: Analyze Current State
- [ ] 1.1 Document all current tool IDs and their types (config vs slash)
- [ ] 1.2 Identify which tools have both config and slash variants
- [ ] 1.3 Verify all slash command templates are identical across tools
- [ ] 1.4 Check for any tool-specific customization in templates (should be generic)

## Phase 2: Update Tool Registry
- [ ] 2.1 Add `SlashVariant` field to `ToolDefinition` struct in `models.go`
- [ ] 2.2 Update each config tool entry to include slash variant reference
  - [ ] 2.2a Claude Code → claude
  - [ ] 2.2b Cline → cline
  - [ ] 2.2c CoStrict → costrukt
  - [ ] 2.2d Qoder → qoder
  - [ ] 2.2e CodeBuddy → codebuddy
  - [ ] 2.2f Qwen → qwen
- [ ] 2.3 Remove slash-only tool registrations from Tools registry
- [ ] 2.4 Add helper method `GetSlashVariantConfigurator(toolID)` to registry

## Phase 3: Update Executor Logic
- [ ] 3.1 Modify `configureTools()` in executor.go
  - [ ] 3.1a After tool configuration, check for slash variant
  - [ ] 3.1b Invoke slash configurator if variant exists
  - [ ] 3.1c Track created/updated files from both configurators
  - [ ] 3.1d Update results to include all created files
- [ ] 3.2 Add error handling for slash variant configuration failures
- [ ] 3.3 Test error recovery (tool config succeeds, slash config fails)

## Phase 4: Update Wizard Display
- [ ] 4.1 Verify wizard still displays 6 tools (no changes needed to navigation)
- [ ] 4.2 Update help text to reflect tool consolidation
- [ ] 4.3 Test keyboard navigation with reduced tool count
- [ ] 4.4 Verify selection/deselection works correctly

## Phase 5: Write Spec Deltas
- [ ] 5.1 Create spec delta for cli-interface
  - [ ] 5.1a MODIFIED: Flat Tool List in Initialization Wizard
    - Update from "17 tools" to "6 primary tools"
    - Document automatic slash command installation
    - Update all scenario descriptions
    - Update help text reference

## Phase 6: Update Tests
- [ ] 6.1 Registry tests
  - [ ] 6.1a Update tool count assertions (17+ → 6)
  - [ ] 6.1b Add tests for SlashVariant field
  - [ ] 6.1c Add tests for helper method `GetSlashVariantConfigurator()`
- [ ] 6.2 Executor tests
  - [ ] 6.2a Test auto-installation of slash variants
  - [ ] 6.2b Test file tracking from both configurators
  - [ ] 6.2c Test error handling when slash config fails
- [ ] 6.3 Wizard tests
  - [ ] 6.3a Verify 6 tools displayed
  - [ ] 6.3b Test navigation and selection unchanged
- [ ] 6.4 Integration tests
  - [ ] 6.4a Run `spectr init` with each tool, verify all files created
  - [ ] 6.4b Verify CLAUDE.md + slash commands for each tool

## Phase 7: Validation & Documentation
- [ ] 7.1 Run `spectr validate consolidate-tool-configuration --strict`
- [ ] 7.2 Fix any validation errors
- [ ] 7.3 Update CHANGELOG with breaking change notice
- [ ] 7.4 Add migration guidance in CLI help/docs

## Phase 8: Manual Testing
- [ ] 8.1 Remove existing spectr/ and .claude/ directories
- [ ] 8.2 Run `spectr init` in clean environment
- [ ] 8.3 Verify tool list shows 6 items (not 17+)
- [ ] 8.4 Select "Claude Code"
- [ ] 8.5 Verify creation of:
  - [ ] 8.5a CLAUDE.md at project root
  - [ ] 8.5b .claude/commands/spectr/proposal.md
  - [ ] 8.5c .claude/commands/spectr/apply.md
  - [ ] 8.5d .claude/commands/spectr/archive.md
- [ ] 8.6 Verify file contents are correct
- [ ] 8.7 Verify spectr/ directory structure created
- [ ] 8.8 Test selecting different tools (Cline, etc.)

## Phase 9: Edge Cases & Robustness
- [ ] 9.1 Test re-running `spectr init` (existing files should update, not duplicate)
- [ ] 9.2 Test with only some tools selected
- [ ] 9.3 Test with all tools selected
- [ ] 9.4 Test with no tools selected (should still initialize spectr/)
- [ ] 9.5 Verify markers work correctly in updated files

## Completion Criteria
✓ All 9 phases completed
✓ All tests passing (unit + integration)
✓ `spectr validate consolidate-tool-configuration --strict` passes
✓ Manual testing on Linux, macOS, Windows (if applicable)
✓ Documentation updated
✓ No regression in other init functionality
✓ Slash command files identical to previous implementation
