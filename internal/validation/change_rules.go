// Package validation provides validation rules for Spectr changes and specs.
//
//nolint:revive // file-length-limit,argument-limit,line-length-limit - validation logic requires comprehensive parameters
package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/conneroisu/spectr/internal/parsers"
)

// DeltaType represents the type of delta operation
type DeltaType string

const (
	DeltaAdded    DeltaType = "ADDED"
	DeltaModified DeltaType = "MODIFIED"
	DeltaRemoved  DeltaType = "REMOVED"
	DeltaRenamed  DeltaType = "RENAMED"
)

// ValidateChangeDeltaSpecs validates all delta spec files in a
// change directory. changeDir should be the path to a change directory
// (e.g., spectr/changes/add-feature). spectrRoot should be the path to the
// spectr/ directory (e.g., /path/to/project/spectr).
// Returns ValidationReport with all issues found, or error for issues
//
//nolint:revive // strictMode is intentional control flag
func ValidateChangeDeltaSpecs(
	changeDir string,
	spectrRoot string,
	strictMode bool,
) (*ValidationReport, error) {
	specsDir := filepath.Join(changeDir, "specs")

	// Check if specs directory exists
	info, err := os.Stat(specsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("specs directory not found: %s", specsDir)
		}

		return nil, fmt.Errorf("failed to access specs directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("specs path is not a directory: %s", specsDir)
	}

	// Find all spec.md files under specs/
	var specFiles []string
	err = filepath.Walk(specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "spec.md" {
			specFiles = append(specFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk specs directory: %w", err)
	}

	if len(specFiles) == 0 {
		return nil, fmt.Errorf("no spec.md files found in specs directory: %s", specsDir)
	}

	// Track all issues across all spec files
	var allIssues []ValidationIssue

	// Track requirement names for duplicate/conflict detection across all files
	addedReqs := make(map[string]string)       // normalized name -> file path
	modifiedReqs := make(map[string]string)    // normalized name -> file path
	removedReqs := make(map[string]string)     // normalized name -> file path
	renamedFromReqs := make(map[string]string) // normalized FROM name -> file path
	renamedToReqs := make(map[string]string)   // normalized TO name -> file path

	// Track total delta count
	totalDeltas := 0

	// Process each spec file
	for _, specPath := range specFiles {
		fileIssues, deltaCount, err := validateSingleDeltaFile(
			specPath,
			addedReqs,
			modifiedReqs,
			removedReqs,
			renamedFromReqs,
			renamedToReqs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to validate %s: %w", specPath, err)
		}

		allIssues = append(allIssues, fileIssues...)
		totalDeltas += deltaCount

		// Validate delta file against base spec
		baseSpecIssues, err := validateDeltaAgainstBaseSpec(specPath, specsDir, spectrRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to validate %s against base spec: %w", specPath, err)
		}
		allIssues = append(allIssues, baseSpecIssues...)
	}

	// Check if there are no deltas at all
	if totalDeltas == 0 {
		allIssues = append(allIssues, ValidationIssue{
			Level: LevelError,
			Path:  specsDir,
			Message: "Change must have at least one delta " +
				"(ADDED, MODIFIED, REMOVED, or RENAMED requirement)",
		})
	}

	// Apply strict mode: convert warnings to errors
	if strictMode {
		for i := range allIssues {
			if allIssues[i].Level == LevelWarning {
				allIssues[i].Level = LevelError
			}
		}
	}

	return NewValidationReport(allIssues), nil
}

// validateSingleDeltaFile validates a single spec.md delta file
// Returns issues, delta count, and error
func validateSingleDeltaFile(
	specPath string,
	addedReqs, modifiedReqs, removedReqs, renamedFromReqs, renamedToReqs map[string]string,
) ([]ValidationIssue, int, error) {
	// Read file
	content, err := os.ReadFile(specPath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse sections
	sections := ExtractSections(string(content))
	var issues []ValidationIssue
	deltaCount := 0

	// Track requirement names within this file for duplicate detection
	fileAddedReqs := make(map[string]bool)
	fileModifiedReqs := make(map[string]bool)
	fileRemovedReqs := make(map[string]bool)
	fileRenamedFromReqs := make(map[string]bool)
	fileRenamedToReqs := make(map[string]bool)

	// Process ADDED Requirements
	if addedContent, hasAdded := sections["ADDED Requirements"]; hasAdded {
		deltaCount++
		addedIssues := validateAddedRequirements(
			addedContent,
			specPath,
			fileAddedReqs,
			addedReqs,
		)
		issues = append(issues, addedIssues...)
	}

	// Process MODIFIED Requirements
	if modifiedContent, hasModified := sections["MODIFIED Requirements"]; hasModified {
		deltaCount++
		modifiedIssues := validateModifiedRequirements(
			modifiedContent,
			specPath,
			fileModifiedReqs,
			modifiedReqs,
		)
		issues = append(issues, modifiedIssues...)
	}

	// Process REMOVED Requirements
	if removedContent, hasRemoved := sections["REMOVED Requirements"]; hasRemoved {
		deltaCount++
		removedIssues := validateRemovedRequirements(
			removedContent,
			specPath,
			fileRemovedReqs,
			removedReqs,
		)
		issues = append(issues, removedIssues...)
	}

	// Process RENAMED Requirements
	if renamedContent, hasRenamed := sections["RENAMED Requirements"]; hasRenamed {
		deltaCount++
		renamedIssues := validateRenamedRequirements(
			renamedContent,
			specPath,
			fileRenamedFromReqs,
			fileRenamedToReqs,
			renamedFromReqs,
			renamedToReqs,
		)
		issues = append(issues, renamedIssues...)
	}

	// Check for cross-section conflicts within this file
	for normalized := range fileAddedReqs {
		if fileModifiedReqs[normalized] {
			// Find actual requirement name for better error message
			reqName := findRequirementNameByNormalized(
				normalized,
				addedReqs,
				specPath,
			)
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  specPath,
				Message: fmt.Sprintf(
					"Requirement '%s' appears in both ADDED and "+
						"MODIFIED sections",
					reqName,
				),
			})
		}
	}

	return issues, deltaCount, nil
}

// RenamedRequirement represents a FROM -> TO rename pair
type RenamedRequirement struct {
	FromName string
	ToName   string
}

// parseRenamedRequirements parses the RENAMED Requirements section
// Expected format:
// - FROM: ### Requirement: OldName
// - TO: ### Requirement: NewName
func parseRenamedRequirements(content string) []RenamedRequirement {
	var renames []RenamedRequirement
	lines := strings.Split(content, "\n")

	fromRegex := regexp.MustCompile(
		`^\s*-\s*FROM:\s*###\s*Requirement:\s*(.+)$`,
	)
	toRegex := regexp.MustCompile(
		`^\s*-\s*TO:\s*###\s*Requirement:\s*(.+)$`,
	)

	var currentFrom string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for FROM line
		if matches := fromRegex.FindStringSubmatch(line); matches != nil {
			currentFrom = strings.TrimSpace(matches[1])

			continue
		}

		// Check for TO line
		matches := toRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		toName := strings.TrimSpace(matches[1])

		// If we have a FROM, pair it with this TO
		if currentFrom != "" {
			renames = append(renames, RenamedRequirement{
				FromName: currentFrom,
				ToName:   toName,
			})
			currentFrom = "" // Reset for next pair
		} else {
			// TO without FROM - add as malformed
			renames = append(renames, RenamedRequirement{
				FromName: "",
				ToName:   toName,
			})
		}
	}

	// If we have a FROM without a TO, add as malformed
	if currentFrom != "" {
		renames = append(renames, RenamedRequirement{
			FromName: currentFrom,
			ToName:   "",
		})
	}

	return renames
}

