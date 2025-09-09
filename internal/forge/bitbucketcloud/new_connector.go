package bitbucketcloud

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

// There are two types of connectors to Bitbucket Cloud:
// - WebConnector is used when the user hasn't configured API credentials
// - APIConnector is used when the user has configured API credentials

type NewConnectorArgs struct {
	AppPassword Option[forgedomain.BitbucketAppPassword]
	ForgeType   Option[forgedomain.ForgeType]
	Log         print.Logger
	RemoteURL   giturl.Parts
	UserName    Option[forgedomain.BitbucketUsername]
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) forgedomain.Connector {
	webConnector := WebConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	if args.UserName.IsSome() && args.AppPassword.IsSome() {
		client := bitbucket.NewBasicAuth(args.UserName.String(), args.AppPassword.String())
		return APIConnector{
			WebConnector: webConnector,
			client:       client,
			log:          args.Log,
		}
	}
	return webConnector
}
