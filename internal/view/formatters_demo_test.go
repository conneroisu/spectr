package view

import (
	"fmt"
	"testing"
)

// TestFormatDashboardText_VisualDemo generates a visual demonstration
// of the dashboard output for manual inspection
func TestFormatDashboardText_VisualDemo(_ *testing.T) {
	// Create sample data matching the design.md specification
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        5,
			TotalRequirements: 42,
			ActiveChanges:     2,
			CompletedChanges:  1,
			TotalTasks:        15,
			CompletedTasks:    8,
		},
		ActiveChanges: []ChangeProgress{
			{
				ID:    "add-view-command",
				Title: "Add View Command",
				Progress: ProgressMetrics{
					Total:      8,
					Completed:  3,
					Percentage: 37,
				},
			},
			{
				ID:    "add-validate-command",
				Title: "Add Validate Command",
				Progress: ProgressMetrics{
					Total:      5,
					Completed:  3,
					Percentage: 60,
				},
			},
		},
		CompletedChanges: []CompletedChange{
			{
				ID:    "add-list-command",
				Title: "Add List Command",
			},
		},
		Specs: []SpecInfo{
			{
				ID:               "cli-framework",
				Title:            "CLI Framework",
				RequirementCount: 15,
			},
			{
				ID:               "validation",
				Title:            "Validation",
				RequirementCount: 12,
			},
		},
	}

	output := FormatDashboardText(data)

	// Print the formatted output for visual inspection
	separator := "================================================================================"
	fmt.Println("\n" + separator)
	fmt.Println("DASHBOARD OUTPUT DEMO")
	fmt.Println(separator)
	fmt.Println(output)
	fmt.Println(separator)
}
