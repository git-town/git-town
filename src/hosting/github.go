package hosting

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
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
	MainBranch string
	log        Log
}

func (c *GitHubConnector) FindProposal(branch, target string) (*Proposal, error) {
	pullRequests, _, err := c.client.PullRequests.List(context.Background(), c.Organization, c.Repository, &github.PullRequestListOptions{
		Head:  c.Organization + ":" + branch,
		Base:  target,
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

func (c *GitHubConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (c *GitHubConnector) HostingServiceName() string {
	return "GitHub"
}

func (c *GitHubConnector) NewProposalURL(branch, parentBranch string) (string, error) {
	toCompare := branch
	if parentBranch != c.MainBranch {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c *GitHubConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.Hostname, c.Organization, c.Repository)
}

func (c *GitHubConnector) SquashMergeProposal(number int, message string) (mergeSHA git.SHA, err error) {
	if number <= 0 {
		return git.SHA{}, fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	c.log.Start(messages.HostingGithubMergingViaAPI, number)
	title, body := ParseCommitMessage(message)
	result, _, err := c.client.PullRequests.Merge(context.Background(), c.Organization, c.Repository, number, body, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: title,
	})
	sha := git.NewSHA(result.GetSHA())
	if err != nil {
		c.log.Failed(err)
		return sha, err
	}
	c.log.Success()
	return sha, nil
}

func (c *GitHubConnector) UpdateProposalTarget(number int, target string) error {
	c.log.Start(messages.HostingGithubUpdatePRViaAPI, number)
	_, _, err := c.client.PullRequests.Edit(context.Background(), c.Organization, c.Repository, number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &target,
		},
	})
	if err != nil {
		c.log.Failed(err)
		return err
	}
	c.log.Success()
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
	MainBranch     string
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
		Target:          pullRequest.Base.GetRef(),
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
