package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/conneroisu/spectr/internal/git"
)

// PRContext holds the context needed for PR creation
type PRContext struct {
	ChangeID     string
	ArchiveName  string
	SkipSpecs    bool
	OpCounts     OperationCounts
	Capabilities []string
	SpectrRoot   string
}

// createPR orchestrates the PR creation workflow after successful archive
func createPR(ctx PRContext) error {
	fmt.Println("\nCreating pull request...")

	platform, err := validateAndDetectPlatform()
	if err != nil {
		return err
	}

	originalBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	branchName := fmt.Sprintf("archive-%s", ctx.ChangeID)
	if err := prepareBranchAndCommit(ctx, branchName); err != nil {
		return err
	}

	prURL, err := pushAndCreatePR(ctx, platform, branchName)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Pull request created: %s\n", prURL)

	restoreOriginalBranch(originalBranch, branchName)

	return nil
}

// validateAndDetectPlatform validates git environment and detects platform.
// It checks git repository status, detects the platform, and verifies CLI tool
// installation. Returns platform and error if any validation fails.
func validateAndDetectPlatform() (git.Platform, error) {
	if err := validateGitEnvironment(); err != nil {
		return git.PlatformUnknown, err
	}

	platform, remoteURL, err := git.DetectPlatform()
	if err != nil {
		msg := "detect git platform: %w. " +
			"Archive completed successfully. Create PR manually"
		return git.PlatformUnknown, fmt.Errorf(msg, err)
	}

	if platform == git.PlatformUnknown {
		msg := "could not detect git hosting platform. Remote URL: %s. " +
			"Archive completed successfully. Create PR manually"
		return git.PlatformUnknown, fmt.Errorf(msg, remoteURL)
	}

	if err := git.CheckCLIToolInstalled(platform); err != nil {
		msg := "%w. Archive completed successfully. Create PR manually"
		return git.PlatformUnknown, fmt.Errorf(msg, err)
	}

	return platform, nil
}

// prepareBranchAndCommit creates a branch, stages files, and commits.
// It creates a new git branch, stages the archive files and specs,
// then commits with an automatically generated message.
func prepareBranchAndCommit(ctx PRContext, branchName string) error {
	if err := createBranch(branchName); err != nil {
		return err
	}

	if err := stageArchiveFiles(ctx); err != nil {
		return err
	}

	commitMsg := buildCommitMessage(ctx)
	if err := git.Commit(commitMsg); err != nil {
		msg := "%v. Archive completed. Branch created. " +
			"Commit manually and push"
		return fmt.Errorf(msg, err)
	}

	return nil
}

// pushAndCreatePR pushes the branch and creates a pull request.
// It pushes the archive branch to origin and creates a PR/MR using the
// platform-specific CLI tool (gh, glab, or tea). Returns PR URL on success.
func pushAndCreatePR(
	ctx PRContext,
	platform git.Platform,
	branchName string,
) (string, error) {
	if err := git.Push(branchName); err != nil {
		msg := "%v. Archive completed. Branch created and committed. " +
			"Push manually"
		return "", fmt.Errorf(msg, err)
	}

	// Build PR options with title and body
	prOpts := git.PROptions{
		Title:  buildPRTitle(ctx.ChangeID),
		Body:   buildPRBody(ctx),
		Branch: branchName,
	}

	prURL, err := git.CreatePR(platform, prOpts)
	if err != nil {
		msg := "%v. Archive completed. Branch pushed. Create PR manually"
		return "", fmt.Errorf(msg, err)
	}

	return prURL, nil
}

// restoreOriginalBranch attempts to restore the original branch.
// After PR creation, this function returns to the branch the user was on
// before the archive. If it fails, a warning is printed but no error returned.
func restoreOriginalBranch(originalBranch, branchName string) {
	if err := git.CheckoutBranch(originalBranch); err != nil {
		msg := "\nWarning: Failed to restore original branch '%s': %v\n" +
			"You are still on branch '%s'. " +
			"Checkout manually with: git checkout %s\n"
		fmt.Fprintf(
			os.Stderr,
			msg,
			originalBranch,
			err,
			branchName,
			originalBranch,
		)
	} else {
		fmt.Printf("Restored original branch: %s\n", originalBranch)
	}
}

// validateGitEnvironment checks if git is properly configured
func validateGitEnvironment() error {
	if err := git.IsGitRepository(); err != nil {
		return err
	}

	return git.HasOriginRemote()
}

