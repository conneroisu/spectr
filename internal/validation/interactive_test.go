package validation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// TestMenuModel_Init tests menu model initialization
func TestMenuModel_Init(t *testing.T) {
	m := menuModel{
		choices:     []string{"Option 1", "Option 2"},
		cursor:      0,
		projectPath: "/test",
		strict:      false,
		jsonOutput:  false,
	}

	cmd := m.Init()
	assert.Zero(t, cmd)
}

// TestMenuModel_Update_Quit tests quitting the menu
func TestMenuModel_Update_Quit(t *testing.T) {
	m := menuModel{
		choices: []string{"Option 1", "Option 2"},
		cursor:  0,
	}

	tests := []struct {
		name string
		key  string
	}{
		{"quit with q", "q"},
		{"quit with ctrl+c", "ctrl+c"},
		{"quit with esc", "esc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "ctrl+c":
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			case "esc":
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			}

			newModel, cmd := m.Update(msg)
			updatedModel, ok := newModel.(menuModel)
			if !ok {
				t.Fatalf("expected menuModel, got %T", newModel)
			}

			assert.True(t, updatedModel.quitting)
			assert.NotZero(t, cmd)
		})
	}
}

// TestMenuModel_Update_Navigation tests menu navigation
func TestMenuModel_Update_Navigation(t *testing.T) {
	m := menuModel{
		choices: []string{"Option 1", "Option 2", "Option 3"},
		cursor:  1,
	}

	// Test moving up
	msg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := m.Update(msg)
	assert.Equal(t, 0, newModel.(menuModel).cursor)

	// Test moving down from cursor 0
	m.cursor = 0
	msg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.Update(msg)
	assert.Equal(t, 1, newModel.(menuModel).cursor)

	// Test k (vim up)
	m.cursor = 1
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	newModel, _ = m.Update(msg)
	assert.Equal(t, 0, newModel.(menuModel).cursor)

	// Test j (vim down) from cursor 0
	m.cursor = 0
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	newModel, _ = m.Update(msg)
	assert.Equal(t, 1, newModel.(menuModel).cursor)
}

// TestMenuModel_Update_NavigationBounds tests that cursor stays within bounds
func TestMenuModel_Update_NavigationBounds(t *testing.T) {
	m := menuModel{
		choices: []string{"Option 1", "Option 2"},
		cursor:  0,
	}

	// Try to move up from first item
	msg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := m.Update(msg)
	assert.Equal(t, 0, newModel.(menuModel).cursor)

	// Move to last item
	m.cursor = 1
	msg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.Update(msg)
	assert.Equal(t, 1, newModel.(menuModel).cursor) // Should stay at last
}

// TestMenuModel_View tests menu rendering
func TestMenuModel_View(t *testing.T) {
	m := menuModel{
		choices: []string{"All (changes and specs)", "All changes", "All specs"},
		cursor:  1,
	}

	view := m.View()

	assert.Contains(t, view, "Validation Menu")
	assert.Contains(t, view, "All (changes and specs)")
	assert.Contains(t, view, "All changes")
	assert.Contains(t, view, "All specs")
	assert.Contains(t, view, "↑/↓ or j/k: navigate")
	assert.Contains(t, view, "Enter: select")
	assert.Contains(t, view, "q: quit")
}

// TestMenuModel_View_Quitting tests that view returns empty when quitting
func TestMenuModel_View_Quitting(t *testing.T) {
	m := menuModel{
		choices:  []string{"Option 1"},
		quitting: true,
	}

	view := m.View()
	assert.Equal(t, "", view)
}

// TestMenuModel_View_Cursor tests cursor rendering
func TestMenuModel_View_Cursor(t *testing.T) {
	m := menuModel{
		choices: []string{"Option 1", "Option 2"},
		cursor:  0,
	}

	view := m.View()
	// Cursor should be on first item
	assert.Contains(t, view, ">")
}

