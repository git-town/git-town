package hosting

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

// GitHubConnector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type GitHubConnector struct {
	client *github.Client
	Config
	mainBranch string
	log        logFn
}

func (c *GitHubConnector) ProposalDetails(number int) (*Proposal, error) {
	pullRequest, _, err := c.client.PullRequests.Get(context.Background(), c.owner, c.repository, number)
	return parseGithubPullRequest(pullRequest), err
}

func (c *GitHubConnector) ProposalForBranch(branch string) (*Proposal, error) {
	pullRequests, _, err := c.client.PullRequests.List(context.Background(), c.owner, c.repository, &github.PullRequestListOptions{
		Head:  c.owner + ":" + branch,
		State: "open",
	})
	if err != nil {
		return nil, err
	}
	if len(pullRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf("found %d pull requests for branch %q", len(pullRequests), branch)
	}
	return parseGithubPullRequest(pullRequests[0]), nil
}

func (c *GitHubConnector) Proposals() ([]Proposal, error) {
	pullRequests, _, err := c.client.PullRequests.List(context.Background(), c.owner, c.repository, &github.PullRequestListOptions{
		State: "open",
	})
	result := make([]Proposal, len(pullRequests))
	if err != nil {
		return result, err
	}
	for p := range pullRequests {
		result[p] = *parseGithubPullRequest(pullRequests[p])
	}
	return result, nil
}

func (c *GitHubConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (c *GitHubConnector) HostingServiceName() string {
	return "GitHub"
}

func (c *GitHubConnector) NewProposalURL(branch, parentBranch string) (string, error) {
	toCompare := branch
	if parentBranch != c.mainBranch {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c *GitHubConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.owner, c.repository)
}

//nolint:nonamedreturns
func (c *GitHubConnector) SquashMergeProposal(number int, message string) (mergeSHA string, err error) {
	if number == 0 {
		return "", fmt.Errorf("no pull request number given")
	}
	if c.log != nil {
		c.log("GitHub API: merging PR #%d\n", number)
	}
	title, body := parseCommitMessage(message)
	result, _, err := c.client.PullRequests.Merge(context.Background(), c.owner, c.repository, number, body, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: title,
	})
	return result.GetSHA(), err
}

func (c *GitHubConnector) UpdateProposalTarget(number int, target string) error {
	if c.log != nil {
		c.log("GitHub API: updating base branch for PR #%d\n", number)
	}
	_, _, err := c.client.PullRequests.Edit(context.Background(), c.owner, c.repository, number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &target,
		},
	})
	return err
}

// NewGithubConnector provides a fully configured GithubConnector instance
// if the current repo is hosted on Github, otherwise nil.
func NewGithubConnector(url giturl.Parts, gitConfig gitConfig, log logFn) *GitHubConnector {
	manualHostName := gitConfig.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if gitConfig.HostingService() != "github" && url.Host != "github.com" {
		return nil
	}
	apiToken := gitConfig.GitHubToken()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return &GitHubConnector{
		client: github.NewClient(httpClient),
		Config: Config{
			apiToken:   apiToken,
			hostname:   url.Host,
			originURL:  gitConfig.OriginURL(),
			owner:      url.Org,
			repository: url.Repo,
		},
		mainBranch: gitConfig.MainBranch(),
		log:        log,
	}
}

// parseGithubPullRequest extracts ChangeRequestInfo from the given GitHub pull-request data.
func parseGithubPullRequest(pullRequest *github.PullRequest) *Proposal {
	if pullRequest == nil {
		return nil
	}
	return &Proposal{
		Number:          pullRequest.GetNumber(),
		Title:           pullRequest.GetTitle(),
		CanMergeWithAPI: pullRequest.GetMergeableState() == "clean",
		// see https://docs.github.com/en/rest/guides/using-the-rest-api-to-interact-with-your-git-database?apiVersion=2022-11-28#checking-mergeability-of-pull-requests for details around the mergeability state
	}
}

//nolint:nonamedreturns
func parseCommitMessage(message string) (title, body string) {
	parts := strings.SplitN(message, "\n", 2)
	title = parts[0]
	if len(parts) == 2 {
		body = parts[1]
	} else {
		body = ""
	}
	return
}
