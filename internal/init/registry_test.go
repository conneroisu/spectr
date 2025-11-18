package init

import (
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	// Test that the registry is not nil
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	// Test that all 17 tools are registered
	allTools := registry.GetAllTools()
	if len(allTools) != 17 {
		t.Errorf("Expected 17 tools, got %d", len(allTools))
	}

	// Test that the tools map is not nil
	if registry.tools == nil {
		t.Error("registry.tools map is nil")
	}
}

func TestGetTool(t *testing.T) {
	registry := NewRegistry()

	tests := []struct {
		name    string
		toolID  string
		wantErr bool
	}{
		{"Get Claude Code", "claude-code", false},
		{"Get Cline", "cline", false},
		{"Get Costrict Config", "costrict-config", false},
		{"Get Qoder Config", "qoder-config", false},
		{"Get CodeBuddy", "codebuddy", false},
		{"Get Qwen", "qwen", false},
		{"Get Claude", "claude", false},
		{"Get Kilocode", "kilocode", false},
		{"Get Qoder Slash", "qoder-slash", false},
		{"Get Cursor", "cursor", false},
		{"Get Aider", "aider", false},
		{"Get Continue", "continue", false},
		{"Get Copilot", "copilot", false},
		{"Get Mentat", "mentat", false},
		{"Get Tabnine", "tabnine", false},
		{"Get Smol", "smol", false},
		{"Get Costrict Slash", "costrict-slash", false},
		{"Get Invalid Tool", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := registry.GetTool(tt.toolID)
			if tt.wantErr {
				verifyErrorCase(t, tt.toolID, tool, err)
			} else {
				verifySuccessCase(t, tt.toolID, tool, err)
			}
		})
	}
}

// verifyErrorCase checks that GetTool properly returns an error
func verifyErrorCase(
	t *testing.T,
	toolID string,
	tool *ToolDefinition,
	err error,
) {
	t.Helper()
	if err == nil {
		t.Errorf("GetTool(%s) expected error, got nil", toolID)
	}
	if tool != nil {
		t.Errorf(
			"GetTool(%s) expected nil tool, got %v",
			toolID,
			tool,
		)
	}
}

// verifySuccessCase checks that GetTool returns a valid tool
func verifySuccessCase(
	t *testing.T,
	toolID string,
	tool *ToolDefinition,
	err error,
) {
	t.Helper()
	if err != nil {
		t.Errorf("GetTool(%s) unexpected error: %v", toolID, err)
	}
	if tool == nil {
		t.Errorf("GetTool(%s) returned nil tool", toolID)

		return
	}
	if tool.ID != toolID {
		t.Errorf(
			"GetTool(%s) returned tool with ID %s",
			toolID,
			tool.ID,
		)
	}
}

func TestGetToolsByType(t *testing.T) {
	registry := NewRegistry()

	// Test config-based tools
	configTools := registry.GetToolsByType(ToolTypeConfig)
	if len(configTools) != 6 {
		t.Errorf("Expected 6 config tools, got %d", len(configTools))
	}

	// Verify all config tools have ConfigPath set
	for _, tool := range configTools {
		if tool.ConfigPath == "" {
			t.Errorf("Config tool %s has empty ConfigPath", tool.ID)
		}
		if tool.Type != ToolTypeConfig {
			t.Errorf("Config tool %s has wrong type: %s", tool.ID, tool.Type)
		}
	}

	// Test slash command tools
	slashTools := registry.GetToolsByType(ToolTypeSlash)
	if len(slashTools) != 11 {
		t.Errorf("Expected 11 slash tools, got %d", len(slashTools))
	}

	// Verify all slash tools have SlashCommand set
	for _, tool := range slashTools {
		if tool.SlashCommand == "" {
			t.Errorf("Slash tool %s has empty SlashCommand", tool.ID)
		}
		if tool.Type != ToolTypeSlash {
			t.Errorf("Slash tool %s has wrong type: %s", tool.ID, tool.Type)
		}
	}
}

