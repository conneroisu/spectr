package view

import (
	"fmt"
	"testing"
)

// TestRenderBarVisualDemo is not a real test - it's a visual demonstration
// of the progress bar rendering at various completion levels.
// Run with: go test -v ./internal/view -run TestRenderBarVisualDemo
func TestRenderBarVisualDemo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping visual demo in short mode")
	}

	fmt.Println("\n" + "Progress Bar Examples:")
	fmt.Println("=====================")
	fmt.Println()

	examples := []struct {
		name      string
		completed int
		total     int
	}{
		{"0% Complete", 0, 10},
		{"25% Complete", 25, 100},
		{"37% Complete (Design Example)", 37, 100},
		{"50% Complete", 5, 10},
		{"60% Complete (Design Example)", 60, 100},
		{"75% Complete", 75, 100},
		{"100% Complete", 10, 10},
		{"Zero Total Edge Case", 0, 0},
		{"3 of 8 Tasks (38%)", 3, 8},
	}

	for _, ex := range examples {
		bar := RenderBar(ex.completed, ex.total)
		fmt.Printf("%-35s %s\n", ex.name+":", bar)
	}
	fmt.Println()
}