// TestItemPickerModel_Init tests item picker initialization
func TestItemPickerModel_Init(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"test"}}
	tbl := table.New(table.WithColumns(columns), table.WithRows(rows))

	m := itemPickerModel{
		table:       tbl,
		items:       []ValidationItem{{Name: "test", ItemType: ItemTypeSpec}},
		projectPath: "/test",
		strict:      false,
	}

	cmd := m.Init()
	assert.Zero(t, cmd)
}

// TestItemPickerModel_Update_Quit tests quitting the item picker
func TestItemPickerModel_Update_Quit(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"test"}}
	tbl := table.New(table.WithColumns(columns), table.WithRows(rows))

	m := itemPickerModel{
		table: tbl,
		items: []ValidationItem{{Name: "test", ItemType: ItemTypeSpec}},
	}

	tests := []struct {
		name string
		key  string
	}{
		{"quit with q", "q"},
		{"quit with ctrl+c", "ctrl+c"},
		{"quit with esc", "esc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "ctrl+c":
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			case "esc":
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			}

			newModel, cmd := m.Update(msg)
			updatedModel, ok := newModel.(itemPickerModel)
			assert.True(t, ok)

			assert.True(t, updatedModel.quitting)
			assert.NotZero(t, cmd)
		})
	}
}

// TestItemPickerModel_View tests item picker rendering
func TestItemPickerModel_View(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 10},
	}
	rows := []table.Row{
		{"test-spec", "spec"},
		{"test-change", "change"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(10),
	)

	m := itemPickerModel{
		table: tbl,
		items: []ValidationItem{
			{Name: "test-spec", ItemType: ItemTypeSpec},
			{Name: "test-change", ItemType: ItemTypeChange},
		},
	}

	view := m.View()

	assert.Contains(t, view, "↑/↓ or j/k: navigate")
	assert.Contains(t, view, "Enter: validate")
	assert.Contains(t, view, "q: back/quit")
}

// TestItemPickerModel_View_Quitting tests view when quitting
func TestItemPickerModel_View_Quitting(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"test"}}
	tbl := table.New(table.WithColumns(columns), table.WithRows(rows))

	m := itemPickerModel{
		table:    tbl,
		quitting: true,
	}

	view := m.View()
	assert.Equal(t, "", view)
}

// TestRunInteractiveValidation_NotTTY tests error when not in a TTY
func TestRunInteractiveValidation_NotTTY(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe (not a TTY)
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	defer func() { _ = r.Close() }()
	defer func() { _ = w.Close() }()

	// Replace stdout with pipe
	os.Stdout = w

	// Restore stdout after test
	defer func() {
		os.Stdout = oldStdout
	}()

	// This should return an error because pipe is not a TTY
	err = RunInteractiveValidation("/test", false, false)

	// Restore stdout before assertions
	os.Stdout = oldStdout

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interactive mode requires a TTY")
}

// TestValidateItems tests the validateItems helper function
func TestValidateItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"test-spec"})
	createValidSpec(t, tmpDir, "test-spec")

	specPath := filepath.Join(tmpDir, SpectrDir, "specs", "test-spec", "spec.md")
	items := []ValidationItem{
		{
			Name:     "test-spec",
			ItemType: ItemTypeSpec,
			Path:     specPath,
		},
	}

	validator := NewValidator(false)
	results, hasFailures := validateItems(validator, items)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "test-spec", results[0].Name)
	assert.False(t, hasFailures)
	assert.True(t, results[0].Valid)
}

// TestValidateItems_MultipleItems tests validating multiple items
func TestValidateItems_MultipleItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"spec1", "spec2"})
	createValidSpec(t, tmpDir, "spec1")
	createValidSpec(t, tmpDir, "spec2")

	items := []ValidationItem{
		{
			Name:     "spec1",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "spec1", "spec.md"),
		},
		{
			Name:     "spec2",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "spec2", "spec.md"),
		},
	}

	validator := NewValidator(false)
	results, hasFailures := validateItems(validator, items)

	assert.Equal(t, 2, len(results))
	assert.False(t, hasFailures)
	assert.True(t, results[0].Valid)
	assert.True(t, results[1].Valid)
}