func TestConfigToolsHaveConfigPath(t *testing.T) {
	registry := NewRegistry()

	expectedConfigTools := map[string]string{
		"claude-code":     ".claude/claude.json",
		"cline":           ".cline/cline_mcp_settings.json",
		"costrict-config": ".costrict/config.json",
		"qoder-config":    ".qoder/config.json",
		"codebuddy":       ".codebuddy/config.json",
		"qwen":            ".qwen/config.json",
	}

	for id, expectedPath := range expectedConfigTools {
		tool, err := registry.GetTool(id)
		if err != nil {
			t.Errorf("Tool %s not found: %v", id, err)

			continue
		}

		if tool.ConfigPath != expectedPath {
			t.Errorf("Tool %s has ConfigPath %s, expected %s", id, tool.ConfigPath, expectedPath)
		}

		if tool.Type != ToolTypeConfig {
			t.Errorf("Tool %s has Type %s, expected %s", id, tool.Type, ToolTypeConfig)
		}

		if tool.SlashCommand != "" {
			t.Errorf("Config tool %s should not have SlashCommand set", id)
		}
	}
}

func TestSlashToolsHaveSlashCommand(t *testing.T) {
	registry := NewRegistry()

	expectedSlashTools := []string{
		"claude",
		"kilocode",
		"qoder-slash",
		"cursor",
		"aider",
		"continue",
		"copilot",
		"mentat",
		"tabnine",
		"smol",
		"costrict-slash",
	}

	for _, id := range expectedSlashTools {
		tool, err := registry.GetTool(id)
		if err != nil {
			t.Errorf("Tool %s not found: %v", id, err)

			continue
		}

		if tool.SlashCommand != "/spectr" {
			t.Errorf("Tool %s has SlashCommand %s, expected /spectr", id, tool.SlashCommand)
		}

		if tool.Type != ToolTypeSlash {
			t.Errorf("Tool %s has Type %s, expected %s", id, tool.Type, ToolTypeSlash)
		}

		if tool.ConfigPath != "" {
			t.Errorf("Slash tool %s should not have ConfigPath set", id)
		}
	}
}

func TestListTools(t *testing.T) {
	registry := NewRegistry()

	toolIDs := registry.ListTools()

	// Test that we get 17 tool IDs
	if len(toolIDs) != 17 {
		t.Errorf("Expected 17 tool IDs, got %d", len(toolIDs))
	}

	// Test that all tool IDs are unique
	seen := make(map[string]bool)
	for _, id := range toolIDs {
		if seen[id] {
			t.Errorf("Duplicate tool ID found: %s", id)
		}
		seen[id] = true
	}
}

func TestToolIDsAreKebabCase(t *testing.T) {
	registry := NewRegistry()

	allTools := registry.GetAllTools()

	for _, tool := range allTools {
		// Check that ID contains only lowercase letters, numbers, and hyphens
		for _, char := range tool.ID {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				t.Errorf(
					"Tool ID %s is not in kebab-case (contains invalid character: %c)",
					tool.ID,
					char,
				)
			}
		}
	}
}

func TestAllToolsHaveRequiredFields(t *testing.T) {
	registry := NewRegistry()

	allTools := registry.GetAllTools()

	for _, tool := range allTools {
		if tool.ID == "" {
			t.Error("Found tool with empty ID")
		}
		if tool.Name == "" {
			t.Errorf("Tool %s has empty Name", tool.ID)
		}
		if tool.Type != ToolTypeConfig && tool.Type != ToolTypeSlash {
			t.Errorf("Tool %s has invalid Type: %s", tool.ID, tool.Type)
		}
		if tool.Priority < 1 || tool.Priority > 17 {
			t.Errorf("Tool %s has invalid Priority: %d (should be 1-17)", tool.ID, tool.Priority)
		}
		if tool.Configured {
			t.Errorf("Tool %s should start with Configured=false", tool.ID)
		}
	}
}

func TestPrioritiesAreUnique(t *testing.T) {
	registry := NewRegistry()

	allTools := registry.GetAllTools()
	priorities := make(map[int]string)

	for _, tool := range allTools {
		if existingTool, exists := priorities[tool.Priority]; exists {
			t.Errorf(
				"Duplicate priority %d found for tools %s and %s",
				tool.Priority,
				existingTool,
				tool.ID,
			)
		}
		priorities[tool.Priority] = tool.ID
	}
}
