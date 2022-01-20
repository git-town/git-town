package hosting

import "errors"

// Driver defines the structure of drivers for the different code hosting services.
type Driver interface {
	// LoadPullRequestInfo loads information about the pull request of the given branch into the given parent branch
	// from the code hosting provider.
	LoadPullRequestInfo(branch, parentBranch string) (PullRequestInfo, error)

	// NewPullRequestURL provides the URL of the page
	// to create a new pull request online.
	NewPullRequestURL(branch, parentBranch string) (string, error)

	// MergePullRequest merges the pull request through the hosting service API.
	MergePullRequest(MergePullRequestOptions) (mergeSha string, err error)

	// RepositoryURL provides the URL where the given repository
	// can be found online.
	RepositoryURL() string

	// HostingServiceName provides the name of the code hosting service.
	HostingServiceName() string
}

// config defines the configuration data needed by the driver package.
type config interface {
	OriginHost() string
	HostingService() string
	GiteaToken() string
	GitHubToken() string
	MainBranch() string
	RemoteOriginURL() string
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
type logFn func(string, ...interface{})

// NewDriver provides the code hosting driver to use based on the git config.
func NewDriver(config config, git gitRunner, log logFn) Driver { //nolint:ireturn
	githubDriver := LoadGithub(config, log)
	if githubDriver != nil {
		return githubDriver
	}
	giteaDriver := LoadGitea(config, log)
	if giteaDriver != nil {
		return giteaDriver
	}
	bitbucketDriver := LoadBitbucket(config, git)
	if bitbucketDriver != nil {
		return bitbucketDriver
	}
	gitlabDriver := LoadGitlab(config)
	if gitlabDriver != nil {
		return gitlabDriver
	}
	return nil
}

// UnsupportedServiceError provides an error message.
func UnsupportedServiceError() error {
	return errors.New(`unsupported hosting service

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab
* Gitea`)
}
