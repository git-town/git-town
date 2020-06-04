package drivers

import "github.com/git-town/git-town/src/git"

// Core provides the public API for the drivers subsystem.

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services.
type CodeHostingDriver interface {

	// LoadPullRequestInfo loads information about the pull request of the given branch into the given parent branch
	// from the code hosting provider.
	LoadPullRequestInfo(branch, parentBranch string) (canMerge bool, defaultCommitMessage string, pullRequestNumber int64, err error)

	// NewPullRequestURL returns the URL of the page
	// to create a new pull request online
	NewPullRequestURL(branch, parentBranch string) string

	// MergePullRequest merges the pull request through the hosting service api
	MergePullRequest(MergePullRequestOptions) (mergeSha string, err error)

	// RepositoryURL returns the URL where the given repository
	// can be found online
	RepositoryURL() string

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
