package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GetActiveChanges finds all active changes in spectr/changes/,
// excluding archive directory
func GetActiveChanges(projectPath string) ([]string, error) {
	changesDir := filepath.Join(projectPath, "spectr", "changes")

	// Check if changes directory exists
	if _, err := os.Stat(changesDir); os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read changes directory: %w", err)
	}

	var changes []string
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip archive directory
		if entry.Name() == "archive" {
			continue
		}

		// Check if proposal.md exists
		proposalPath := filepath.Join(changesDir, entry.Name(), "proposal.md")
		if _, err := os.Stat(proposalPath); err == nil {
			changes = append(changes, entry.Name())
		}
	}

	// Sort alphabetically for consistency
	sort.Strings(changes)

	return changes, nil
}

// GetActiveChangeIDs returns a list of active change IDs
// (directory names under spectr/changes/, excluding archive/)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetActiveChangeIDs(projectRoot string) ([]string, error) {
	return GetActiveChanges(projectRoot)
}
