package init

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewInitExecutor(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Test creating executor with valid path
	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	if executor.projectPath != tmpDir {
		t.Errorf("Expected project path to be %s, got %s", tmpDir, executor.projectPath)
	}

	if executor.registry == nil {
		t.Error("Expected registry to be initialized")
	}

	if executor.tm == nil {
		t.Error("Expected template manager to be initialized")
	}
}

func TestNewInitExecutorInvalidPath(t *testing.T) {
	// Test with non-existent path
	_, err := NewInitExecutor("/this/path/does/not/exist/12345")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}

func TestExecuteBasicInitialization(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Execute with no tools selected
	result, err := executor.Execute(make([]string, 0))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify directory structure was created
	spectrDir := filepath.Join(tmpDir, "spectr")
	if !FileExists(spectrDir) {
		t.Error("Expected spectr/ directory to be created")
	}

	specsDir := filepath.Join(spectrDir, "specs")
	if !FileExists(specsDir) {
		t.Error("Expected spectr/specs/ directory to be created")
	}

	changesDir := filepath.Join(spectrDir, "changes")
	if !FileExists(changesDir) {
		t.Error("Expected spectr/changes/ directory to be created")
	}

	// Verify project.md was created
	projectFile := filepath.Join(spectrDir, "project.md")
	if !FileExists(projectFile) {
		t.Error("Expected spectr/project.md to be created")
	}

	// Verify AGENTS.md was created
	agentsFile := filepath.Join(spectrDir, "AGENTS.md")
	if !FileExists(agentsFile) {
		t.Error("Expected spectr/AGENTS.md to be created")
	}

	// Verify created files are tracked in result
	if len(result.CreatedFiles) == 0 {
		t.Error("Expected some created files to be tracked")
	}

	// Check that essential files are in the created list
	hasProjectMd := false
	hasAgentsMd := false
	for _, file := range result.CreatedFiles {
		if file == "spectr/project.md" {
			hasProjectMd = true
		}
		if file == "spectr/AGENTS.md" {
			hasAgentsMd = true
		}
	}

	if !hasProjectMd {
		t.Error("Expected spectr/project.md in created files")
	}
	if !hasAgentsMd {
		t.Error("Expected spectr/AGENTS.md in created files")
	}
}

func TestExecuteWithToolConfiguration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Execute with Claude Code tool selected
	result, err := executor.Execute([]string{"claude-code"})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify CLAUDE.md was created
	claudeFile := filepath.Join(tmpDir, "CLAUDE.md")
	if !FileExists(claudeFile) {
		t.Error("Expected CLAUDE.md to be created")
	}

	// Verify file is tracked in result
	hasClaudeMd := false
	for _, file := range result.CreatedFiles {
		if file == ".claude/claude.json" || file == "CLAUDE.md" {
			hasClaudeMd = true
		}
	}

	if !hasClaudeMd {
		t.Error("Expected CLAUDE.md in created files")
	}

	// Verify file content contains Spectr markers
	content, err := os.ReadFile(claudeFile)
	if err != nil {
		t.Fatalf("Failed to read CLAUDE.md: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, SpectrStartMarker) {
		t.Error("Expected CLAUDE.md to contain Spectr start marker")
	}
	if !contains(contentStr, SpectrEndMarker) {
		t.Error("Expected CLAUDE.md to contain Spectr end marker")
	}
}

func TestExecuteWithMultipleTools(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Execute with multiple tools
	tools := []string{"claude-code", "cline", "claude"}
	result, err := executor.Execute(tools)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify config-based tools were created
	claudeFile := filepath.Join(tmpDir, "CLAUDE.md")
	clineFile := filepath.Join(tmpDir, "CLINE.md")

	if !FileExists(claudeFile) {
		t.Error("Expected CLAUDE.md to be created")
	}
	if !FileExists(clineFile) {
		t.Error("Expected CLINE.md to be created")
	}

	// Verify slash command files were created
	slashFiles := []string{
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
		".claude/commands/spectr/archive.md",
	}

	for _, file := range slashFiles {
		fullPath := filepath.Join(tmpDir, file)
		if !FileExists(fullPath) {
			t.Errorf("Expected %s to be created", file)
		}
	}

	// Verify we have some created files tracked
	if len(result.CreatedFiles) < 5 {
		t.Errorf("Expected at least 5 created files, got %d", len(result.CreatedFiles))
	}
}

func TestExecuteIdempotency(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// First execution
	result1, err := executor.Execute(make([]string, 0))
	if err != nil {
		t.Fatalf("First execute failed: %v", err)
	}

	// Second execution (should detect existing Spectr)
	result2, err := executor.Execute(make([]string, 0))
	if err != nil {
		t.Fatalf("Second execute failed: %v", err)
	}

	// Should have warnings about existing files
	if len(result2.Errors) == 0 {
		t.Error("Expected warnings about existing Spectr on second execution")
	}

	// Verify first execution created files
	if len(result1.CreatedFiles) == 0 {
		t.Error("Expected created files in first execution")
	}

	// Second execution should have fewer created files (most already exist)
	if len(result2.CreatedFiles) >= len(result1.CreatedFiles) {
		t.Error("Expected fewer created files in second execution")
	}
}

func TestGetConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	testCases := []struct {
		toolID      string
		expectNil   bool
		description string
	}{
		{"claude-code", false, "Claude Code config"},
		{"cline", false, "Cline config"},
		{"claude", false, "Claude slash"},
		{"cursor", false, "Cursor slash"},
		{"invalid-tool", true, "Invalid tool ID"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			configurator := executor.getConfigurator(tc.toolID)
			if tc.expectNil && configurator != nil {
				t.Errorf("Expected nil configurator for %s", tc.toolID)
			}
			if !tc.expectNil && configurator == nil {
				t.Errorf("Expected non-nil configurator for %s", tc.toolID)
			}
		})
	}
}

func TestCreateDirectoryStructure(t *testing.T) {
	tmpDir := t.TempDir()
	executor, err := NewInitExecutor(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	result := &ExecutionResult{
		CreatedFiles: make([]string, 0),
		UpdatedFiles: make([]string, 0),
		Errors:       make([]string, 0),
	}

	spectrDir := filepath.Join(tmpDir, "spectr")
	err = executor.createDirectoryStructure(spectrDir, result)
	if err != nil {
		t.Fatalf("createDirectoryStructure failed: %v", err)
	}

	// Verify directories exist
	dirs := []string{
		spectrDir,
		filepath.Join(spectrDir, "specs"),
		filepath.Join(spectrDir, "changes"),
	}

	for _, dir := range dirs {
		if !FileExists(dir) {
			t.Errorf("Expected directory %s to exist", dir)
		}
	}

	// Verify result tracks created directories
	if len(result.CreatedFiles) != 3 {
		t.Errorf("Expected 3 created directories, got %d", len(result.CreatedFiles))
	}
}
