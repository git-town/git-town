package forgejo

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Detect indicates whether the current repository is hosted on a Forgejo server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "codeberg.org"
}

type NewConnectorArgs struct {
	APIToken  Option[forgedomain.ForgejoToken]
	Browser   Option[configdomain.Browser]
	ConfigDir configdomain.RepoConfigDir
	Log       print.Logger
	RemoteURL giturl.Parts
}

// NewConnector provides a new connector instance for the Forgejo API.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint:ireturn
	webConnector := WebConnector{
		HostedRepoInfo: forgedomain.HostedRepoInfo{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		browser: args.Browser,
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
	if args.APIToken.IsSome() {
		apiConnector := APIConnector{
			APIToken:     args.APIToken,
			WebConnector: webConnector,
			_client:      MutableNone[forgejo.Client](),
			log:          args.Log,
			remoteURL:    args.RemoteURL,
		}
		return &CachedAPIConnector{
			api:   &apiConnector,
			cache: forgedomain.APICache{},
		}
	}
	return webConnector
}
