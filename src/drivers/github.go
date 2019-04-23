package drivers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"github.com/Originate/git-town/src/git"
	"github.com/google/go-github/github"
)

type githubCodeHostingDriver struct {
	originURL  string
	hostname   string
	apiToken   string
	client     *github.Client
	owner      string
	repository string
}

func (d *githubCodeHostingDriver) CanBeUsed(driverType string) bool {
	return driverType == "github" || d.hostname == "github.com"
}

func (d *githubCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, string, error) {
	if d.apiToken == "" {
		return false, "", nil
	}
	d.connect()
	pullRequests, err := d.getPullRequests(branch, parentBranch)
	if err != nil {
		return false, "", err
	}
	if len(pullRequests) != 1 {
		return false, "", nil
	}
	return true, d.getDefaultCommitMessage(pullRequests[0]), nil
}

func (d *githubCodeHostingDriver) GetNewPullRequestURL(branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", d.GetRepositoryURL(), url.PathEscape(toCompare))
}

func (d *githubCodeHostingDriver) GetRepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", d.hostname, d.owner, d.repository)
}

func (d *githubCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	d.connect()
	err := d.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	return d.mergePullRequest(options)
}

func (d *githubCodeHostingDriver) HostingServiceName() string {
	return "GitHub"
}

func (d *githubCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.GetURLHostname(originURL)
	d.client = nil
	repositoryParts := strings.SplitN(git.GetURLRepositoryName(originURL), "/", 2)
	if len(repositoryParts) == 2 {
		d.owner = repositoryParts[0]
		d.repository = repositoryParts[1]
	}
}

func (d *githubCodeHostingDriver) SetOriginHostname(originHostname string) {
	d.hostname = originHostname
}

func (d *githubCodeHostingDriver) GetAPITokenKey() string {
	return "git-town.github-token"
}

func (d *githubCodeHostingDriver) SetAPIToken(apiToken string) {
	d.apiToken = apiToken
}

func init() {
	registry.RegisterDriver(&githubCodeHostingDriver{})
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

func (d *githubCodeHostingDriver) getPullRequestNumber(options MergePullRequestOptions) (int, error) {
	pullRequests, err := d.getPullRequests(options.Branch, options.ParentBranch)
	if err != nil {
		return -1, err
	}
	if len(pullRequests) == 0 {
		return -1, errors.New("No pull request found")
	}
	if len(pullRequests) > 1 {
		pullRequestNumbersAsStrings := make([]string, len(pullRequests))
		for i, pullRequest := range pullRequests {
			pullRequestNumbersAsStrings[i] = strconv.Itoa(*pullRequest.Number)
		}
		return -1, fmt.Errorf("Multiple pull requests found: %s", strings.Join(pullRequestNumbersAsStrings, ", "))
	}
	return *pullRequests[0].Number, nil
}

func (d *githubCodeHostingDriver) getPullRequests(branch, parentBranch string) ([]*github.PullRequest, error) {
	pullRequests, _, err := d.client.PullRequests.List(context.Background(), d.owner, d.repository, &github.PullRequestListOptions{
		Base:  parentBranch,
		Head:  d.owner + ":" + branch,
		State: "open",
	})
	return pullRequests, err
}

func (d *githubCodeHostingDriver) mergePullRequest(options MergePullRequestOptions) (string, error) {
	pullRequestNumber, err := d.getPullRequestNumber(options)
	if err != nil {
		return "", err
	}
	if options.LogRequests {
		printLog(fmt.Sprintf("GitHub API: Merging PR #%d", pullRequestNumber))
	}
	commitMessageParts := strings.SplitN(options.CommitMessage, "\n", 2)
	githubCommitTitle := commitMessageParts[0]
	githubCommitMessage := ""
	if len(commitMessageParts) == 2 {
		githubCommitMessage = commitMessageParts[1]
	}
	result, _, err := d.client.PullRequests.Merge(context.Background(), d.owner, d.repository, pullRequestNumber, githubCommitMessage, &github.PullRequestOptions{
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
			printLog(fmt.Sprintf("GitHub API: Updating base branch for PR #%d", *pullRequest.Number))
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
