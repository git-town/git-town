// Package hosting provides support for interacting with code hosting services.
// Commands like "new-pull-request", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, Gitlab, Bitbucket, etc.
// Drivers implement the CodeHostingDriver interface.
package hosting

import (
	"errors"

	"github.com/git-town/git-town/v7/src/giturl"
)

// Config contains the information needed for all platform connectors.
type Config struct {
	// bearer token to authenticate with the API
	apiToken string

	// hostname override
	hostname string

	// where the "origin" remote points to
	originURL string

	// the organization within the hosting platform that owns the repo
	// TODO: rename no "organization"
	owner string

	// repo name within the organization
	repository string
}

// Connector describes the activities that Git Town performs on code hosting platforms via their API.
// Individual implementations exist to talk to specific hosting platforms.
// They all conform to this interface.
type Connector interface {
	// ChangeRequestForBranch provides details about the change request for the given branch.
	ChangeRequestForBranch(branch string) (*ChangeRequestInfo, error)

	// DefaultCommitMessage provides the commit message template to use
	// for change requests on the respective hosting platform.
	DefaultCommitMessage(crInfo ChangeRequestInfo) string

	// HostingServiceName provides the name of the code hosting service
	// supported by the respective connector implementation.
	HostingServiceName() string

	// SquashMergeChangeRequest squash-merges the change request with the given number
	// using the given commit message.
	SquashMergeChangeRequest(number int, message string) (mergeSHA string, err error)

	// NewChangeRequestURL provides the URL of the page
	// to create a new pull request online.
	NewChangeRequestURL(branch, parentBranch string) (string, error)

	// RepositoryURL provides the URL where the current repository can be found online.
	RepositoryURL() string

	// UpdateChangeRequestTarget updates the target branch of the given change request.
	UpdateChangeRequestTarget(number int, target string) error
}

// ChangeRequestInfo contains information about a change request
// on a code hosting platform.
type ChangeRequestInfo struct {
	// the change request Number
	Number int

	// textual title of the change request
	Title string

	// whether this change request can be merged programmatically
	CanMergeWithAPI bool
}

// config defines the configuration data needed by the driver package.
type gitConfig interface {
	// OriginOverride provides the override for the origin URL from the Git Town configuration.
	OriginOverride() string

	// HostingService provides the name of the hosting service that runs at the origin remote.
	HostingService() string

	// GiteaToken provides the personal access token for Gitea stored in the Git configuration.
	GiteaToken() string

	// GitHubToken provides the personal access token for GitHub stored in the Git configuration.
	GitHubToken() string

	// GitLabToken provides the personal access token for GitLab stored in the Git configuration.
	GitLabToken() string

	// MainBranch provides the name of the main branch.
	MainBranch() string

	// OriginURL provides the URL of the origin remote.
	OriginURL() string
}

// runner defines the runner methods used by the driver package.
type gitRunner interface {
	ShaForBranch(string) (string, error)
}

// logFn defines a function with fmt.Printf API that CodeHostingDriver instances can use to give updates on activities they do.
type logFn func(string, ...interface{})

// NewDriver provides an instance of the code hosting driver to use based on the git config.
//
//nolint:ireturn,nolintlint
func NewConnector(config gitConfig, git gitRunner, log logFn) (Connector, error) {
	url := giturl.Parse(config.OriginURL())
	if url == nil {
		return nil, nil //nolint:nilnil  // "nil, nil" is a legitimate return value here
	}
	githubConnector := NewGithubConnector(*url, config, log)
	if githubConnector != nil {
		return githubConnector, nil
	}
	gitlabConnector, err := NewGitlabConnector(*url, config, log)
	if err != nil {
		return nil, err
	}
	if gitlabConnector != nil {
		return gitlabConnector, nil
	}
	bitbucketConnector := NewBitbucketConnector(*url, config, git)
	if bitbucketConnector != nil {
		return bitbucketConnector, nil
	}
	giteaConnector := NewGiteaConnector(*url, config, log)
	if giteaConnector != nil {
		return giteaConnector, nil
	}
	return nil, nil //nolint:nilnil  // "nil, nil" is a legitimate return value here
}

// UnsupportedServiceError communicates that the origin remote runs an unknown code hosting service.
func UnsupportedServiceError() error {
	return errors.New(`unsupported hosting service

This command requires hosting on one of these services:
* Bitbucket
* GitHub
* GitLab
* Gitea`)
}
