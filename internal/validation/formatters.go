// Package validation provides validation result printing functions.
// This file contains functions for outputting validation results
// in both JSON and human-readable formats.
package validation

import (
	"encoding/json"
	"fmt"
	"os"
)

// BulkResult represents the result of validating a single item
type BulkResult struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Valid  bool              `json:"valid"`
	Report *ValidationReport `json:"report,omitempty"`
	Error  string            `json:"error,omitempty"`
}

// PrintJSONReport prints a single validation report as JSON
func PrintJSONReport(
	report *ValidationReport,
) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)

		return
	}
	fmt.Println(string(data))
}

// PrintHumanReport prints a single validation report in human format
func PrintHumanReport(
	itemName string,
	report *ValidationReport,
) {
	if report.Valid {
		fmt.Printf("✓ %s valid\n", itemName)

		return
	}

	issueCount := len(report.Issues)
	fmt.Printf("✗ %s has %d issue(s):\n", itemName, issueCount)

	for _, issue := range report.Issues {
		fmt.Printf(
			"  [%s] %s: %s\n",
			issue.Level,
			issue.Path,
			issue.Message,
		)
	}
}

// PrintBulkJSONResults prints bulk validation results as JSON
func PrintBulkJSONResults(results []BulkResult) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)

		return
	}
	fmt.Println(string(data))
}

// PrintBulkHumanResults prints bulk validation results in human format
func PrintBulkHumanResults(results []BulkResult) {
	passCount := 0
	failCount := 0

	for _, result := range results {
		if result.Valid {
			fmt.Printf("✓ %s (%s)\n", result.Name, result.Type)
			passCount++
		} else {
			if result.Error != "" {
				fmt.Printf(
					"✗ %s (%s): %s\n",
					result.Name,
					result.Type,
					result.Error,
				)
			} else {
				issueCount := len(result.Report.Issues)
				fmt.Printf(
					"✗ %s (%s) has %d issue(s):\n",
					result.Name,
					result.Type,
					issueCount,
				)
				for _, issue := range result.Report.Issues {
					fmt.Printf(
						"  [%s] %s: %s\n",
						issue.Level,
						issue.Path,
						issue.Message,
					)
				}
			}
			failCount++
		}
	}

	fmt.Printf(
		"\n%d passed, %d failed, %d total\n",
		passCount,
		failCount,
		len(results),
	)
}
