package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
	anonConnector := WebConnector{
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	if proposalURLOverride, hasProposalOverride := args.ProposalOverride.Get(); hasProposalOverride {
		return TestConnector{
			WebConnector: anonConnector,
			log:          args.Log,
			override:     proposalURLOverride,
		}, nil
	}
	apiToken, hasAPIToken := args.APIToken.Get()
	if !hasAPIToken {
		return anonConnector, nil
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	githubClient := github.NewClient(httpClient)
	if args.RemoteURL.Host != "github.com" {
		url := "https://" + args.RemoteURL.Host
		var err error
		githubClient, err = githubClient.WithEnterpriseURLs(url, url)
		if err != nil {
			return anonConnector, fmt.Errorf(messages.GitHubEnterpriseInitializeError, err)
		}
	}
	return APIConnector{
		APIToken:     args.APIToken,
		WebConnector: anonConnector,
		client:       NewMutable(githubClient),
		log:          args.Log,
	}, nil
}

func RepositoryURL(hostNameWithStandardPort string, organization string, repository string) string {
	return fmt.Sprintf("https://%s/%s/%s", hostNameWithStandardPort, organization, repository)
}

// parsePullRequest extracts standardized proposal data from the given GitHub pull-request.
func parsePullRequest(pullRequest *github.PullRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		Body:         NewOption(pullRequest.GetBody()),
		Number:       pullRequest.GetNumber(),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.GetRef()),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.GetRef()),
		Title:        pullRequest.GetTitle(),
		MergeWithAPI: pullRequest.GetMergeableState() == "clean",
		URL:          *pullRequest.HTMLURL,
	}
}