// createBranch creates a new git branch for the archive.
// It checks for existing branches with the same name and creates a new
// branch with the pattern "archive-{change-id}".
func createBranch(branchName string) error {
	// Check if branch already exists
	if git.BranchExists(branchName) {
		msg := "branch '%s' already exists. " +
			"Delete it first with: git branch -D %s"
		return fmt.Errorf(msg, branchName, branchName)
	}

	if err := git.CreateBranch(branchName); err != nil {
		return err
	}

	fmt.Printf("Created branch: %s\n", branchName)

	return nil
}

// stageArchiveFiles stages the archived directory and updated specs
func stageArchiveFiles(ctx PRContext) error {
	paths := []string{
		filepath.Join(ctx.SpectrRoot, "changes", "archive", ctx.ArchiveName),
	}

	// Add specs directory if specs were updated
	if !ctx.SkipSpecs {
		paths = append(paths, filepath.Join(ctx.SpectrRoot, "specs"))
	}

	if err := git.StageFiles(paths); err != nil {
		return err
	}

	fmt.Println("Staged files for commit")
	return nil
}

// buildCommitMessage generates the commit message for the archive
func buildCommitMessage(ctx PRContext) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Archive: %s\n\n", ctx.ChangeID))
	sb.WriteString(
		fmt.Sprintf(
			"Completed change '%s' archived to changes/archive/%s/\n",
			ctx.ChangeID,
			ctx.ArchiveName,
		),
	)

	if !ctx.SkipSpecs && ctx.OpCounts.Total() > 0 {
		sb.WriteString("\nSpec operations applied:\n")
		sb.WriteString(fmt.Sprintf("+ %d added\n", ctx.OpCounts.Added))
		sb.WriteString(fmt.Sprintf("~ %d modified\n", ctx.OpCounts.Modified))
		sb.WriteString(fmt.Sprintf("- %d removed\n", ctx.OpCounts.Removed))
		sb.WriteString(fmt.Sprintf("→ %d renamed\n", ctx.OpCounts.Renamed))
	}

	sb.WriteString(fmt.Sprintf("\nChange-Id: %s\n", ctx.ChangeID))

	return sb.String()
}

// buildPRTitle generates the PR title
func buildPRTitle(changeID string) string {
	return fmt.Sprintf("Archive: %s", changeID)
}

// buildPRBody generates the PR body description.
// It creates a detailed markdown body showing archive summary, spec updates,
// and review notes for the pull request.
func buildPRBody(ctx PRContext) string {
	var sb strings.Builder

	sb.WriteString("## Archive Summary\n\n")
	sb.WriteString(
		fmt.Sprintf("Archived completed change: `%s`\n\n", ctx.ChangeID),
	)
	sb.WriteString(
		fmt.Sprintf(
			"Location: `spectr/changes/archive/%s/`\n\n",
			ctx.ArchiveName,
		),
	)

	sb.WriteString("## Spec Updates\n\n")
	switch {
	case ctx.SkipSpecs:
		sb.WriteString("Spec updates skipped (--skip-specs flag used)\n\n")
	case ctx.OpCounts.Total() > 0:
		sb.WriteString("Spec operations applied:\n")
		sb.WriteString(fmt.Sprintf("- **+ %d added**\n", ctx.OpCounts.Added))
		sb.WriteString(
			fmt.Sprintf("- **~ %d modified**\n", ctx.OpCounts.Modified),
		)
		sb.WriteString(
			fmt.Sprintf("- **- %d removed**\n", ctx.OpCounts.Removed),
		)
		sb.WriteString(
			fmt.Sprintf("- **→ %d renamed**\n\n", ctx.OpCounts.Renamed),
		)

		writeCapabilities(&sb, ctx.Capabilities)
	default:
		sb.WriteString("No spec updates (no delta specs found)\n\n")
	}

	writeReviewNotes(&sb)

	return sb.String()
}

// writeCapabilities appends capability list to the string builder
func writeCapabilities(sb *strings.Builder, capabilities []string) {
	if len(capabilities) == 0 {
		return
	}

	sb.WriteString("Updated capabilities:\n")
	for _, cap := range capabilities {
		sb.WriteString(fmt.Sprintf("- `%s`\n", cap))
	}
	sb.WriteString("\n")
}

// writeReviewNotes appends review notes section to the string builder
func writeReviewNotes(sb *strings.Builder) {
	sb.WriteString("## Review Notes\n\n")
	sb.WriteString(
		"This PR archives a completed change and updates " +
			"specifications to reflect the implemented functionality. " +
			"Please review:\n",
	)
	sb.WriteString("1. Archived change structure and completeness\n")
	sb.WriteString("2. Spec delta accuracy and correctness\n")
	sb.WriteString("3. Merged spec content\n\n")

	sb.WriteString("---\n")
	sb.WriteString("Generated by `spectr archive --pr`\n")
}
