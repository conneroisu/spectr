package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Change with prefix",
			content:  "# Change: Add Feature\n\nMore content",
			expected: "Add Feature",
		},
		{
			name:     "Spec with prefix",
			content:  "# Spec: Authentication\n\nMore content",
			expected: "Authentication",
		},
		{
			name:     "No prefix",
			content:  "# CLI Framework\n\nMore content",
			expected: "CLI Framework",
		},
		{
			name:     "Multiple headings",
			content:  "# First Heading\n## Second Heading\n# Third Heading",
			expected: "First Heading",
		},
		{
			name:     "Extra whitespace",
			content:  "#   Change:   Trim Whitespace   \n\nMore content",
			expected: "Trim Whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "test.md")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			title, err := ExtractTitle(filePath)
			if err != nil {
				t.Fatalf("ExtractTitle failed: %v", err)
			}
			if title != tt.expected {
				t.Errorf("Expected title %q, got %q", tt.expected, title)
			}
		})
	}
}

func TestExtractTitle_NoHeading(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")
	content := "Some content without heading\n\nMore content"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	title, err := ExtractTitle(filePath)
	if err != nil {
		t.Fatalf("ExtractTitle failed: %v", err)
	}
	if title != "" {
		t.Errorf("Expected empty title, got %q", title)
	}
}

func TestCountTasks(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedTotal     int
		expectedCompleted int
	}{
		{
			name: "Mixed tasks",
			content: `## Tasks
- [ ] Task 1
- [x] Task 2
- [ ] Task 3
- [X] Task 4`,
			expectedTotal:     4,
			expectedCompleted: 2,
		},
		{
			name: "All completed",
			content: `## Tasks
- [x] Task 1
- [X] Task 2`,
			expectedTotal:     2,
			expectedCompleted: 2,
		},
		{
			name: "All incomplete",
			content: `## Tasks
- [ ] Task 1
- [ ] Task 2`,
			expectedTotal:     2,
			expectedCompleted: 0,
		},
		{
			name: "With indentation",
			content: `## Tasks
  - [ ] Indented task 1
    - [x] Nested task 2`,
			expectedTotal:     2,
			expectedCompleted: 1,
		},
		{
			name: "Mixed content",
			content: `## Tasks
Some text
- [ ] Task 1
More text
- [x] Task 2
Not a task line`,
			expectedTotal:     2,
			expectedCompleted: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "tasks.md")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			status, err := CountTasks(filePath)
			if err != nil {
				t.Fatalf("CountTasks failed: %v", err)
			}
			if status.Total != tt.expectedTotal {
				t.Errorf("Expected total %d, got %d", tt.expectedTotal, status.Total)
			}
			if status.Completed != tt.expectedCompleted {
				t.Errorf("Expected completed %d, got %d", tt.expectedCompleted, status.Completed)
			}
		})
	}
}

func TestCountTasks_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.md")

	status, err := CountTasks(filePath)
	if err != nil {
		t.Fatalf("CountTasks should not error on missing file: %v", err)
	}
	if status.Total != 0 || status.Completed != 0 {
		t.Errorf("Expected zero status, got total=%d, completed=%d", status.Total, status.Completed)
	}
}

func TestCountDeltas(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	specsDir := filepath.Join(changeDir, "specs", "test-spec")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Test Spec

## ADDED Requirements
### Requirement: New Feature

## MODIFIED Requirements
### Requirement: Updated Feature

## REMOVED Requirements
### Requirement: Old Feature
`
	specPath := filepath.Join(specsDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	count, err := CountDeltas(changeDir)
	if err != nil {
		t.Fatalf("CountDeltas failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 deltas, got %d", count)
	}
}

func TestCountDeltas_NoSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	count, err := CountDeltas(tmpDir)
	if err != nil {
		t.Fatalf("CountDeltas should not error on missing specs: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 deltas, got %d", count)
	}
}

func TestCountRequirements(t *testing.T) {
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")

	content := `# Test Spec

### Requirement: Feature 1
Description

### Requirement: Feature 2
Description

## Another Section

### Requirement: Feature 3
Description
`
	if err := os.WriteFile(specPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	count, err := CountRequirements(specPath)
	if err != nil {
		t.Fatalf("CountRequirements failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 requirements, got %d", count)
	}
}

func TestCountRequirements_NoRequirements(t *testing.T) {
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")

	content := `# Test Spec

Some content without requirements
`
	if err := os.WriteFile(specPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	count, err := CountRequirements(specPath)
	if err != nil {
		t.Fatalf("CountRequirements failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 requirements, got %d", count)
	}
}
