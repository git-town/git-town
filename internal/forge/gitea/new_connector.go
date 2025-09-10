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

type NewConnectorArgs struct {
	APIToken  Option[forgedomain.GiteaToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

// NewGiteaConfig provides Gitea configuration data if the current repo is hosted on Gitea,
// otherwise nil.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint:ireturn
	anonConnector := AnonConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	apiToken, hasAPIToken := args.APIToken.Get()
	if !hasAPIToken {
		return anonConnector
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient := gitea.NewClientWithHTTP("https://"+args.RemoteURL.Host, httpClient)
	return AuthConnector{
		APIToken:      args.APIToken,
		AnonConnector: anonConnector,
		client:        giteaClient,
		log:           args.Log,
	}
}
