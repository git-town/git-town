package forge

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/azuredevops"
	"github.com/git-town/git-town/v22/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v22/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/forgejo"
	"github.com/git-town/git-town/v22/internal/forge/gh"
	"github.com/git-town/git-town/v22/internal/forge/gitea"
	"github.com/git-town/git-town/v22/internal/forge/github"
	"github.com/git-town/git-town/v22/internal/forge/gitlab"
	"github.com/git-town/git-town/v22/internal/forge/glab"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// NewConnector provides an instance of the forge connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (Option[forgedomain.Connector], error) {
	remoteURL, hasRemoteURL := args.RemoteURL.Get()
	forgeType, hasForgeType := Detect(remoteURL, args.ForgeType).Get()
	if !hasRemoteURL || !hasForgeType {
		return None[forgedomain.Connector](), nil
	}
	proposalOverride := forgedomain.ReadProposalOverride()
	var connector forgedomain.Connector
	var err error
	switch forgeType {
	case forgedomain.ForgeTypeAzureDevOps:
		connector = azuredevops.NewConnector(azuredevops.NewConnectorArgs{
			Browser:          args.Browser,
			ProposalOverride: proposalOverride,
			RemoteURL:        remoteURL,
		})
	case forgedomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword:      args.BitbucketAppPassword,
			Browser:          args.Browser,
			Log:              args.Log,
			ProposalOverride: proposalOverride,
			RemoteURL:        remoteURL,
			UserName:         args.BitbucketUsername,
		})
	case forgedomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			AppPassword:      args.BitbucketAppPassword,
			Browser:          args.Browser,
			Log:              args.Log,
			ProposalOverride: proposalOverride,
			RemoteURL:        remoteURL,
			UserName:         args.BitbucketUsername,
		})
	case forgedomain.ForgeTypeForgejo:
		connector = forgejo.NewConnector(forgejo.NewConnectorArgs{
			APIToken:         args.ForgejoToken,
			Browser:          args.Browser,
			Log:              args.Log,
			ProposalOverride: proposalOverride,
			RemoteURL:        remoteURL,
		})
	case forgedomain.ForgeTypeGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:         args.GiteaToken,
			Browser:          args.Browser,
			Log:              args.Log,
			ProposalOverride: proposalOverride,
			RemoteURL:        remoteURL,
		})
	case forgedomain.ForgeTypeGitHub:
		if testHome, inTestMode := args.TestHome.Get(); inTestMode {
			connector = github.MockConnector{
				WebConnector: github.NewWebConnector(remoteURL, args.Browser),
				Proposals:    mockproposals.Load(testHome.String()),
			}
		} else {
		switch args.GitHubConnectorType.GetOr(forgedomain.GitHubConnectorTypeAPI) {
		case forgedomain.GitHubConnectorTypeAPI:
			connector, err = github.NewConnector(github.NewConnectorArgs{
				APIToken:         args.GitHubToken,
				Browser:          args.Browser,
				Log:              args.Log,
				ProposalOverride: proposalOverride,
				RemoteURL:        remoteURL,
			})
		case forgedomain.GitHubConnectorTypeGh:
			connector = &gh.CachedConnector{
				Connector: gh.Connector{
					Backend:  args.Backend,
					Frontend: args.Frontend,
				},
				Cache: forgedomain.APICache{},
			}
		}
	case forgedomain.ForgeTypeGitLab:
		switch args.GitLabConnectorType.GetOr(forgedomain.GitLabConnectorTypeAPI) {
		case forgedomain.GitLabConnectorTypeAPI:
			connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
				APIToken:         args.GitLabToken,
				Browser:          args.Browser,
				Log:              args.Log,
				ProposalOverride: proposalOverride,
				RemoteURL:        remoteURL,
			})
		case forgedomain.GitLabConnectorTypeGlab:
			connector = &glab.CachedConnector{
				Connector: glab.Connector{
					Backend:  args.Backend,
					Frontend: args.Frontend,
				},
				Cache: forgedomain.APICache{},
			}
		}
	}
	return NewOption(connector), err
}

type NewConnectorArgs struct {
	Backend              subshelldomain.Querier
	BitbucketAppPassword Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername    Option[forgedomain.BitbucketUsername]
	Browser              Option[configdomain.Browser]
	ForgeType            Option[forgedomain.ForgeType]
	ForgejoToken         Option[forgedomain.ForgejoToken]
	Frontend             subshelldomain.Runner
	GitHubConnectorType  Option[forgedomain.GitHubConnectorType]
	GitHubToken          Option[forgedomain.GitHubToken]
	GitLabConnectorType  Option[forgedomain.GitLabConnectorType]
	GitLabToken          Option[forgedomain.GitLabToken]
	GiteaToken           Option[forgedomain.GiteaToken]
	Log                  print.Logger
	RemoteURL            Option[giturl.Parts]
	TestHome             Option[configdomain.TestHome]
}
