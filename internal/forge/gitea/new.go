package gitea

import (
	"context"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"golang.org/x/oauth2"
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
		Data: forgedomain.Data{
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
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken.String()})
		httpClient := oauth2.NewClient(context.Background(), tokenSource)
		giteaClient := gitea.NewClientWithHTTP("https://"+args.RemoteURL.Host, httpClient)
		return AuthConnector{
			APIToken:     args.APIToken,
			WebConnector: webConnector,
			client:       giteaClient,
			log:          args.Log,
		}
	}
	return webConnector
}
