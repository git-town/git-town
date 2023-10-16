package hosting

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting/bitbucket"
	"github.com/git-town/git-town/v9/src/hosting/common"
	"github.com/git-town/git-town/v9/src/hosting/gitea"
	"github.com/git-town/git-town/v9/src/hosting/github"
	"github.com/git-town/git-town/v9/src/hosting/gitlab"
)

type Connector common.Connector

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (common.Connector, error) {
	githubConnector, err := github.NewConnector(github.NewConnectorArgs{
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
	bitbucketConnector, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
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
	giteaConnector, err := gitea.NewConnector(gitea.NewConnectorArgs{
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
	GetSHAForBranch common.SHAForBranchFunc
	GiteaAPIToken   string
	GithubAPIToken  string
	GitlabAPIToken  string
	MainBranch      domain.LocalBranchName
	Log             common.Log
}
