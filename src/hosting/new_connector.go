package hosting

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/bitbucket"
	"github.com/git-town/git-town/v11/src/hosting/common"
	"github.com/git-town/git-town/v11/src/hosting/gitea"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/hosting/gitlab"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	githubConnector, err := github.NewConnector(github.NewConnectorArgs{
		HostingService: args.HostingService,
		APIToken:       args.GithubAPIToken,
		MainBranch:     args.MainBranch,
		OriginURL:      args.OriginURL,
		Log:            args.Log,
	})
	if githubConnector != nil || err != nil {
		return githubConnector, err
	}
	gitlabConnector, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
		HostingService: args.HostingService,
		OriginURL:      args.OriginURL,
		APIToken:       args.GitlabAPIToken,
		Log:            args.Log,
	})
	if gitlabConnector != nil || err != nil {
		return gitlabConnector, err
	}
	bitbucketConnector, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
		OriginURL:       args.OriginURL,
		HostingService:  args.HostingService,
		GetSHAForBranch: args.GetSHAForBranch,
	})
	if bitbucketConnector != nil || err != nil {
		return bitbucketConnector, err
	}
	giteaConnector, err := gitea.NewConnector(gitea.NewConnectorArgs{
		OriginURL:      args.OriginURL,
		HostingService: args.HostingService,
		APIToken:       args.GiteaAPIToken,
		Log:            args.Log,
	})
	if giteaConnector != nil || err != nil {
		return giteaConnector, err
	}
	return nil, nil //nolint:nilnil
}

type NewConnectorArgs struct {
	HostingService  configdomain.Hosting
	OriginURL       *giturl.Parts
	GetSHAForBranch common.SHAForBranchFunc
	GiteaAPIToken   configdomain.GiteaToken
	GithubAPIToken  configdomain.GitHubToken
	GitlabAPIToken  configdomain.GitLabToken
	MainBranch      domain.LocalBranchName
	Log             common.Log
}
