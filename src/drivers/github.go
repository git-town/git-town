package drivers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/drivers/helpers"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// githubCodeHostingDriver provides access to the API of GitHub installations.
type githubCodeHostingDriver struct {
	apiToken   string
	client     *github.Client
	config     config
	hostname   string
	originURL  string
	owner      string
	repository string
	log        logFn
}

// LoadGithub provides a GitHub driver instance if the given repo configuration is for a Github repo,
// otherwise nil.
func LoadGithub(config config, log logFn) CodeHostingDriver {
	driverType := config.GetCodeHostingDriverName()
	originURL := config.GetRemoteOriginURL()
	hostname := helpers.GetURLHostname(originURL)
	configuredHostName := config.GetCodeHostingOriginHostname()
	if configuredHostName != "" {
		hostname = configuredHostName
	}
	if driverType != "github" && hostname != "github.com" {
		return nil
	}
	repositoryParts := strings.SplitN(helpers.GetURLRepositoryName(originURL), "/", 2)
	if len(repositoryParts) != 2 {
		return nil
	}
	owner := repositoryParts[0]
	repository := repositoryParts[1]
	return &githubCodeHostingDriver{
		apiToken:   config.GetGitHubToken(),
		config:     config,
		hostname:   hostname,
		log:        log,
		originURL:  originURL,
		owner:      owner,
		repository: repository,
	}
}

func (d *githubCodeHostingDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	if d.apiToken == "" {
		return result, nil
	}
	d.connect()
	pullRequests, err := d.getPullRequests(branch, parentBranch)
	if err != nil {
		return result, err
	}
	if len(pullRequests) != 1 {
		return result, nil
	}
	result.CanMergeWithAPI = true
	result.DefaultCommitMessage = d.getDefaultCommitMessage(pullRequests[0])
	result.PullRequestNumber = int64(pullRequests[0].GetNumber())
	return result, nil
}

func (d *githubCodeHostingDriver) NewPullRequestURL(branch string, parentBranch string) (string, error) {
	toCompare := branch
	if parentBranch != d.config.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", d.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (d *githubCodeHostingDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", d.hostname, d.owner, d.repository)
}

func (d *githubCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	d.connect()
	err = d.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	return d.mergePullRequest(options)
}

func (d *githubCodeHostingDriver) HostingServiceName() string {
	return "GitHub"
}

// Helper

func (d *githubCodeHostingDriver) connect() {
	if d.client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: d.apiToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		d.client = github.NewClient(tc)
	}
}

func (d *githubCodeHostingDriver) getDefaultCommitMessage(pullRequest *github.PullRequest) string {
	return fmt.Sprintf("%s (#%d)", *pullRequest.Title, *pullRequest.Number)
}

func (d *githubCodeHostingDriver) getPullRequests(branch, parentBranch string) ([]*github.PullRequest, error) {
	pullRequests, _, err := d.client.PullRequests.List(context.Background(), d.owner, d.repository, &github.PullRequestListOptions{
		Base:  parentBranch,
		Head:  d.owner + ":" + branch,
		State: "open",
	})
	return pullRequests, err
}

func (d *githubCodeHostingDriver) mergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
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

func (d *githubCodeHostingDriver) updatePullRequestsAgainst(options MergePullRequestOptions) error {
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
