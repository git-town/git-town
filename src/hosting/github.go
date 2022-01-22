package hosting

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// GithubDriver provides access to the API of GitHub installations.
type GithubDriver struct {
	apiToken   string
	client     *github.Client
	config     config
	hostname   string
	log        logFn
	originURL  string
	owner      string
	repository string
}

// NewGithubDriver provides a GitHub driver instance if the given repo configuration is for a Github repo,
// otherwise nil.
func NewGithubDriver(config config, log logFn) *GithubDriver {
	driverType := config.HostingService()
	originURL := config.OriginURL()
	hostname := giturl.Host(originURL)
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		hostname = manualHostName
	}
	if driverType != "github" && hostname != "github.com" {
		return nil
	}
	repositoryParts := strings.SplitN(giturl.Repo(originURL), "/", 2)
	if len(repositoryParts) != 2 {
		return nil
	}
	owner := repositoryParts[0]
	repository := repositoryParts[1]
	return &GithubDriver{
		apiToken:   config.GitHubToken(),
		config:     config,
		hostname:   hostname,
		log:        log,
		originURL:  originURL,
		owner:      owner,
		repository: repository,
	}
}

func (d *GithubDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	if d.apiToken == "" {
		return result, nil
	}
	d.connect()
	pullRequests, err := d.loadPullRequests(branch, parentBranch)
	if err != nil {
		return result, err
	}
	if len(pullRequests) != 1 {
		return result, nil
	}
	result.CanMergeWithAPI = true
	result.DefaultCommitMessage = d.defaultCommitMessage(pullRequests[0])
	result.PullRequestNumber = int64(pullRequests[0].GetNumber())
	return result, nil
}

func (d *GithubDriver) NewPullRequestURL(branch string, parentBranch string) (string, error) {
	toCompare := branch
	if parentBranch != d.config.MainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", d.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (d *GithubDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", d.hostname, d.owner, d.repository)
}

func (d *GithubDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	d.connect()
	err = d.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	return d.mergePullRequest(options)
}

func (d *GithubDriver) HostingServiceName() string {
	return "GitHub"
}

// Helper

func (d *GithubDriver) connect() {
	if d.client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: d.apiToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		d.client = github.NewClient(tc)
	}
}

func (d *GithubDriver) defaultCommitMessage(pullRequest *github.PullRequest) string {
	return fmt.Sprintf("%s (#%d)", *pullRequest.Title, *pullRequest.Number)
}

func (d *GithubDriver) loadPullRequests(branch, parentBranch string) ([]*github.PullRequest, error) {
	pullRequests, _, err := d.client.PullRequests.List(context.Background(), d.owner, d.repository, &github.PullRequestListOptions{
		Base:  parentBranch,
		Head:  d.owner + ":" + branch,
		State: "open",
	})
	return pullRequests, err
}

func (d *GithubDriver) mergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	if options.PullRequestNumber == 0 {
		return "", fmt.Errorf("cannot merge via Github since there is no pull request")
	}
	if options.LogRequests {
		d.log("GitHub API: Merging PR #%d\n", options.PullRequestNumber)
	}
	commitMessageParts := strings.SplitN(options.CommitMessage, "\n", 2)
	githubCommitTitle := commitMessageParts[0]
	githubCommitMessage := ""
	if len(commitMessageParts) == 2 {
		githubCommitMessage = commitMessageParts[1]
	}
	result, _, err := d.client.PullRequests.Merge(context.Background(), d.owner, d.repository, int(options.PullRequestNumber), githubCommitMessage, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: githubCommitTitle,
	})
	if err != nil {
		return "", err
	}
	return *result.SHA, nil
}

func (d *GithubDriver) updatePullRequestsAgainst(options MergePullRequestOptions) error {
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
