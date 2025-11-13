package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "github.com"
}

type NewConnectorArgs struct {
	APIToken         Option[forgedomain.GitHubToken]
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
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken.String()})
		httpClient := oauth2.NewClient(context.Background(), tokenSource)
		githubClient := github.NewClient(httpClient)
		if args.RemoteURL.Host != "github.com" {
			url := "https://" + args.RemoteURL.Host
			var err error
			githubClient, err = githubClient.WithEnterpriseURLs(url, url)
			if err != nil {
				return webConnector, fmt.Errorf(messages.GitHubEnterpriseInitializeError, err)
			}
		}
		apiConnector := APIConnector{
			WebConnector: webConnector,
			client:       NewMutable(githubClient),
			log:          args.Log,
		}
		return CachedAPIConnector{
			api:   apiConnector,
			cache: forgedomain.ProposalCache{},
		}, nil
	}
	return webConnector, nil
}

func RepositoryURL(hostNameWithStandardPort string, organization string, repository string) string {
	return fmt.Sprintf("https://%s/%s/%s", hostNameWithStandardPort, organization, repository)
}

// parsePullRequest extracts standardized proposal data from the given GitHub pull-request.
func parsePullRequest(pullRequest *github.PullRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		Active:       pullRequest.GetState() == "open",
		Body:         NewOption(pullRequest.GetBody()),
		Number:       pullRequest.GetNumber(),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.GetRef()),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.GetRef()),
		Title:        pullRequest.GetTitle(),
		MergeWithAPI: pullRequest.GetMergeableState() == "clean",
		URL:          *pullRequest.HTMLURL,
	}
}
