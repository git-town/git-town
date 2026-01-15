package gitlab

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Detect indicates whether the current repository is hosted on a GitLab server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "gitlab.com"
}

type NewConnectorArgs struct {
	APIToken  Option[forgedomain.GitlabToken]
	Browser   Option[configdomain.Browser]
	ConfigDir configdomain.RepoConfigDir
	Log       print.Logger
	RemoteURL giturl.Parts
}

func NewConnector(args NewConnectorArgs) (forgedomain.Connector, error) { //nolint: ireturn
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
		}, nil
	}
	if apiToken, hasAPIToken := args.APIToken.Get(); hasAPIToken {
		client, err := gitlab.NewClient(apiToken.String(), gitlab.WithBaseURL(webConnector.baseURL()))
		if err != nil {
			return webConnector, err
		}
		apiConnector := APIConnector{
			WebConnector: webConnector,
			client:       client,
			log:          args.Log,
		}
		return &CachedAPIConnector{
			api:   apiConnector,
			cache: forgedomain.APICache{},
		}, nil
	}
	return webConnector, nil
}
