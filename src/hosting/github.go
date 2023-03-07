package hosting

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/config"
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
		return nil, fmt.Errorf("found %d pull requests from branch %q into branch %q", len(pullRequests), branch, target)
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

//nolint:nonamedreturns
func (c *GitHubConnector) SquashMergeProposal(number int, message string) (mergeSHA string, err error) {
	if number <= 0 {
		return "", fmt.Errorf("no pull request number given")
	}
	if c.log != nil {
		c.log("GitHub API: merging PR #%d\n", number)
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
		c.log("GitHub API: updating base branch for PR #%d\n", number)
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
	apiToken := gitConfig.GitHubToken()
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
		MainBranch: gitConfig.MainBranch(),
		log:        log,
	}, nil
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

//nolint:nonamedreturns
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
