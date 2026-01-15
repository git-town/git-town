package bitbucketcloud

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

// Detect indicates whether the current repository is hosted on Bitbucket Cloud.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "bitbucket.org"
}

type NewConnectorArgs struct {
	AppPassword Option[forgedomain.BitbucketAppPassword]
	Browser     Option[configdomain.Browser]
	ConfigDir   configdomain.RepoConfigDir
	Log         print.Logger
	RemoteURL   giturl.Parts
	UserName    Option[forgedomain.BitbucketUsername]
}

// NewConnector provides the correct connector for talking to Bitbucket Cloud.
func NewConnector(args NewConnectorArgs) forgedomain.Connector { //nolint: ireturn
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
