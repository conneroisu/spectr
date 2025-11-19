# Change: Consolidate Tool Configuration

## Why

The current init wizard presents 17+ AI tool options split between config-based tools (Claude Code, Cline, etc.) and slash command tools (Claude, Kilocode, etc.). This creates confusion:

1. **Unclear separation**: Users don't understand the difference between "Claude Code" and "Claude"
2. **Double selection burden**: To get full Claude integration, users must select TWO items
3. **UI complexity**: 17+ items in wizard is verbose and hard to navigate
4. **Inconsistent UX**: Some tools have both variants, some don't—no clear pattern

The OpenSpec project demonstrates the better pattern: when users select a tool, they automatically get all integration variants (config files + slash commands).

## What Changes

### Consolidation Strategy
- **Remove** 11+ standalone slash command tool registrations from wizard
- **Keep** 6 core config-based tools: Claude Code, Cline, CoStrict, Qoder, CodeBuddy, Qwen
- **Merge** slash command installation into config tool setup
- **Result**: When "Claude Code" is selected, users get BOTH `CLAUDE.md` + three slash commands in `.claude/commands/spectr/`

### Files Created Per Tool
When "Claude Code" is selected:
```
/ (project root)
└── CLAUDE.md                          ← NEW (config file with spectr markers)
└── .claude/
    └── commands/
        └── spectr/
            ├── proposal.md             ← NEW (slash command, auto-installed)
            ├── apply.md                ← NEW (slash command, auto-installed)
            └── archive.md              ← NEW (slash command, auto-installed)
```

### Key Changes by Component

**Tool Registry** (`internal/init/registry.go`)
- Add `SlashVariant` field to `ToolDefinition` to map tools to configurators
- Remove 11+ tool entries for slash-command-only tools
- Reduce public tool list from 17+ to 6 items

**Executor** (`internal/init/executor.go`)
- After configuring each tool, check if it has a slash variant
- Automatically call slash command configurator—no user prompt
- Track created/updated files from both config and slash configurators
- Update results display to show all generated files

**Wizard** (`internal/init/wizard.go`)
- Fewer items to display (6 instead of 17+)
- Same navigation and selection mechanics
- Help text updated to reflect tool count

**Spec Updates** (`spectr/specs/cli-interface/spec.md`)
- **MODIFIED**: "Flat Tool List in Initialization Wizard" requirement
  - Change from 17+ tools to 6 primary tools
  - Document automatic slash command installation
  - Update scenario wording

## Impact

### Affected Specifications
- `cli-interface`: Flat Tool List in Initialization Wizard requirement

### Breaking Changes
**BREAKING**: Slash command tool IDs (claude, kilocode, cline-slash, etc.) are removed from the tool registry and wizard. Code or scripts that reference these IDs will fail.

**Migration Path**:
- If users previously selected "Claude" separately, they should select "Claude Code" instead
- All other behavior is identical—same files are created in same locations

### Backward Compatibility
- Existing installations are unaffected (configurators check `IsConfigured()`)
- Running `spectr init` again does not duplicate files (markers + update logic)
- Generated slash command files are identical to previous implementation

### Benefits
1. **Simpler UX**: 6 tools instead of 17+
2. **Complete integration**: One selection = full tooling
3. **No confusion**: Clear tool names and automatic setup
4. **Consistency**: Matches OpenSpec pattern users may already use
5. **Maintainability**: Fewer tool definitions to manage

## Affected Code

### Files Modified
- `internal/init/registry.go` - Tool definitions
- `internal/init/executor.go` - Configuration logic
- `internal/init/wizard.go` - Display logic
- `spectr/specs/cli-interface/spec.md` - Spec deltas

### Internal Changes (No Public Impact)
- `internal/init/models.go` - Add SlashVariant field if needed
- `internal/init/configurator.go` - No changes (existing pattern works)
- `internal/init/templates.go` - No changes (templates already support all commands)

## Implementation Dependencies

1. ✓ Tool registry refactored with slash variant mapping
2. ✓ Executor updated to auto-install slash configurators
3. ✓ Spec deltas written
4. ✓ Tests updated to reflect new tool count
5. ✓ Validation passes strict mode

## Risks & Mitigations

| Risk | Severity | Mitigation |
|------|----------|-----------|
| Scripts expecting old tool IDs fail | Medium | Document in CHANGELOG, announce in release notes |
| Users confused by missing tools | Low | Explain consolidation in CLI help text |
| Slash command files already exist | Low | Existing configurator logic handles updates via markers |
| Spec delta conflicts | Low | Only cli-interface spec affected |

## Rollback Plan

If needed:
1. Revert registry changes to restore slash tool entries
2. Revert executor changes to skip auto-installation
3. Users can select tools independently as before
4. No data loss (files unchanged)
