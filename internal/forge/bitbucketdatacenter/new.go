package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Detect always return false because we can't guess a self-hosted URL.
func Detect(_ giturl.Parts) bool {
	return false
}

// NewConnector provides the Bitbucket connector instance to use.
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
	if args.UserName.IsSome() && args.AppPassword.IsSome() {
		apiConnector := APIConnector{
			WebConnector: webConnector,
			log:          args.Log,
			token:        args.AppPassword.GetOrZero().String(),
			username:     args.UserName.GetOrZero().String(),
		}
		return &CachedAPIConnector{
			api:   apiConnector,
			cache: forgedomain.ProposalCache{},
		}
	}
	return webConnector
}

type NewConnectorArgs struct {
	AppPassword      Option[forgedomain.BitbucketAppPassword]
	Log              print.Logger
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
	UserName         Option[forgedomain.BitbucketUsername]
}
