package hosting

import (
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/bitbucketcloud"
	"github.com/git-town/git-town/v16/internal/hosting/gitea"
	"github.com/git-town/git-town/v16/internal/hosting/github"
	"github.com/git-town/git-town/v16/internal/hosting/gitlab"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// NewConnector provides an instance of the code hosting connector to use based on the given gitConfig.
func NewConnector(config config.UnvalidatedConfig, remote gitdomain.Remote, log print.Logger) (Option[hostingdomain.Connector], error) {
	remoteURL, hasRemoteURL := config.NormalConfig.RemoteURL(remote).Get()
	hostingPlatform := config.NormalConfig.HostingPlatform
	platform, hasPlatform := Detect(remoteURL, hostingPlatform).Get()
	if !hasRemoteURL || !hasPlatform {
		return None[hostingdomain.Connector](), nil
	}
	var connector hostingdomain.Connector
	switch platform {
	case configdomain.HostingPlatformBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword:     config.NormalConfig.BitbucketAppPassword,
			HostingPlatform: hostingPlatform,
			Log:             log,
			RemoteURL:       remoteURL,
			UserName:        config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case configdomain.HostingPlatformGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  config.NormalConfig.GiteaToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), nil
	case configdomain.HostingPlatformGitHub:
		var err error
		connector, err = github.NewConnector(github.NewConnectorArgs{
			APIToken:  github.GetAPIToken(config.NormalConfig.GitHubToken),
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	case configdomain.HostingPlatformGitLab:
		var err error
		connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  config.NormalConfig.GitLabToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	}
	return None[hostingdomain.Connector](), nil
}
