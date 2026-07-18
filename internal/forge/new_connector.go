package forge

import (
	"github.com/git-town/git-town/v23/internal/browser/browserdomain"
	"github.com/git-town/git-town/v23/internal/cli/print"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/azuredevops"
	"github.com/git-town/git-town/v23/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v23/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/forge/forgejo"
	"github.com/git-town/git-town/v23/internal/forge/gh"
	"github.com/git-town/git-town/v23/internal/forge/gitea"
	"github.com/git-town/git-town/v23/internal/forge/github"
	"github.com/git-town/git-town/v23/internal/forge/gitlab"
	"github.com/git-town/git-town/v23/internal/forge/glab"
	"github.com/git-town/git-town/v23/internal/git/giturl"
	"github.com/git-town/git-town/v23/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// NewConnector provides an instance of the forge connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (Option[forgedomain.Connector], Option[forgedomain.DetectedForgeType], error) {
	remoteURL, hasRemoteURL := args.RemoteURL.Get()
	detectedForgeType, hasForgeType := Detect(remoteURL, args.ForgeType).Get()
	if !hasRemoteURL || !hasForgeType {
		return None[forgedomain.Connector](), None[forgedomain.DetectedForgeType](), nil
	}
	var connector forgedomain.Connector
	var err error
	switch detectedForgeType.ForgeType() {
	case forgedomain.ForgeTypeAzuredevops:
		connector = azuredevops.WebConnector{
			HostedRepoInfo: forgedomain.HostedRepoInfo{
				Hostname:     remoteURL.Host,
				Organization: remoteURL.Org,
				Repository:   remoteURL.Repo,
			},
			BrowserEnabled:    args.BrowserEnabled,
			BrowserExecutable: args.BrowserExecutable,
		}
	case forgedomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			APIToken:          args.BitbucketAPIToken,
			BrowserEnabled:    args.BrowserEnabled,
			BrowserExecutable: args.BrowserExecutable,
			ConfigDir:         args.ConfigDir,
			Log:               args.Log,
			RemoteURL:         remoteURL,
			UserName:          args.BitbucketUsername,
		})
	case forgedomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			APIToken:          args.BitbucketAPIToken,
			BrowserEnabled:    args.BrowserEnabled,
			BrowserExecutable: args.BrowserExecutable,
			ConfigDir:         args.ConfigDir,
			Log:               args.Log,
			RemoteURL:         remoteURL,
			UserName:          args.BitbucketUsername,
		})
	case forgedomain.ForgeTypeForgejo:
		connector = forgejo.NewConnector(forgejo.NewConnectorArgs{
			APIToken:          args.ForgejoToken,
			BrowserEnabled:    args.BrowserEnabled,
			BrowserExecutable: args.BrowserExecutable,
			ConfigDir:         args.ConfigDir,
			Log:               args.Log,
			RemoteURL:         remoteURL,
		})
	case forgedomain.ForgeTypeGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:          args.GiteaToken,
			BrowserEnabled:    args.BrowserEnabled,
			BrowserExecutable: args.BrowserExecutable,
			ConfigDir:         args.ConfigDir,
			Log:               args.Log,
			RemoteURL:         remoteURL,
		})
	case forgedomain.ForgeTypeGithub:
		switch args.GithubConnectorType.GetOr(forgedomain.GithubConnectorTypeAPI) {
		case forgedomain.GithubConnectorTypeAPI:
			connector, err = github.NewConnector(github.NewConnectorArgs{
				APIToken:          args.GithubToken,
				BrowserEnabled:    args.BrowserEnabled,
				BrowserExecutable: args.BrowserExecutable,
				ConfigDir:         args.ConfigDir,
				Log:               args.Log,
				RemoteURL:         remoteURL,
			})
		case forgedomain.GithubConnectorTypeGh:
			connector = &gh.CachedConnector{
				Connector: gh.Connector{
					Backend:        args.Backend,
					BrowserEnabled: args.BrowserEnabled,
					Frontend:       args.Frontend,
					Log:            args.Log,
				},
				Cache: forgedomain.APICache{},
			}
		}
	case forgedomain.ForgeTypeGitlab:
		switch args.GitlabConnectorType.GetOr(forgedomain.GitlabConnectorTypeAPI) {
		case forgedomain.GitlabConnectorTypeAPI:
			connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
				APIToken:          args.GitlabToken,
				BrowserEnabled:    args.BrowserEnabled,
				BrowserExecutable: args.BrowserExecutable,
				ConfigDir:         args.ConfigDir,
				Log:               args.Log,
				RemoteURL:         remoteURL,
			})
		case forgedomain.GitlabConnectorTypeGlab:
			connector = &glab.CachedConnector{
				Connector: glab.Connector{
					Backend:        args.Backend,
					BrowserEnabled: args.BrowserEnabled,
					Frontend:       args.Frontend,
					Log:            args.Log,
				},
				Cache: forgedomain.APICache{},
			}
		}
	}
	return NewOption(connector), Some(detectedForgeType), err
}

type NewConnectorArgs struct {
	Backend             subshelldomain.Querier
	BitbucketAPIToken   Option[forgedomain.BitbucketAPIToken]
	BitbucketUsername   Option[forgedomain.BitbucketUsername]
	BrowserEnabled      browserdomain.BrowserEnabled
	BrowserExecutable   Option[browserdomain.BrowserExecutable]
	ConfigDir           configdomain.RepoConfigDir
	ForgeType           Option[forgedomain.ForgeType]
	ForgejoToken        Option[forgedomain.ForgejoToken]
	Frontend            subshelldomain.Runner
	GiteaToken          Option[forgedomain.GiteaToken]
	GithubConnectorType Option[forgedomain.GithubConnectorType]
	GithubToken         Option[forgedomain.GithubToken]
	GitlabConnectorType Option[forgedomain.GitlabConnectorType]
	GitlabToken         Option[forgedomain.GitlabToken]
	Log                 print.Logger
	RemoteURL           Option[giturl.Parts]
}
