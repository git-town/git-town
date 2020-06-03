package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem.

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services.
type CodeHostingDriver interface {

	// CanMergePullRequest returns whether or not MergePullRequest should be
	// called when shipping. If true, also returns the default commit message
	CanMergePullRequest(branch, parentBranch string) (canMerge bool, defaultCommitMessage string, pullRequestNumber int64, err error)

	// GetNewPullRequestURL returns the URL of the page
	// to create a new pull request online
	GetNewPullRequestURL(branch, parentBranch string) string

	// MergePullRequest merges the pull request through the hosting service api
	MergePullRequest(MergePullRequestOptions) (mergeSha string, err error)

	// GetRepositoryURL returns the URL where the given repository
	// can be found online
	GetRepositoryURL() string

	// HostingServiceName returns the name of the code hosting service
	HostingServiceName() string
}

// MergePullRequestOptions defines the options to the MergePullRequest function.
type MergePullRequestOptions struct {
	Branch            string
	PullRequestNumber int64
	CommitMessage     string
	LogRequests       bool
	ParentBranch      string
}

// Load returns the code hosting driver to use based on the git config.
// nolint:interfacer  // for Gitea support later
func Load(config *git.Configuration) CodeHostingDriver {
	driver := LoadGithub(config)
	if driver != nil {
		return driver
	}
	driver = LoadBitbucket(config)
	if driver != nil {
		return driver
	}
	driver = LoadGitlab(config)
	if driver != nil {
		return driver
	}
	return nil
}

// UnsupportedHostingError provides an error message.
func UnsupportedHostingError() string {
	return `Unsupported hosting service

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab`
}
