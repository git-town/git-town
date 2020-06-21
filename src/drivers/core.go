package drivers

import "errors"

// Core provides the public API for the drivers subsystem.

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services.
type CodeHostingDriver interface {

	// LoadPullRequestInfo loads information about the pull request of the given branch into the given parent branch
	// from the code hosting provider.
	LoadPullRequestInfo(branch, parentBranch string) (PullRequestInfo, error)

	// NewPullRequestURL returns the URL of the page
	// to create a new pull request online.
	NewPullRequestURL(branch, parentBranch string) (string, error)

	// MergePullRequest merges the pull request through the hosting service API.
	MergePullRequest(MergePullRequestOptions) (mergeSha string, err error)

	// RepositoryURL returns the URL where the given repository
	// can be found online.
	RepositoryURL() string

	// HostingServiceName returns the name of the code hosting service.
	HostingServiceName() string
}

// config defines the configuration data needed by the driver package.
type config interface {
	GetCodeHostingOriginHostname() string
	GetCodeHostingDriverName() string
	GetGiteaToken() string
	GetGitHubToken() string
	GetMainBranch() string
	GetRemoteOriginURL() string
}

// runner defines the runner methods used by the driver package.
type gitRunner interface {
	ShaForBranch(string) (string, error)
}

// PullRequestInfo contains information about a pull request.
type PullRequestInfo struct {
	CanMergeWithAPI      bool
	DefaultCommitMessage string
	PullRequestNumber    int64
}

// MergePullRequestOptions defines the options to the MergePullRequest function.
type MergePullRequestOptions struct {
	Branch            string
	PullRequestNumber int64
	CommitMessage     string
	LogRequests       bool
	ParentBranch      string
}

// logFn defines a function with fmt.Printf API that CodeHostingDriver instances can use to give updates on activities they do.
type logFn func(...interface{})

// Load returns the code hosting driver to use based on the git config.
// nolint:interfacer  // for Gitea support later
func Load(config config, git gitRunner, log logFn) CodeHostingDriver {
	driver := LoadGithub(config, log)
	if driver != nil {
		return driver
	}
	driver = LoadGitea(config, log)
	if driver != nil {
		return driver
	}
	driver = LoadBitbucket(config, git)
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
func UnsupportedHostingError() error {
	return errors.New(`unsupported hosting service

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab
* Gitea`)
}
