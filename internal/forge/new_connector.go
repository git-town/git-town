package forge

import (
	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/config"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v19/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v19/internal/forge/codeberg"
	"github.com/git-town/git-town/v19/internal/forge/forgedomain"
	"github.com/git-town/git-town/v19/internal/forge/gitea"
	"github.com/git-town/git-town/v19/internal/forge/github"
	"github.com/git-town/git-town/v19/internal/forge/gitlab"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// NewConnector provides an instance of the forge connector to use based on the given gitConfig.
func NewConnector(config config.UnvalidatedConfig, remote gitdomain.Remote, log print.Logger) (Option[forgedomain.Connector], error) {
	remoteURL, hasRemoteURL := config.NormalConfig.RemoteURL(remote).Get()
	forgeType := config.NormalConfig.ForgeType
	platform, hasPlatform := Detect(remoteURL, forgeType).Get()
	if !hasRemoteURL || !hasPlatform {
		return None[forgedomain.Connector](), nil
	}
	var connector forgedomain.Connector
	switch platform {
	case configdomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword: config.NormalConfig.BitbucketAppPassword,
			ForgeType:   forgeType,
			Log:         log,
			RemoteURL:   remoteURL,
			UserName:    config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case configdomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			AppPassword:     config.NormalConfig.BitbucketAppPassword,
			HostingPlatform: forgeType,
			Log:             log,
			RemoteURL:       remoteURL,
			UserName:        config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case configdomain.ForgeTypeCodeberg:
		var err error
		connector, err = codeberg.NewConnector(codeberg.NewConnectorArgs{
			APIToken:  config.NormalConfig.CodebergToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
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
			APIToken:  config.NormalConfig.GitHubToken,
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
