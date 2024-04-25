package hosting

import (
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/hosting/bitbucket"
	"github.com/git-town/git-town/v14/src/hosting/gitea"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/git-town/git-town/v14/src/hosting/gitlab"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (hostingdomain.Connector, error) {
	switch Detect(args.OriginURL, args.HostingPlatform) {
	case configdomain.HostingPlatformBitbucket:
		return bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			OriginURL:       args.OriginURL,
		})
	case configdomain.HostingPlatformGitea:
		return gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:        args.GiteaToken,
			HostingPlatform: args.HostingPlatform,
			Log:             args.Log,
			OriginURL:       args.OriginURL,
		})
	case configdomain.HostingPlatformGitHub:
		return github.NewConnector(github.NewConnectorArgs{
			APIToken:        github.GetAPIToken(args.GitHubToken),
			HostingPlatform: args.HostingPlatform,
			Log:             args.Log,
			MainBranch:      args.MainBranch,
			OriginURL:       args.OriginURL,
		})
	case configdomain.HostingPlatformGitLab:
		return gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:        args.GitLabToken,
			HostingPlatform: args.HostingPlatform,
			Log:             args.Log,
			OriginURL:       args.OriginURL,
		})
	case configdomain.HostingPlatformNone:
		return nil, nil
	}
	return nil, nil
}

type NewConnectorArgs struct {
	*configdomain.FullConfig
	HostingPlatform configdomain.HostingPlatform
	Log             print.Logger
	OriginURL       *giturl.Parts
}
