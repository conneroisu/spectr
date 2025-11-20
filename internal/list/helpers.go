package list

import (
	"encoding/base64"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// ClipboardWriter is an interface for writing to clipboard
type ClipboardWriter interface {
	WriteAll(text string) error
}

// realClipboardWriter implements ClipboardWriter using the actual clipboard library
type realClipboardWriter struct{}

// WriteAll writes text to the system clipboard
func (r realClipboardWriter) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// mockClipboardWriter implements ClipboardWriter for testing
type mockClipboardWriter struct {
	CopiedText string
}

// WriteAll stores text in memory instead of writing to clipboard
func (m *mockClipboardWriter) WriteAll(text string) error {
	m.CopiedText = text
	return nil
}

// copyToClipboard copies text to clipboard using native or OSC 52
func copyToClipboard(text string, writer ClipboardWriter) error {
	// Try native clipboard first
	err := writer.WriteAll(text)
	if err == nil {
		return nil
	}

	// Fallback to OSC 52 for SSH sessions
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	osc52 := "\x1b]52;c;" + encoded + "\x07"
	fmt.Print(osc52)

	// OSC 52 doesn't report errors, consider it successful
	return nil
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
