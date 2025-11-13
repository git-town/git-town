package gitlab

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Detect indicates whether the current repository is hosted on a GitLab server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "gitlab.com"
}

type NewConnectorArgs struct {
	APIToken         Option[forgedomain.GitLabToken]
	Log              print.Logger
	ProposalOverride Option[forgedomain.ProposalOverride]
	RemoteURL        giturl.Parts
}

func NewConnector(args NewConnectorArgs) (forgedomain.Connector, error) { //nolint: ireturn
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
			cache: forgedomain.ProposalCache{},
		}, nil
	}
	return webConnector, nil
}
