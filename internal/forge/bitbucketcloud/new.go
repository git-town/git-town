package bitbucketcloud

import (
	"github.com/git-town/git-town/v24/internal/browser/browserdomain"
	"github.com/git-town/git-town/v24/internal/cli/print"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/giturl"
	"github.com/git-town/git-town/v24/internal/subshell"
	"github.com/git-town/git-town/v24/internal/test/mockproposals"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

const (
	// sshHostname is the hostname that Bitbucket Cloud uses for Git access over SSH.
	sshHostname = "ssh.bitbucket.org"
	// webHostname is the hostname of the Bitbucket Cloud website.
	webHostname = "bitbucket.org"
)

// Detect indicates whether the current repository is hosted on Bitbucket Cloud.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == webHostname || remoteURL.Host == sshHostname
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

// NewConnector provides the correct connector for talking to Bitbucket Cloud.
func NewConnector(args NewConnectorArgs) (forgedomain.Connector, error) { //nolint: ireturn
	webConnector := WebConnector{
		browserEnabled:    args.BrowserEnabled,
		browserExecutable: args.BrowserExecutable,
		Hostname:          webHostnameFor(args.RemoteURL.Host),
		Organization:      args.RemoteURL.Org,
		Repository:        args.RemoteURL.Repo,
	}
	if subshell.IsInTest() {
		proposalsPath := mockproposals.NewMockProposalPath(args.ConfigDir)
		proposals := mockproposals.Load(proposalsPath)
		return &MockConnector{
			cache:         forgedomain.APICache{},
			log:           args.Log,
			Proposals:     proposals,
			ProposalsPath: proposalsPath,
			WebConnector:  webConnector,
		}, nil
	}
	userName, hasUserName := args.UserName.Get()
	apiToken, hasAPIToken := args.APIToken.Get()
	if hasUserName && hasAPIToken {
		client, err := bitbucket.NewBasicAuth(userName.String(), apiToken.String())
		if err != nil {
			return nil, err
		}
		apiConnector := APIConnector{
			client:       NewMutable(client),
			log:          args.Log,
			WebConnector: webConnector,
		}
		return &CachedAPIConnector{
			api:   apiConnector,
			cache: forgedomain.APICache{},
		}, nil
	}
	return webConnector, nil
}

// webHostnameFor provides the hostname of the website for the given remote hostname.
// Bitbucket Cloud serves Git over SSH at a different hostname than its website.
func webHostnameFor(remoteHost string) string {
	if remoteHost == sshHostname {
		return webHostname
	}
	return remoteHost
}
