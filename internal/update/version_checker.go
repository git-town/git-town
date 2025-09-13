package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/config"
	"github.com/hashicorp/go-version"
)

// Info contains information about an available update
type Info struct {
	CurrentVersion string
	LatestVersion  string
	UpdateURL      string
}

// VersionChecker checks for Git Town updates
type VersionChecker struct {
	githubClient *GitHubClient
}

// CheckForUpdate checks if a newer version is available on GitHub
func (self *VersionChecker) CheckForUpdate(ctx context.Context) (*Info, error) {
	release, err := self.githubClient.GetLatestRelease(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching latest release: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := config.GitTownVersion

	updateInfo := &Info{
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		UpdateURL:      release.HTMLURL,
	}

	return updateInfo, nil
}

// IsUpdateAvailable checks if there is a new version available
func (self Info) IsUpdateAvailable() bool {
	current, err := version.NewVersion(self.CurrentVersion)
	if err != nil {
		return false
	}

	latest, err := version.NewVersion(self.LatestVersion)
	if err != nil {
		return false
	}

	return latest.GreaterThan(current)
}

// NewVersionChecker creates a new version checker
func NewVersionChecker() *VersionChecker {
	return &VersionChecker{
		githubClient: NewGitHubClient(),
	}
}
