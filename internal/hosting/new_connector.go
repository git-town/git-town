package hosting

import (
	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting/bitbucket"
	"github.com/git-town/git-town/v15/internal/hosting/gitea"
	"github.com/git-town/git-town/v15/internal/hosting/github"
	"github.com/git-town/git-town/v15/internal/hosting/gitlab"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (Option[hostingdomain.Connector], error) {
	url, hasURL := giturl.Parse(args.OriginURL).Get()
	if !hasURL {
		return None[hostingdomain.Connector](), nil
	}
	platform, hasPlatform := Detect(url, args.HostingPlatform).Get()
	if !hasPlatform {
		return None[hostingdomain.Connector](), nil
	}
	var connector hostingdomain.Connector
	switch platform {
	case configdomain.HostingPlatformBitbucket:
		connector = bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingPlatform: args.HostingPlatform,
			OriginURL:       url,
		})
		return Some(connector), nil
	case configdomain.HostingPlatformGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  args.Config.GiteaToken,
			Log:       args.Log,
			OriginURL: url,
		})
		return Some(connector), nil
	case configdomain.HostingPlatformGitHub:
		var err error
		connector, err = github.NewConnector(github.NewConnectorArgs{
			APIToken:  github.GetAPIToken(args.Config.GitHubToken),
			Log:       args.Log,
			OriginURL: url,
		})
		return Some(connector), err
	case configdomain.HostingPlatformGitLab:
		var err error
		connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  args.Config.GitLabToken,
			Log:       args.Log,
			OriginURL: url,
		})
		return Some(connector), err
	}
	return None[hostingdomain.Connector](), nil
}

type NewConnectorArgs struct {
	Config          configdomain.UnvalidatedConfig
	HostingPlatform Option[configdomain.HostingPlatform]
	Log             print.Logger
	OriginURL       string
}
