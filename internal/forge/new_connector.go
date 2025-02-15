package forge

import (
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/config"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v18/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/forge/gitea"
	"github.com/git-town/git-town/v18/internal/forge/github"
	"github.com/git-town/git-town/v18/internal/forge/gitlab"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// NewConnector provides an instance of the forge connector to use based on the given gitConfig.
func NewConnector(config config.UnvalidatedConfig, remote gitdomain.Remote, log print.Logger) (Option[forgedomain.Connector], error) {
	remoteURL, hasRemoteURL := config.NormalConfig.RemoteURL(remote).Get()
	hostingPlatform := config.NormalConfig.HostingPlatform
	platform, hasPlatform := Detect(remoteURL, hostingPlatform).Get()
	if !hasRemoteURL || !hasPlatform {
		return None[forgedomain.Connector](), nil
	}
	var connector forgedomain.Connector
	switch platform {
	case configdomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword:     config.NormalConfig.BitbucketAppPassword,
			HostingPlatform: hostingPlatform,
			Log:             log,
			RemoteURL:       remoteURL,
			UserName:        config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case configdomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			AppPassword:     config.NormalConfig.BitbucketAppPassword,
			HostingPlatform: hostingPlatform,
			Log:             log,
			RemoteURL:       remoteURL,
			UserName:        config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case configdomain.ForgeTypeGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  config.NormalConfig.GiteaToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), nil
	case configdomain.ForgeTypeGitHub:
		var err error
		connector, err = github.NewConnector(github.NewConnectorArgs{
			APIToken:  github.GetAPIToken(config.NormalConfig.GitHubToken),
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	case configdomain.ForgeTypeGitLab:
		var err error
		connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  config.NormalConfig.GitLabToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	}
	return None[forgedomain.Connector](), nil
}
