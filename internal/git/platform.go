// Package git provides utilities for git operations and pull request creation.
// It supports GitHub, GitLab, and Gitea/Forgejo platforms.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Platform represents a git hosting platform
type Platform string

const (
	// PlatformGitHub represents GitHub hosting
	PlatformGitHub Platform = "github"
	// PlatformGitLab represents GitLab hosting
	PlatformGitLab Platform = "gitlab"
	// PlatformGitea represents Gitea/Forgejo hosting
	PlatformGitea Platform = "gitea"
	// PlatformUnknown represents an unknown platform
	PlatformUnknown Platform = "unknown"
)

// DetectPlatform detects the git hosting platform from the origin remote URL
func DetectPlatform() (Platform, string, error) {
	// Get the origin remote URL
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return PlatformUnknown, "", fmt.Errorf("get origin remote URL: %w", err)
	}

	url := strings.TrimSpace(string(output))
	if url == "" {
		return PlatformUnknown, "", fmt.Errorf("origin remote URL is empty")
	}

	// Detect platform from URL
	platform := detectPlatformFromURL(url)

	return platform, url, nil
}

// detectPlatformFromURL detects the platform from a git remote URL
func detectPlatformFromURL(url string) Platform {
	url = strings.ToLower(url)

	// Check for GitHub
	if strings.Contains(url, "github.com") {
		return PlatformGitHub
	}

	// Check for GitLab
	if strings.Contains(url, "gitlab.com") || strings.Contains(url, "gitlab") {
		return PlatformGitLab
	}

	// Check for Gitea/Forgejo
	if strings.Contains(url, "gitea") || strings.Contains(url, "forgejo") {
		return PlatformGitea
	}

	return PlatformUnknown
}

// GetCLITool returns the appropriate CLI tool for the platform
func GetCLITool(platform Platform) (string, error) {
	switch platform {
	case PlatformGitHub:
		return "gh", nil
	case PlatformGitLab:
		return "glab", nil
	case PlatformGitea:
		return "tea", nil
	case PlatformUnknown:
		fallthrough
	default:
		return "", fmt.Errorf("unknown platform: %s", platform)
	}
}

// CheckCLIToolInstalled checks if the CLI tool for the platform is installed
func CheckCLIToolInstalled(platform Platform) error {
	tool, err := GetCLITool(platform)
	if err != nil {
		return err
	}

	// Check if tool is in PATH
	_, err = exec.LookPath(tool)
	if err != nil {
		installURL := getInstallURL(platform)
		return fmt.Errorf(
			"%s not found in PATH. Install from: %s",
			tool,
			installURL,
		)
	}

	return nil
}

// getInstallURL returns the installation URL for the platform's CLI tool
func getInstallURL(platform Platform) string {
	switch platform {
	case PlatformGitHub:
		return "https://cli.github.com/"
	case PlatformGitLab:
		return "https://gitlab.com/gitlab-org/cli"
	case PlatformGitea:
		return "https://gitea.com/gitea/tea"
	case PlatformUnknown:
		fallthrough
	default:
		return ""
	}
}
