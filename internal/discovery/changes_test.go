package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestGetActiveChanges(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create active changes
	testChanges := []string{"add-feature", "fix-bug", "update-docs"}
	for _, name := range testChanges {
		changeDir := filepath.Join(changesDir, name)
		if err := os.MkdirAll(changeDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		proposalPath := filepath.Join(changeDir, "proposal.md")
		if err := os.WriteFile(
			proposalPath,
			[]byte("# Test"),
			testFilePerm,
		); err != nil {
			t.Fatal(err)
		}
	}

	// Create archive directory (should be excluded)
	archiveDir := filepath.Join(changesDir, "archive", "old-change")
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(archiveDir, "proposal.md"),
		[]byte("# Old"),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}

	// Create hidden directory (should be excluded)
	hiddenDir := filepath.Join(changesDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(hiddenDir, "proposal.md"),
		[]byte("# Hidden"),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}

	// Create directory without proposal.md (should be excluded)
	emptyDir := filepath.Join(changesDir, "incomplete")
	if err := os.MkdirAll(emptyDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test discovery
	changes, err := GetActiveChanges(tmpDir)
	if err != nil {
		t.Fatalf("GetActiveChanges failed: %v", err)
	}

	if len(changes) != len(testChanges) {
		t.Errorf("Expected %d changes, got %d", len(testChanges), len(changes))
	}

	// Verify all expected changes are found
	changeMap := make(map[string]bool)
	for _, c := range changes {
		changeMap[c] = true
	}
	for _, expected := range testChanges {
		if !changeMap[expected] {
			t.Errorf("Expected change %s not found", expected)
		}
	}

	// Verify archived and hidden changes are not included
	if changeMap["old-change"] {
		t.Error("Archived change should not be included")
	}
	if changeMap[".hidden"] {
		t.Error("Hidden directory should not be included")
	}
	if changeMap["incomplete"] {
		t.Error("Incomplete change should not be included")
	}
}

func TestGetActiveChanges_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	changes, err := GetActiveChanges(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("Expected empty result, got %d changes", len(changes))
	}
}
func TestGetActiveChangeIDs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("Expected empty result, got %d changes", len(changes))
	}
}
func TestGetActiveChangeIDs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create changes in non-alphabetical order to test sorting
	testChanges := []string{
		"zebra-feature",
		"add-feature",
		"middle-feature",
	}
	for _, name := range testChanges {
		createChangeDir(t, changesDir, name, "# Test")
	}

	// Create archive directory (should be excluded)
	archiveDir := filepath.Join(changesDir, "archive")
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	createChangeDir(t, archiveDir, "archived-change", "# Archived")

	// Test GetActiveChangeIDs
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("GetActiveChangeIDs failed: %v", err)
	}

	if len(changes) != len(testChanges) {
		t.Errorf("Expected %d changes, got %d", len(testChanges), len(changes))
	}

	// Verify sorting
	// (should be: add-feature, middle-feature, zebra-feature)
	expectedSorted := []string{
		"add-feature",
		"middle-feature",
		"zebra-feature",
	}
	verifyOrdering(t, changes, expectedSorted, "change")

	// Verify archive is excluded
	verifyChangesExcluded(t, changes, []string{"archived-change"})
}
