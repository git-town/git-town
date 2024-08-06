package hosting

import (
	"github.com/git-town/git-town/v14/internal/cli/print"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/git/giturl"
	. "github.com/git-town/git-town/v14/internal/gohacks/prelude"
	"github.com/git-town/git-town/v14/internal/hosting/bitbucket"
	"github.com/git-town/git-town/v14/internal/hosting/gitea"
	"github.com/git-town/git-town/v14/internal/hosting/github"
	"github.com/git-town/git-town/v14/internal/hosting/gitlab"
	"github.com/git-town/git-town/v14/internal/hosting/hostingdomain"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (Option[hostingdomain.Connector], error) {
	platform, hasPlatform := Detect(args.OriginURL, args.HostingPlatform).Get()
	if !hasPlatform {
		return None[hostingdomain.Connector](), nil
	}
	var connector hostingdomain.Connector
	switch platform {
	case configdomain.HostingPlatformBitbucket:
		connector = bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			OriginURL:       args.OriginURL,
		})
		return Some(connector), nil
	case configdomain.HostingPlatformGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  args.Config.GiteaToken,
			Log:       args.Log,
			OriginURL: args.OriginURL,
		})
		return Some(connector), nil
	case configdomain.HostingPlatformGitHub:
		var err error
		connector, err = github.NewConnector(github.NewConnectorArgs{
			APIToken:  github.GetAPIToken(args.Config.GitHubToken),
			Log:       args.Log,
			OriginURL: args.OriginURL,
		})
		return Some(connector), err
	case configdomain.HostingPlatformGitLab:
		var err error
		connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  args.Config.GitLabToken,
			Log:       args.Log,
			OriginURL: args.OriginURL,
		})
		return Some(connector), err
	}
	return None[hostingdomain.Connector](), nil
}

type NewConnectorArgs struct {
	Config          configdomain.UnvalidatedConfig
	HostingPlatform Option[configdomain.HostingPlatform]
	Log             print.Logger
	OriginURL       giturl.Parts
}
