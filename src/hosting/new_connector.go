package hosting

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/bitbucket"
	"github.com/git-town/git-town/v11/src/hosting/gitea"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/hosting/gitlab"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (hostingdomain.Connector, error) {
	githubConnector, err := github.NewConnector(github.NewConnectorArgs{
		HostingPlatform: args.HostingPlatform,
		APIToken:        github.GetAPIToken(args.GitHubToken),
		MainBranch:      args.MainBranch,
		OriginURL:       args.OriginURL,
		Log:             args.Log,
	})
	if githubConnector != nil || err != nil {
		return githubConnector, err
	}
	gitlabConnector, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
		HostingPlatform: args.HostingPlatform,
		OriginURL:       args.OriginURL,
		APIToken:        args.GitLabToken,
		Log:             args.Log,
	})
	if gitlabConnector != nil || err != nil {
		return gitlabConnector, err
	}
	bitbucketConnector, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
		OriginURL:       args.OriginURL,
		HostingPlatform: args.HostingPlatform,
	})
	if bitbucketConnector != nil || err != nil {
		return bitbucketConnector, err
	}
	giteaConnector, err := gitea.NewConnector(gitea.NewConnectorArgs{
		OriginURL:       args.OriginURL,
		HostingPlatform: args.HostingPlatform,
		APIToken:        args.GiteaToken,
		Log:             args.Log,
	})
	if giteaConnector != nil || err != nil {
		return giteaConnector, err
	}
	return nil, nil
}

type NewConnectorArgs struct {
	*configdomain.FullConfig
	HostingPlatform configdomain.HostingPlatform
	Log             hostingdomain.Log
	OriginURL       *giturl.Parts
}
