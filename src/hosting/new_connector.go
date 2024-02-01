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
	platform := detect(args.OriginURL, args.HostingPlatform)
	switch platform {
	case hostingdomain.PlatformBitbucket:
		return bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			OriginURL:       args.OriginURL,
			HostingPlatform: args.HostingPlatform,
		})
	case hostingdomain.PlatformGitea:
		return gitea.NewConnector(gitea.NewConnectorArgs{
			OriginURL:       args.OriginURL,
			HostingPlatform: args.HostingPlatform,
			APIToken:        args.GiteaToken,
			Log:             args.Log,
		})
	case hostingdomain.PlatformGithub:
		return github.NewConnector(github.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			APIToken:        github.GetAPIToken(args.GitHubToken),
			MainBranch:      args.MainBranch,
			OriginURL:       args.OriginURL,
			Log:             args.Log,
		})
	case hostingdomain.PlatformGitlab:
		return gitlab.NewConnector(gitlab.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			OriginURL:       args.OriginURL,
			APIToken:        args.GitLabToken,
			Log:             args.Log,
		})
	case hostingdomain.PlatformNone:
		return nil, nil
	}
	return nil, nil
}

type NewConnectorArgs struct {
	*configdomain.FullConfig
	HostingPlatform configdomain.HostingPlatform
	Log             hostingdomain.Log
	OriginURL       *giturl.Parts
}
