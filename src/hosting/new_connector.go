package hosting

import (
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting/bitbucket"
	"github.com/git-town/git-town/v14/src/hosting/gitea"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/git-town/git-town/v14/src/hosting/gitlab"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (hostingdomain.Connector, error) { // TODO wrap in option here to avoid the nil interface problem?
	platform, hasPlatform := Detect(args.OriginURL, args.HostingPlatform).Get()
	if !hasPlatform {
		return nil, nil
	}
	switch platform {
	case configdomain.HostingPlatformBitbucket:
		return bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			OriginURL:       args.OriginURL,
		}), nil
	case configdomain.HostingPlatformGitea:
		return gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  args.GiteaToken,
			Log:       args.Log,
			OriginURL: args.OriginURL,
		})
	case configdomain.HostingPlatformGitHub:
		return github.NewConnector(github.NewConnectorArgs{
			APIToken:   github.GetAPIToken(args.GitHubToken),
			Log:        args.Log,
			MainBranch: args.MainBranch,
			OriginURL:  args.OriginURL,
		})
	case configdomain.HostingPlatformGitLab:
		return gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  args.GitLabToken,
			Log:       args.Log,
			OriginURL: args.OriginURL,
		})
	}
	return nil, nil
}

type NewConnectorArgs struct {
	configdomain.FullConfig
	HostingPlatform Option[configdomain.HostingPlatform]
	Log             print.Logger
	OriginURL       giturl.Parts
}
