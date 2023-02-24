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

// GithubConfig contains connection information for Github-based hosting platforms.
type GithubConfig struct {
	apiToken   string
	config     config
	hostname   string
	originURL  string
	owner      string
	repository string
}

// NewGithubConfig provides GitHub configuration data if the current repo is hosted on Github,
// otherwise nil.
func NewGithubConfig(url giturl.Parts, config config) *GithubConfig {
	driverType := config.HostingService()
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if driverType != "github" && url.Host != "github.com" {
		return nil
	}
	return &GithubConfig{
		apiToken:   config.GitHubToken(),
		config:     config,
		hostname:   url.Host,
		originURL:  config.OriginURL(),
		owner:      url.Org,
		repository: url.Repo,
	}
}

func (c GithubConfig) defaultProposalMessage(pullRequest *github.PullRequest) string {
	return fmt.Sprintf("%s (#%d)", *pullRequest.Title, *pullRequest.Number)
}

// Driver provides a working driver instance ready to talk to the GitHub API.
func (c GithubConfig) Driver(log logFn) GithubDriver {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.apiToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return GithubDriver{
		client:       github.NewClient(httpClient),
		GithubConfig: c,
		log:          log,
	}
}

func (c GithubConfig) HostingServiceName() string {
	return "GitHub"
}

func (c GithubConfig) NewProposalURL(branch string, parentBranch string) (string, error) {
	toCompare := branch
	if parentBranch != c.config.MainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c GithubConfig) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.hostname, c.owner, c.repository)
}

/*************************************************************************
 *
 *  GITHUB DRIVER
 *
 */

// GithubDriver provides access to the GitHub API.
type GithubDriver struct {
	GithubConfig
	client *github.Client
	log    logFn
}

func (d *GithubDriver) ProposalDetails(branch, parentBranch string) (*Proposal, error) {
	if d.apiToken == "" {
		return nil, nil //nolint:nilnil // we really want to return nil here
	}
	pullRequests, err := d.loadPullRequests(branch, parentBranch)
	if err != nil {
		return nil, err
	}
	if len(pullRequests) < 1 {
		return nil, fmt.Errorf("no pull request from branch %q to branch %q found", branch, parentBranch)
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf("found %d pull requests from branch %q to branch %q", len(pullRequests), branch, parentBranch)
	}
	return &Proposal{
		CanMergeWithAPI:        true,
		DefaultProposalMessage: d.defaultProposalMessage(pullRequests[0]),
		ProposalNumber:         pullRequests[0].GetNumber(),
	}, nil
}

func (d *GithubDriver) loadPullRequests(branch, parentBranch string) ([]*github.PullRequest, error) {
	pullRequests, _, err := d.client.PullRequests.List(context.Background(), d.owner, d.repository, &github.PullRequestListOptions{
		Base:  parentBranch,
		Head:  d.owner + ":" + branch,
		State: "open",
	})
	return pullRequests, err
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (d *GithubDriver) SquashMergeProposal(options SquashMergeProposalOptions) (mergeSha string, err error) {
	err = d.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	if options.ProposalNumber == 0 {
		return "", fmt.Errorf("cannot merge via Github since there is no pull request")
	}
	if options.LogRequests {
		d.log("GitHub API: Merging PR #%d\n", options.ProposalNumber)
	}
	commitMessageParts := strings.SplitN(options.CommitMessage, "\n", 2)
	githubCommitTitle := commitMessageParts[0]
	githubCommitMessage := ""
	if len(commitMessageParts) == 2 {
		githubCommitMessage = commitMessageParts[1]
	}
	result, _, err := d.client.PullRequests.Merge(context.Background(), d.owner, d.repository, options.ProposalNumber, githubCommitMessage, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: githubCommitTitle,
	})
	if err != nil {
		return "", err
	}
	return *result.SHA, nil
}

func (d *GithubDriver) updatePullRequestsAgainst(options SquashMergeProposalOptions) error {
	pullRequests, _, err := d.client.PullRequests.List(context.Background(), d.owner, d.repository, &github.PullRequestListOptions{
		Base:  options.Branch,
		State: "open",
	})
	if err != nil {
		return err
	}
	for _, pullRequest := range pullRequests {
		if options.LogRequests {
			d.log("GitHub API: Updating base branch for PR #%d\n", *pullRequest.Number)
		}
		_, _, err = d.client.PullRequests.Edit(context.Background(), d.owner, d.repository, *pullRequest.Number, &github.PullRequest{
			Base: &github.PullRequestBranch{
				Ref: &options.ParentBranch,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
