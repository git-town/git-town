package drivers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"github.com/Originate/git-town/src/git"
	"github.com/google/go-github/github"
)

type githubCodeHostingDriver struct {
	originURL  string
	hostname   string
	client     *github.Client
	owner      string
	repository string
}

func (d *githubCodeHostingDriver) CanBeUsed() bool {
	return d.hostname == "github.com" || strings.Contains(d.hostname, "github")
}

func (d *githubCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, error) {
	if os.Getenv("GIT_TOWN_GITHUB_TOKEN") == "" {
		return false, nil
	}
	d.connect()
	pullRequestNumbers, err := d.getPullRequestNumbers(branch, parentBranch)
	if err != nil {
		return false, err
	}
	return len(pullRequestNumbers) == 1, nil
}

func (d *githubCodeHostingDriver) GetNewPullRequestURL(branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", d.GetRepositoryURL(), toCompare)
}

func (d *githubCodeHostingDriver) GetRepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", d.owner, d.repository)
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
	return "Github"
}

func (d *githubCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.GetURLHostname(originURL)
	d.client = nil
	if d.CanBeUsed() {
		repositoryParts := strings.SplitN(git.GetURLRepositoryName(originURL), "/", 2)
		d.owner = repositoryParts[0]
		d.repository = repositoryParts[1]
	} else {
		d.owner = ""
		d.repository = ""
	}
}

func init() {
	registry.RegisterDriver(&githubCodeHostingDriver{})
}

// Helper

func (d *githubCodeHostingDriver) connect() {
	if d.client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GIT_TOWN_GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		d.client = github.NewClient(tc)
	}
}

func (d *githubCodeHostingDriver) getPullRequestNumber(options MergePullRequestOptions) (int, error) {
	pullRequestNumbers, err := d.getPullRequestNumbers(options.Branch, options.ParentBranch)
	if err != nil {
		return -1, err
	}
	if len(pullRequestNumbers) == 0 {
		return -1, errors.New("No pull request found")
	}
	if len(pullRequestNumbers) > 1 {
		pullRequestNumbersAsStrings := make([]string, len(pullRequestNumbers))
		for i, number := range pullRequestNumbers {
			pullRequestNumbersAsStrings[i] = strconv.Itoa(number)
		}
		return -1, fmt.Errorf("Multiple pull requests found: %s", strings.Join(pullRequestNumbersAsStrings, ", "))
	}
	return pullRequestNumbers[0], nil
}

func (d *githubCodeHostingDriver) getPullRequestNumbers(branch, parentBranch string) ([]int, error) {
	pullRequests, _, err := d.client.PullRequests.List(context.Background(), d.owner, d.repository, &github.PullRequestListOptions{
		Base:  parentBranch,
		Head:  d.owner + ":" + branch,
		State: "open",
	})
	if err != nil {
		return []int{}, err
	}
	result := make([]int, len(pullRequests))
	for i, pullRequest := range pullRequests {
		result[i] = *pullRequest.Number
	}
	return result, nil
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