// TestValidateItems_WithFailures tests validation with some failures
func TestValidateItems_WithFailures(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"good-spec", "bad-spec"})
	createValidSpec(t, tmpDir, "good-spec")

	// Create invalid spec
	badSpecDir := filepath.Join(tmpDir, SpectrDir, "specs", "bad-spec")
	err := os.MkdirAll(badSpecDir, testDirPerm)
	assert.NoError(t, err)
	badSpecPath := filepath.Join(badSpecDir, "spec.md")
	err = os.WriteFile(badSpecPath, []byte("# Bad\nInvalid content"), testFilePerm)
	assert.NoError(t, err)

	items := []ValidationItem{
		{
			Name:     "good-spec",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "good-spec", "spec.md"),
		},
		{
			Name:     "bad-spec",
			ItemType: ItemTypeSpec,
			Path:     badSpecPath,
		},
	}

	validator := NewValidator(true) // strict mode
	results, hasFailures := validateItems(validator, items)

	assert.Equal(t, 2, len(results))
	assert.True(t, hasFailures)
	assert.True(t, results[0].Valid)  // good-spec should be valid
	assert.False(t, results[1].Valid) // bad-spec should be invalid
}

// TestValidateItems_EmptyList tests validation with empty item list
func TestValidateItems_EmptyList(t *testing.T) {
	validator := NewValidator(false)
	results, hasFailures := validateItems(validator, make([]ValidationItem, 0))

	assert.Equal(t, 0, len(results))
	assert.False(t, hasFailures)
}

// TestTruncateString tests the truncateString helper
func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "shorter than max",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "equal to max",
			input:    "exact",
			maxLen:   5,
			expected: "exact",
		},
		{
			name:     "longer than max",
			input:    "this is a very long string",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "maxLen less than ellipsis length",
			input:    "test",
			maxLen:   2,
			expected: "te",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   5,
			expected: "",
		},
		{
			name:     "path truncation",
			input:    "/very/long/path/to/some/file.md",
			maxLen:   20,
			expected: "/very/long/path/t...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
			assert.True(t, len(result) <= tt.maxLen)
		})
	}
}

// TestApplyTableStyles tests the applyTableStyles helper
func TestApplyTableStyles(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"test"}}
	tbl := table.New(table.WithColumns(columns), table.WithRows(rows))

	// This should not panic
	assert.NotPanics(t, func() {
		applyTableStyles(&tbl)
	})
}

// TestValidationResultMsg tests the validationResultMsg structure
func TestValidationResultMsg(t *testing.T) {
	msg := validationResultMsg{
		results: []BulkResult{
			{Name: "test", Type: ItemTypeSpec, Valid: true},
		},
		hasFailures: false,
		err:         nil,
	}

	assert.Equal(t, 1, len(msg.results))
	assert.False(t, msg.hasFailures)
	assert.Zero(t, msg.err)
}

// TestValidationResultMsg_WithError tests validation result with error
func TestValidationResultMsg_WithError(t *testing.T) {
	testErr := errors.New("test error")
	msg := validationResultMsg{
		results:     nil,
		hasFailures: true,
		err:         testErr,
	}

	assert.Zero(t, msg.results)
	assert.True(t, msg.hasFailures)
	assert.Equal(t, testErr, msg.err)
}

// TestMenuModel_Update_ValidationResult tests handling validation result message
func TestMenuModel_Update_ValidationResult(t *testing.T) {
	m := menuModel{
		choices:    []string{"Option 1"},
		jsonOutput: false,
	}

	// Create a validation result message
	msg := validationResultMsg{
		results: []BulkResult{
			{Name: "test", Type: ItemTypeSpec, Valid: true},
		},
		hasFailures: false,
		err:         nil,
	}

	newModel, cmd := m.Update(msg)
	updatedModel, ok := newModel.(menuModel)
	assert.True(t, ok)

	assert.True(t, updatedModel.quitting)
	assert.NotZero(t, cmd)
}

