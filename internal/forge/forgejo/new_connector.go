package forgejo

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type NewConnectorArgs struct {
	APIToken  Option[forgedomain.ForgejoToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

// NewConnector provides a new connector instance.
func NewConnector(args NewConnectorArgs) (forgedomain.Connector, error) { //nolint:ireturn
	webConnector := AnonConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	hasAuth := args.APIToken.IsSome()
	if !hasAuth {
		return webConnector, nil
	}
	forgejoClient, err := forgejo.NewClient("https://"+args.RemoteURL.Host, forgejo.SetToken(args.APIToken.String()))
	return AuthConnector{
		APIToken:      args.APIToken,
		AnonConnector: webConnector,
		client:        forgejoClient,
		log:           args.Log,
	}, err
}
