package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v24/internal/browser/browserdomain"
	"github.com/git-town/git-town/v24/internal/cli/print"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/giturl"
	"github.com/git-town/git-town/v24/internal/subshell"
	"github.com/git-town/git-town/v24/internal/test/mockproposals"
	. "github.com/git-town/git-town/v24/pkg/prelude"
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
		browserEnabled:    args.BrowserEnabled,
		browserExecutable: args.BrowserExecutable,
	}
	if subshell.IsInTest() {
		proposalsPath := mockproposals.NewMockProposalPath(args.ConfigDir)
		proposals := mockproposals.Load(proposalsPath)
		return &MockConnector{
			Proposals:     proposals,
			ProposalsPath: proposalsPath,
			WebConnector:  webConnector,
			cache:         forgedomain.APICache{},
			log:           args.Log,
		}
	}
	if args.UserName.IsSome() && args.APIToken.IsSome() {
		apiConnector := APIConnector{
			WebConnector: webConnector,
			log:          args.Log,
			token:        args.APIToken.GetOrZero().String(),
			username:     args.UserName.GetOrZero().String(),
		}
		return &CachedAPIConnector{
			api:   apiConnector,
			cache: forgedomain.APICache{},
		}
	}
	return webConnector
}

type NewConnectorArgs struct {
	APIToken          Option[forgedomain.BitbucketAPIToken]
	BrowserEnabled    browserdomain.BrowserEnabled
	BrowserExecutable Option[browserdomain.BrowserExecutable]
	ConfigDir         configdomain.RepoConfigDir
	Log               print.Logger
	RemoteURL         giturl.Parts
	UserName          Option[forgedomain.BitbucketUsername]
}
