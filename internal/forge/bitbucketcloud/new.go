package bitbucketcloud

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

// Detect indicates whether the current repository is hosted on a Bitbucket server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "bitbucket.org"
}

type NewConnectorArgs struct {
	AppPassword      Option[forgedomain.BitbucketAppPassword]
	ForgeType        Option[forgedomain.ForgeType]
	Log              print.Logger
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
	UserName         Option[forgedomain.BitbucketUsername]
}

// NewConnector provides the correct connector for talking to Bitbucket Cloud.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint: ireturn
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
	userName, hasUserName := args.UserName.Get()
	appPassword, hasAppPassword := args.AppPassword.Get()
	hasAuth := hasUserName && hasAppPassword
	if hasAuth {
		return APIConnector{
			WebConnector: webConnector,
			client:       NewMutable(bitbucket.NewBasicAuth(userName.String(), appPassword.String())),
			log:          args.Log,
		}
	}
	return webConnector
}
