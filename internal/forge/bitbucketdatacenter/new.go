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

// NewConnector provides the Bitbucket connector instance to use.
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
	AppPassword      Option[forgedomain.BitbucketAppPassword]
	ForgeType        Option[forgedomain.ForgeType]
	Log              print.Logger
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
	UserName         Option[forgedomain.BitbucketUsername]
}
