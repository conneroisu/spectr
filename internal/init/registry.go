package init

import "fmt"

// ToolRegistry manages the collection of available AI tool definitions
type ToolRegistry struct {
	tools map[string]*ToolDefinition
}

// NewRegistry creates and initializes a new ToolRegistry with all
// 17 AI tool definitions
func NewRegistry() *ToolRegistry {
	registry := &ToolRegistry{
		tools: make(map[string]*ToolDefinition),
	}

	// Config-based tools (6 tools)
	registry.registerTool(&ToolDefinition{
		ID:         "claude-code",
		Name:       "Claude Code",
		Type:       ToolTypeConfig,
		ConfigPath: ".claude/claude.json",
		Priority:   1,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "cline",
		Name:       "Cline",
		Type:       ToolTypeConfig,
		ConfigPath: ".cline/cline_mcp_settings.json",
		Priority:   2,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "costrict-config",
		Name:       "Costrict",
		Type:       ToolTypeConfig,
		ConfigPath: ".costrict/config.json",
		Priority:   3,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "qoder-config",
		Name:       "Qoder",
		Type:       ToolTypeConfig,
		ConfigPath: ".qoder/config.json",
		Priority:   4,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "codebuddy",
		Name:       "CodeBuddy",
		Type:       ToolTypeConfig,
		ConfigPath: ".codebuddy/config.json",
		Priority:   5,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "qwen",
		Name:       "Qwen",
		Type:       ToolTypeConfig,
		ConfigPath: ".qwen/config.json",
		Priority:   6,
		Configured: false,
	})

	// Slash command tools (11 tools)
	registry.registerTool(&ToolDefinition{
		ID:           "claude",
		Name:         "Claude",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     7,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "kilocode",
		Name:         "Kilocode",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     8,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "qoder-slash",
		Name:         "Qoder (Slash)",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     9,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "cursor",
		Name:         "Cursor",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     10,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "aider",
		Name:         "Aider",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     11,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "continue",
		Name:         "Continue",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     12,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "copilot",
		Name:         "Copilot",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     13,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "mentat",
		Name:         "Mentat",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     14,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "tabnine",
		Name:         "Tabnine",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     15,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "smol",
		Name:         "Smol",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     16,
		Configured:   false,
	})

	registry.registerTool(&ToolDefinition{
		ID:           "costrict-slash",
		Name:         "Costrict (Slash)",
		Type:         ToolTypeSlash,
		SlashCommand: "/spectr",
		Priority:     17,
		Configured:   false,
	})

	return registry
}

// registerTool adds a tool to the registry
func (r *ToolRegistry) registerTool(tool *ToolDefinition) {
	r.tools[tool.ID] = tool
}

// GetTool retrieves a tool by its ID
// Returns an error if the tool ID is not found
func (r *ToolRegistry) GetTool(id string) (*ToolDefinition, error) {
	tool, exists := r.tools[id]
	if !exists {
		return nil, fmt.Errorf("tool with ID '%s' not found", id)
	}

	return tool, nil
}

// GetAllTools returns all registered tools as a slice
func (r *ToolRegistry) GetAllTools() []*ToolDefinition {
	tools := make([]*ToolDefinition, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// GetToolsByType returns all tools of a specific type
func (r *ToolRegistry) GetToolsByType(toolType ToolType) []*ToolDefinition {
	tools := make([]*ToolDefinition, 0)
	for _, tool := range r.tools {
		if tool.Type == toolType {
			tools = append(tools, tool)
		}
	}

	return tools
}

// ListTools returns a list of all tool IDs
func (r *ToolRegistry) ListTools() []string {
	ids := make([]string, 0, len(r.tools))
	for id := range r.tools {
		ids = append(ids, id)
	}

	return ids
}