// TestMenuModel_Update_ValidationResultWithError tests handling validation error
func TestMenuModel_Update_ValidationResultWithError(t *testing.T) {
	m := menuModel{
		choices: []string{"Option 1"},
	}

	msg := validationResultMsg{
		results:     nil,
		hasFailures: true,
		err:         errors.New("test error"),
	}

	newModel, cmd := m.Update(msg)
	updatedModel, ok := newModel.(menuModel)
	assert.True(t, ok)

	assert.True(t, updatedModel.quitting)
	assert.NotZero(t, cmd)
}

// TestItemPickerModel_Update_ValidationResult tests handling validation in picker
func TestItemPickerModel_Update_ValidationResult(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"test"}}
	tbl := table.New(table.WithColumns(columns), table.WithRows(rows))

	m := itemPickerModel{
		table:      tbl,
		jsonOutput: false,
	}

	msg := validationResultMsg{
		results: []BulkResult{
			{Name: "test", Type: ItemTypeSpec, Valid: true},
		},
		hasFailures: false,
		err:         nil,
	}

	newModel, cmd := m.Update(msg)
	updatedModel, ok := newModel.(itemPickerModel)
	assert.True(t, ok)

	assert.True(t, updatedModel.quitting)
	assert.NotZero(t, cmd)
}

// TestMenuModel_HandleSelection_PickItem tests transitioning to item picker
func TestMenuModel_HandleSelection_PickItem(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"test-spec"})

	m := menuModel{
		choices:     []string{"All", "Changes", "Specs", "Pick specific item"},
		cursor:      3,
		selected:    3,
		projectPath: tmpDir,
		strict:      false,
		jsonOutput:  false,
	}

	newModel, _ := m.handleSelection()

	// Should transition to itemPickerModel
	_, ok := newModel.(itemPickerModel)
	assert.True(t, ok)
}

// TestMenuModel_HandleSelection_NoItems tests picking item with no items
func TestMenuModel_HandleSelection_NoItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, nil) // Empty project

	m := menuModel{
		choices:     []string{"All", "Changes", "Specs", "Pick specific item"},
		cursor:      3,
		selected:    3,
		projectPath: tmpDir,
		strict:      false,
		jsonOutput:  false,
	}

	newModel, cmd := m.handleSelection()

	// Should quit when no items
	updatedModel, ok := newModel.(menuModel)
	assert.True(t, ok)
	assert.True(t, updatedModel.quitting)
	assert.NotZero(t, cmd)
}

// TestConstants tests that constants are defined correctly
func TestConstants(t *testing.T) {
	assert.Equal(t, 35, validationIDWidth)
	assert.Equal(t, 10, validationTypeWidth)
	assert.Equal(t, 55, validationPathWidth)
	assert.Equal(t, 50, menuWidth)
	assert.Equal(t, 8, menuHeight)
	assert.Equal(t, 53, validationPathTruncate)
	assert.Equal(t, 3, ellipsisMinLength)
	assert.Equal(t, 10, tableHeight)
}

// TestMenuModel_View_AllChoices tests that all menu choices are displayed
func TestMenuModel_View_AllChoices(t *testing.T) {
	m := menuModel{
		choices: []string{
			"All (changes and specs)",
			"All changes",
			"All specs",
			"Pick specific item",
		},
		cursor: 0,
	}

	view := m.View()

	for _, choice := range m.choices {
		assert.Contains(t, view, choice)
	}
}

// TestTruncateString_ExactBoundary tests truncation at exact boundary
func TestTruncateString_ExactBoundary(t *testing.T) {
	// String exactly at truncation point
	input := strings.Repeat("a", validationPathTruncate)
	result := truncateString(input, validationPathTruncate)
	assert.Equal(t, input, result)

	// String one character over
	input = strings.Repeat("a", validationPathTruncate+1)
	result = truncateString(input, validationPathTruncate)
	assert.Equal(t, validationPathTruncate, len(result))
	assert.True(t, strings.HasSuffix(result, "..."))
}
