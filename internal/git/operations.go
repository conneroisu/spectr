package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not in a git repository. Initialize git with 'git init'")
	}

	return nil
}

// HasOriginRemote checks if the origin remote is configured
func HasOriginRemote() error {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("no 'origin' remote configured. Run 'git remote add origin <url>'")
	}
	if strings.TrimSpace(string(output)) == "" {
		return fmt.Errorf("origin remote URL is empty")
	}

	return nil
}

// CreateBranch creates a new git branch with the given name
func CreateBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create branch: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// StageFiles stages files for commit
func StageFiles(paths []string) error {
	args := append([]string{"add"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stage files: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Commit creates a git commit with the given message
func Commit(message string) error {
	// Use heredoc-style message passing via stdin
	cmd := exec.Command("git", "commit", "-F", "-")
	cmd.Stdin = strings.NewReader(message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("commit: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Push pushes the current branch to origin
func Push(branchName string) error {
	cmd := exec.Command("git", "push", "-u", "origin", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("push branch: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetCurrentBranch returns the name of the current git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CheckoutBranch switches to the specified git branch
func CheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("checkout branch: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// BranchExists checks if a branch exists locally
func BranchExists(branchName string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branchName)
	err := cmd.Run()

	return err == nil
}
