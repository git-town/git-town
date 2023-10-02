package hosting

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

// GitHubConnector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type GitHubConnector struct {
	client *github.Client
	CommonConfig
	MainBranch domain.LocalBranchName
	log        Log
}

func (gc *GitHubConnector) FindProposal(branch, target domain.LocalBranchName) (*Proposal, error) {
	pullRequests, _, err := gc.client.PullRequests.List(context.Background(), gc.Organization, gc.Repository, &github.PullRequestListOptions{
		Head:  gc.Organization + ":" + branch.String(),
		Base:  target.String(),
		State: "open",
	})
	if err != nil {
		return nil, err
	}
	if len(pullRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf(messages.ProposalMultipleFound, len(pullRequests), branch, target)
	}
	proposal := parsePullRequest(pullRequests[0])
	return &proposal, nil
}

func (gc *GitHubConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (gc *GitHubConnector) HostingServiceName() string {
	return "GitHub"
}

func (gc *GitHubConnector) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	toCompare := branch.String()
	if parentBranch != gc.MainBranch {
		toCompare = parentBranch.String() + "..." + branch.String()
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", gc.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (gc *GitHubConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", gc.Hostname, gc.Organization, gc.Repository)
}

func (gc *GitHubConnector) SquashMergeProposal(number int, message string) (mergeSHA domain.SHA, err error) {
	if number <= 0 {
		return domain.SHA{}, fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	gc.log.Start(messages.HostingGithubMergingViaAPI, number)
	title, body := ParseCommitMessage(message)
	result, _, err := gc.client.PullRequests.Merge(context.Background(), gc.Organization, gc.Repository, number, body, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: title,
	})
	sha := domain.NewSHA(result.GetSHA())
	if err != nil {
		gc.log.Failed(err)
		return sha, err
	}
	gc.log.Success()
	return sha, nil
}

func (gc *GitHubConnector) UpdateProposalTarget(number int, target domain.LocalBranchName) error {
	gc.log.Start(messages.HostingGithubUpdatePRViaAPI, number)
	targetName := target.String()
	_, _, err := gc.client.PullRequests.Edit(context.Background(), gc.Organization, gc.Repository, number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &(targetName),
		},
	})
	if err != nil {
		gc.log.Failed(err)
		return err
	}
	gc.log.Success()
	return nil
}

// NewGithubConnector provides a fully configured GithubConnector instance
// if the current repo is hosted on Github, otherwise nil.
func NewGithubConnector(args NewGithubConnectorArgs) (*GitHubConnector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "github.com" && args.HostingService != config.HostingGitHub) {
		return nil, nil //nolint:nilnil
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return &GitHubConnector{
		client: github.NewClient(httpClient),
		CommonConfig: CommonConfig{
			APIToken:     args.APIToken,
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		MainBranch: args.MainBranch,
		log:        args.Log,
	}, nil
}

type NewGithubConnectorArgs struct {
	HostingService config.Hosting
	OriginURL      *giturl.Parts
	APIToken       string
	MainBranch     domain.LocalBranchName
	Log            Log
}

// getGitHubApiToken returns the GitHub API token to use.
// It first checks the GITHUB_TOKEN environment variable.
// If that is not set, it checks the GITHUB_AUTH_TOKEN environment variable.
// If that is not set, it checks the git config.
func GetGitHubAPIToken(gitConfig gitTownConfig) string {
	apiToken := os.ExpandEnv("$GITHUB_TOKEN")
	if apiToken == "" {
		apiToken = os.ExpandEnv("$GITHUB_AUTH_TOKEN")
	}
	if apiToken == "" {
		apiToken = gitConfig.GitHubToken()
	}
	return apiToken
}

// parsePullRequest extracts standardized proposal data from the given GitHub pull-request.
func parsePullRequest(pullRequest *github.PullRequest) Proposal {
	return Proposal{
		Number:          pullRequest.GetNumber(),
		Target:          domain.NewLocalBranchName(pullRequest.Base.GetRef()),
		Title:           pullRequest.GetTitle(),
		CanMergeWithAPI: pullRequest.GetMergeableState() == "clean",
	}
}

func ParseCommitMessage(message string) (title, body string) {
	parts := strings.SplitN(message, "\n", 2)
	title = parts[0]
	if len(parts) == 2 {
		body = parts[1]
	} else {
		body = ""
	}
	for strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return
}