// findRequirementNameByNormalized finds the original requirement name
// given a normalized version
func findRequirementNameByNormalized(
	normalized string,
	_ map[string]string,
	_ string,
) string {
	// This is a helper to provide better error messages
	// In practice, we'd need to iterate through requirements again
	// For now, return the normalized version as it's close enough
	return normalized
}

// validateDeltaAgainstBaseSpec validates a delta file against the base spec
// Returns validation issues for pre-merge validation errors
func validateDeltaAgainstBaseSpec(
	deltaSpecPath string,
	specsDir string,
	spectrRoot string,
) ([]ValidationIssue, error) {
	// Extract capability name from delta spec path
	// Path structure: .../changes/<change-id>/specs/<capability>/spec.md
	// We want to extract <capability>
	relPath, err := filepath.Rel(specsDir, deltaSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	// Extract capability (should be the directory name before spec.md)
	capability := filepath.Dir(relPath)
	if capability == "." || capability == "" {
		return nil, fmt.Errorf("invalid delta spec path structure: %s", deltaSpecPath)
	}

	// Construct base spec path
	baseSpecPath := filepath.Join(spectrRoot, "specs", capability, "spec.md")

	// Check if base spec exists
	_, err = os.Stat(baseSpecPath)
	baseExists := err == nil

	// Parse delta spec
	deltaPlan, err := parsers.ParseDeltaSpec(deltaSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse delta spec: %w", err)
	}

	// Validate delta against base spec
	if err := ValidatePreMerge(baseSpecPath, deltaPlan, baseExists); err != nil {
		return []ValidationIssue{
			{
				Level:   LevelError,
				Path:    deltaSpecPath,
				Message: err.Error(),
			},
		}, nil
	}

	return nil, nil
}
