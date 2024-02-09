package hosting

import (
	"github.com/git-town/git-town/v11/src/cli/print"
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
	switch Detect(args.OriginURL, args.HostingPlatform) {
	case configdomain.HostingPlatformBitbucket:
		return bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			OriginURL:       args.OriginURL,
			HostingPlatform: args.HostingPlatform,
		})
	case configdomain.HostingPlatformGitea:
		return gitea.NewConnector(gitea.NewConnectorArgs{
			OriginURL:       args.OriginURL,
			HostingPlatform: args.HostingPlatform,
			APIToken:        args.GiteaToken,
			Log:             args.Log,
		})
	case configdomain.HostingPlatformGitHub:
		return github.NewConnector(github.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			APIToken:        github.GetAPIToken(args.GitHubToken),
			MainBranch:      args.MainBranch,
			OriginURL:       args.OriginURL,
			Log:             args.Log,
		})
	case configdomain.HostingPlatformGitLab:
		return gitlab.NewConnector(gitlab.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			OriginURL:       args.OriginURL,
			APIToken:        args.GitLabToken,
			Log:             args.Log,
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
