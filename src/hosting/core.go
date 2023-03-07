// Package hosting provides support for interacting with code hosting services.
// Commands like "new-pull-request", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, Gitlab, Bitbucket, etc.
// Implementations of connectors for particular code hosting platforms conform to the Connector interface.
package hosting

import (
	"errors"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/giturl"
)

// Connector describes the activities that Git Town performs on code hosting platforms via their API.
// Individual implementations exist to talk to specific hosting platforms.
// They all conform to this interface.
type Connector interface {
	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective hosting platform is prepopulated with.
	DefaultProposalMessage(proposal Proposal) string

	// FindProposal provides details about the proposal for the given branch into the given target branch.
	// Returns nil if no proposal exists.
	FindProposal(branch, target string) (*Proposal, error)

	// HostingServiceName provides the name of the code hosting service
	// supported by the respective connector implementation.
	HostingServiceName() string

	// SquashMergeProposal squash-merges the proposal with the given number
	// using the given commit message.
	SquashMergeProposal(number int, message string) (mergeSHA string, err error)

	// NewProposalURL provides the URL of the page
	// to create a new proposal online.
	NewProposalURL(branch, parentBranch string) (string, error)

	// RepositoryURL provides the URL where the current repository can be found online.
	RepositoryURL() string

	// UpdateProposalTarget updates the target branch of the given proposal.
	UpdateProposalTarget(number int, target string) error
}

// CommonConfig contains data needed by all platform connectors.
type CommonConfig struct {
	// bearer token to authenticate with the API
	APIToken string

	// Hostname override
	Hostname string

	// the Organization within the hosting platform that owns the repo
	Organization string

	// repo name within the organization
	Repository string
}

// Proposal contains information about a change request
// on a code hosting platform.
// Alternative names are "pull request" or "merge request".
type Proposal struct {
	// the number used to identify the proposal on the hosting platform
	Number int

	// name of the target branch ("base") of this proposal
	Target string

	// textual title of the proposal
	Title string

	// whether this proposal can be merged via the API
	CanMergeWithAPI bool
}

// gitTownConfig defines the configuration data needed by the hosting package.
// This extra interface is necessary to access config.GitTown without creating a cyclic dependency.
type gitTownConfig interface {
	// OriginOverride provides the override for the origin URL from the Git Town configuration.
	OriginOverride() string

	// HostingService provides the name of the hosting service that runs at the origin remote.
	HostingService() (config.HostingService, error)

	// GiteaToken provides the personal access token for Gitea stored in the Git configuration.
	GiteaToken() string

	// GitHubToken provides the personal access token for GitHub stored in the Git configuration.
	GitHubToken() string

	// GitLabToken provides the personal access token for GitLab stored in the Git configuration.
	GitLabToken() string

	// MainBranch provides the name of the main branch.
	MainBranch() string

	// OriginURL provides the URL of the origin remote.
	OriginURL() *giturl.Parts
}

// runner defines the runner methods used by the hosting package.
type gitRunner interface {
	ShaForBranch(string) (string, error)
}

// logFn defines a function with fmt.Printf API that Connector instances can use to give updates on activities they do.
type logFn func(string, ...interface{})

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
//
//nolint:ireturn,nolintlint
func NewConnector(config gitTownConfig, git gitRunner, log logFn) (Connector, error) {
	githubConnector, err := NewGithubConnector(config, log)
	if err != nil {
		return nil, err
	}
	if githubConnector != nil {
		return githubConnector, nil
	}
	gitlabConnector, err := NewGitlabConnector(config, log)
	if err != nil {
		return nil, err
	}
	if gitlabConnector != nil {
		return gitlabConnector, nil
	}
	bitbucketConnector, err := NewBitbucketConnector(config, git)
	if err != nil {
		return nil, err
	}
	if bitbucketConnector != nil {
		return bitbucketConnector, nil
	}
	giteaConnector, err := NewGiteaConnector(config, log)
	if err != nil {
		return nil, err
	}
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
