// Package cmd provides command-line interface implementations for Spectr.
// This file contains the archive command for archiving completed changes.
package cmd

import (
	"fmt"

	"github.com/conneroisu/spectr/internal/archive"
)

// ArchiveCmd represents the archive command
type ArchiveCmd struct {
	ChangeID    string `arg:"" optional:"" help:"Change ID to archive"`
	Yes         bool   `name:"yes" short:"y" help:"Skip confirmation"`
	SkipSpecs   bool   `name:"skip-specs" help:"Skip spec updates"`
	NoValidate  bool   `name:"no-validate" help:"Skip validation"`
	Interactive bool   `short:"I" name:"interactive" help:"Interactive mode"`
}

// Run executes the archive command
func (c *ArchiveCmd) Run() error {
	// Create archiver with flags
	archiver, err := archive.NewArchiver(
		c.Yes,
		c.SkipSpecs,
		c.NoValidate,
		c.Interactive,
	)
	if err != nil {
		return fmt.Errorf("failed to create archiver: %w", err)
	}

	// Execute archive
	err = archiver.Archive(c.ChangeID)
	if err != nil {
		return fmt.Errorf("archive failed: %w", err)
	}

	return nil
}
