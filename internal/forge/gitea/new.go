package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Detect indicates whether the current repository is hosted on a gitea server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "gitea.com"
}

type NewConnectorArgs struct {
	APIToken         Option[forgedomain.GiteaToken]
	Log              print.Logger
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
}

// NewConnector provides a connector instance that talks to the gitea API.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint:ireturn
	webConnector := WebConnector{
		HostedRepoInfo: forgedomain.HostedRepoInfo{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	if proposalURLOverride, hasProposalOverride := args.ProposalOverride.Get(); hasProposalOverride {
		return TestConnector{
			WebConnector: webConnector,
			log:          args.Log,
			override:     proposalURLOverride,
		}
	}
	if apiToken, hasAPIToken := args.APIToken.Get(); hasAPIToken {
		authConnector := AuthConnector{
			APIToken:     apiToken,
			RemoteURL:    args.RemoteURL,
			WebConnector: webConnector,
			_client:      MutableNone[gitea.Client](),
			log:          args.Log,
		}
		return &CachedAPIConnector{
			api:   &authConnector,
			cache: forgedomain.ProposalCache{},
		}
	}
	return webConnector
}
