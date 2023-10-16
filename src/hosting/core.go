// Package hosting provides support for interacting with code hosting services.
// Commands like "new-pull-request", "repo", and "ship" use this package
// to know how to perform Git Town operations on GitHub, Gitlab, Bitbucket, etc.
// Implementations of connectors for particular code hosting platforms conform to the Connector interface.
package hosting

import (
	"errors"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting/common"
	"github.com/git-town/git-town/v9/src/hosting/gitlab"
)

type Connector common.Connector

// gitTownConfig defines the configuration data needed by the hosting package.
// This extra interface is necessary to access config.GitTown without creating a cyclic dependency.
type gitTownConfig interface {
	// OriginOverride provides the override for the origin URL from the Git Town configuration.
	OriginOverride() string

	// HostingService provides the name of the hosting service that runs at the origin remote.
	HostingService() (config.Hosting, error)

	// GiteaToken provides the personal access token for Gitea stored in the Git configuration.
	GiteaToken() string

	// GitHubToken provides the personal access token for GitHub stored in the Git configuration.
	GitHubToken() string

	// GitLabToken provides the personal access token for GitLab stored in the Git configuration.
	GitLabToken() string

	// MainBranch provides the name of the main branch.
	MainBranch() domain.LocalBranchName

	// OriginURL provides the URL of the origin remote.
	OriginURL() *giturl.Parts
}

type SHAForBranchFunc func(domain.BranchName) (domain.SHA, error)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (common.Connector, error) {
	githubConnector, err := NewGithubConnector(NewGithubConnectorArgs{
		HostingService: args.HostingService,
		APIToken:       args.GithubAPIToken,
		MainBranch:     args.MainBranch,
		OriginURL:      args.OriginURL,
		Log:            args.Log,
	})
	if err != nil {
		return nil, err
	}
	if githubConnector != nil {
		return githubConnector, nil
	}
	gitlabConnector, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
		HostingService: args.HostingService,
		OriginURL:      args.OriginURL,
		APIToken:       args.GitlabAPIToken,
		Log:            args.Log,
	})
	if err != nil {
		return nil, err
	}
	if gitlabConnector != nil {
		return gitlabConnector, nil
	}
	bitbucketConnector, err := NewBitbucketConnector(NewBitbucketConnectorArgs{
		OriginURL:       args.OriginURL,
		HostingService:  args.HostingService,
		GetSHAForBranch: args.GetSHAForBranch,
	})
	if err != nil {
		return nil, err
	}
	if bitbucketConnector != nil {
		return bitbucketConnector, nil
	}
	giteaConnector, err := NewGiteaConnector(NewGiteaConnectorArgs{
		OriginURL:      args.OriginURL,
		HostingService: args.HostingService,
		APIToken:       args.GiteaAPIToken,
		Log:            args.Log,
	})
	if err != nil {
		return nil, err
	}
	if giteaConnector != nil {
		return giteaConnector, nil
	}
	return nil, nil
}

type NewConnectorArgs struct {
	HostingService  config.Hosting
	OriginURL       *giturl.Parts
	GetSHAForBranch SHAForBranchFunc
	GiteaAPIToken   string
	GithubAPIToken  string
	GitlabAPIToken  string
	MainBranch      domain.LocalBranchName
	Log             common.Log
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
