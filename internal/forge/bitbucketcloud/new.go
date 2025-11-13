package bitbucketcloud

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

// Detect indicates whether the current repository is hosted on Bitbucket Cloud.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "bitbucket.org"
}

type NewConnectorArgs struct {
	AppPassword      Option[forgedomain.BitbucketAppPassword]
	Log              print.Logger
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
	UserName         Option[forgedomain.BitbucketUsername]
}

// NewConnector provides the correct connector for talking to Bitbucket Cloud.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint: ireturn
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
	userName, hasUserName := args.UserName.Get()
	appPassword, hasAppPassword := args.AppPassword.Get()
	if hasUserName && hasAppPassword {
		apiConnector := APIConnector{
			WebConnector: webConnector,
			client:       NewMutable(bitbucket.NewBasicAuth(userName.String(), appPassword.String())),
			log:          args.Log,
		}
		return &CachedAPIConnector{
			api:   apiConnector,
			cache: forgedomain.APICache{},
		}
	}
	return webConnector
}
