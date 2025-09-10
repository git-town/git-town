package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Detect always return false because we can't guess a self-hosted URL.
func Detect(_ giturl.Parts) bool {
	return false
}

// NewConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint:ireturn
	webConnector := WebConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	hasAuth := args.UserName.IsSome() && args.AppPassword.IsSome()
	if !hasAuth {
		return webConnector
	}
	return APIConnector{
		WebConnector: webConnector,
		log:          args.Log,
		token:        args.AppPassword.String(),
		username:     args.UserName.String(),
	}
}

type NewConnectorArgs struct {
	AppPassword Option[forgedomain.BitbucketAppPassword]
	ForgeType   Option[forgedomain.ForgeType]
	Log         print.Logger
	RemoteURL   giturl.Parts
	UserName    Option[forgedomain.BitbucketUsername]
}
