package forge

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v21/internal/forge/codeberg"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gh"
	"github.com/git-town/git-town/v21/internal/forge/gitea"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
	"github.com/git-town/git-town/v21/internal/forge/glab"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// NewConnector provides an instance of the forge connector to use based on the given gitConfig.
// TODO: replace most arguments with OpenRepoResult
// or replace config.UnvalidatedConfig with config.UnvalidatedConfig.NormalConfig
func NewConnector(config config.UnvalidatedConfig, remote gitdomain.Remote, log print.Logger, frontend subshelldomain.Runner, backend subshelldomain.Querier) (Option[forgedomain.Connector], error) {
	remoteURL, hasRemoteURL := config.NormalConfig.RemoteURL(remote).Get()
	forgeType := config.NormalConfig.ForgeType
	platform, hasPlatform := Detect(remoteURL, forgeType).Get()
	if !hasRemoteURL || !hasPlatform {
		return None[forgedomain.Connector](), nil
	}
	var connector forgedomain.Connector
	var err error
	switch platform {
	case forgedomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword: config.NormalConfig.BitbucketAppPassword,
			ForgeType:   forgeType,
			Log:         log,
			RemoteURL:   remoteURL,
			UserName:    config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			AppPassword:     config.NormalConfig.BitbucketAppPassword,
			HostingPlatform: forgeType,
			Log:             log,
			RemoteURL:       remoteURL,
			UserName:        config.NormalConfig.BitbucketUsername,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeCodeberg:
		connector, err = codeberg.NewConnector(codeberg.NewConnectorArgs{
			APIToken:  config.NormalConfig.CodebergToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	case forgedomain.ForgeTypeGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  config.NormalConfig.GiteaToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeGitHub:
		if githubConnectorType, hasGitHubConnectorType := config.NormalConfig.GitHubConnectorType.Get(); hasGitHubConnectorType {
			switch githubConnectorType {
			case forgedomain.GitHubConnectorTypeAPI:
				connector, err = github.NewConnector(github.NewConnectorArgs{
					APIToken:  config.NormalConfig.GitHubToken,
					Log:       log,
					RemoteURL: remoteURL,
				})
				return Some(connector), err
			case forgedomain.GitHubConnectorTypeGh:
				connector = gh.Connector{
					Backend:  backend,
					Frontend: frontend,
				}
				return Some(connector), err
			}
		}
		// no GitHubConnectorType specified --> use the API connector
		connector, err = github.NewConnector(github.NewConnectorArgs{
			APIToken:  config.NormalConfig.GitHubToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	case forgedomain.ForgeTypeGitLab:
		if gitLabConnectorType, hasGitLabConnectorType := config.NormalConfig.GitLabConnectorType.Get(); hasGitLabConnectorType {
			switch gitLabConnectorType {
			case forgedomain.GitLabConnectorTypeAPI:
				connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
					APIToken:  config.NormalConfig.GitLabToken,
					Log:       log,
					RemoteURL: remoteURL,
				})
				return Some(connector), err
			case forgedomain.GitLabConnectorTypeGlab:
				connector = glab.Connector{
					Backend:  backend,
					Frontend: frontend,
				}
				return Some(connector), err
			}
		}
		// no GitLabConnectorType specified --> use the API connector
		connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  config.NormalConfig.GitLabToken,
			Log:       log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	}
	return None[forgedomain.Connector](), nil
}
