//nolint:revive // file-length-limit - interactive functions logically grouped
package validation

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

const (
	// Table column widths for validation items
	validationIDWidth   = 35
	validationTypeWidth = 10
	validationPathWidth = 55

	// Menu dimensions
	menuWidth  = 50
	menuHeight = 8

	// Truncation settings
	validationPathTruncate = 53
	ellipsisMinLength      = 3

	// Table height
	tableHeight = 10
)

// menuModel represents the bubbletea model for the menu screen
type menuModel struct {
	choices     []string
	cursor      int
	selected    int
	quitting    bool
	projectPath string
	strict      bool
	jsonOutput  bool
}

// itemPickerModel represents the bubbletea model for item selection
type itemPickerModel struct {
	table       table.Model
	items       []ValidationItem
	selectedID  string
	quitting    bool
	projectPath string
	strict      bool
	jsonOutput  bool
}

// validationResultMsg is sent when validation completes
type validationResultMsg struct {
	results     []BulkResult
	hasFailures bool
	err         error
}

// Init initializes the menu model
func (m menuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the menu
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true

			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			m.selected = m.cursor

			return m.handleSelection()
		}

	case validationResultMsg:
		// Handle validation result
		m.quitting = true
		if msg.err != nil {
			fmt.Fprintf(os.Stderr, "Validation error: %v\n", msg.err)

			return m, tea.Quit
		}

		// Print results
		if m.jsonOutput {
			PrintBulkJSONResults(msg.results)
		} else {
			PrintBulkHumanResults(msg.results)
		}

		return m, tea.Quit
	}

	return m, nil
}

// handleSelection processes the menu selection
func (m menuModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.selected {
	case 0: // All (changes and specs)
		return m, m.validateAll()
	case 1: // All changes
		return m, m.validateChanges()
	case 2: // All specs
		return m, m.validateSpecs()
	case 3: // Pick specific item
		return m.showItemPicker()
	}

	return m, tea.Quit
}

// validateAll validates all items
func (m menuModel) validateAll() tea.Cmd {
	return func() tea.Msg {
		items, err := GetAllItems(m.projectPath)
		if err != nil {
			return validationResultMsg{err: err}
		}

		validator := NewValidator(m.strict)
		results, hasFailures := validateItems(validator, items)

		return validationResultMsg{
			results:     results,
			hasFailures: hasFailures,
		}
	}
}

// validateChanges validates all changes
func (m menuModel) validateChanges() tea.Cmd {
	return func() tea.Msg {
		items, err := GetChangeItems(m.projectPath)
		if err != nil {
			return validationResultMsg{err: err}
		}

		validator := NewValidator(m.strict)
		results, hasFailures := validateItems(validator, items)

		return validationResultMsg{
			results:     results,
			hasFailures: hasFailures,
		}
	}
}

// validateSpecs validates all specs
func (m menuModel) validateSpecs() tea.Cmd {
	return func() tea.Msg {
		items, err := GetSpecItems(m.projectPath)
		if err != nil {
			return validationResultMsg{err: err}
		}

		validator := NewValidator(m.strict)
		results, hasFailures := validateItems(validator, items)

		return validationResultMsg{
			results:     results,
			hasFailures: hasFailures,
		}
	}
}

// showItemPicker transitions to the item picker screen
func (m menuModel) showItemPicker() (tea.Model, tea.Cmd) {
	items, err := GetAllItems(m.projectPath)
	if err != nil {
		m.quitting = true
		fmt.Fprintf(os.Stderr, "Error loading items: %v\n", err)

		return m, tea.Quit
	}

	if len(items) == 0 {
		m.quitting = true
		fmt.Println("No items to validate")

		return m, tea.Quit
	}

	// Build table
	columns := []table.Column{
		{Title: "Name", Width: validationIDWidth},
		{Title: "Type", Width: validationTypeWidth},
		{Title: "Path", Width: validationPathWidth},
	}

	rows := make([]table.Row, len(items))
	for i, item := range items {
		rows[i] = table.Row{
			item.Name,
			item.ItemType,
			truncateString(item.Path, validationPathTruncate),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	applyTableStyles(&t)

	picker := itemPickerModel{
		table:       t,
		items:       items,
		projectPath: m.projectPath,
		strict:      m.strict,
		jsonOutput:  m.jsonOutput,
	}

	return picker, nil
}

// View renders the menu
func (m menuModel) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginBottom(1)

	choiceStyle := lipgloss.NewStyle().
		PaddingLeft(2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true).
		PaddingLeft(2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	s := titleStyle.Render("Validation Menu") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		if m.cursor == i {
			s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
		} else {
			s += choiceStyle.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("↑/↓ or j/k: navigate | Enter: select | q: quit")

	return s
}

// Init initializes the item picker model
func (m itemPickerModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the item picker
func (m itemPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			return m.handleSelection()
		}

	case validationResultMsg:
		// Handle validation result
		m.quitting = true
		if msg.err != nil {
			fmt.Fprintf(os.Stderr, "Validation error: %v\n", msg.err)

			return m, tea.Quit
		}

		// Print results
		if m.jsonOutput {
			PrintBulkJSONResults(msg.results)
		} else {
			PrintBulkHumanResults(msg.results)
		}

		return m, tea.Quit
	}

	// Update table with key events
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

// handleSelection validates the selected item
func (m itemPickerModel) handleSelection() (tea.Model, tea.Cmd) {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[cursor]
	m.selectedID = item.Name

	// Validate the selected item
	return m, func() tea.Msg {
		validator := NewValidator(m.strict)
		result, err := ValidateSingleItem(validator, item)

		return validationResultMsg{
			results:     []BulkResult{result},
			hasFailures: err != nil || !result.Valid,
			err:         err,
		}
	}
}

// View renders the item picker
func (m itemPickerModel) View() string {
	if m.quitting {
		return ""
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	view := m.table.View() + "\n"
	view += helpStyle.Render("↑/↓ or j/k: navigate | Enter: validate | q: back/quit") + "\n"

	return view
}

// RunInteractiveValidation runs the interactive validation TUI
func RunInteractiveValidation(projectPath string, strict bool, jsonOutput bool) error {
	// Check if running in a TTY
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return fmt.Errorf("interactive mode requires a TTY")
	}

	choices := []string{
		"All (changes and specs)",
		"All changes",
		"All specs",
		"Pick specific item",
	}

	m := menuModel{
		choices:     choices,
		cursor:      0,
		projectPath: projectPath,
		strict:      strict,
		jsonOutput:  jsonOutput,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive mode: %w", err)
	}

	// Check if we got a validation result
	if fm, ok := finalModel.(menuModel); ok {
		if fm.quitting {
			return nil
		}
	}

	return nil
}

// validateItems validates a list of items and returns results
func validateItems(
	validator *Validator,
	items []ValidationItem,
) ([]BulkResult, bool) {
	results := make([]BulkResult, 0, len(items))
	hasFailures := false

	for _, item := range items {
		result, err := ValidateSingleItem(validator, item)
		results = append(results, result)

		if err != nil || !result.Valid {
			hasFailures = true
		}
	}

	return results, hasFailures
}

// truncateString truncates a string and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= ellipsisMinLength {
		return s[:maxLen]
	}

	return s[:maxLen-ellipsisMinLength] + "..."
}

// applyTableStyles applies default styling to a table
func applyTableStyles(t *table.Model) {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("99"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	t.SetStyles(s)
}
