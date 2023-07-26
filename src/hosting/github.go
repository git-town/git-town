package hosting

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
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
	log        logFn
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

func (c *GitHubConnector) SquashMergeProposal(number int, message string) (mergeSHA string, err error) {
	if number <= 0 {
		return "", fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	if c.log != nil {
		c.log(messages.HostingGithubMergingViaAPI, number)
	}
	title, body := ParseCommitMessage(message)
	result, _, err := c.client.PullRequests.Merge(context.Background(), c.Organization, c.Repository, number, body, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: title,
	})
	return result.GetSHA(), err
}

func (c *GitHubConnector) UpdateProposalTarget(number int, target string) error {
	if c.log != nil {
		c.log(messages.HostingGithubUpdateBasebranchViaAPI, number)
	}
	_, _, err := c.client.PullRequests.Edit(context.Background(), c.Organization, c.Repository, number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &target,
		},
	})
	return err
}

// NewGithubConnector provides a fully configured GithubConnector instance
// if the current repo is hosted on Github, otherwise nil.
func NewGithubConnector(gitConfig gitTownConfig, log logFn) (*GitHubConnector, error) {
	hostingService, err := gitConfig.HostingService()
	if err != nil {
		return nil, err
	}
	url := gitConfig.OriginURL()
	if url == nil || (url.Host != "github.com" && hostingService != config.HostingServiceGitHub) {
		return nil, nil //nolint:nilnil
	}
	apiToken := getGitHubAPIToken(gitConfig)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return &GitHubConnector{
		client: github.NewClient(httpClient),
		CommonConfig: CommonConfig{
			APIToken:     apiToken,
			Hostname:     url.Host,
			Organization: url.Org,
			Repository:   url.Repo,
		},
		MainBranch: gitConfig.MainBranch(), // TODO: inject mainBranch as argument
		log:        log,
	}, nil
}

// getGitHubApiToken returns the GitHub API token to use.
// It first checks the GITHUB_TOKEN environment variable.
// If that is not set, it checks the GITHUB_AUTH_TOKEN environment variable.
// If that is not set, it checks the git config.
func getGitHubAPIToken(gitConfig gitTownConfig) string {
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
