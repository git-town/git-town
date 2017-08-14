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

// GithubCodeHostingDriver provides tools for working with repositories
// hosted on Github
type GithubCodeHostingDriver struct {
	client     *github.Client
	owner      string
	repository string
}

// NewGithubCodeHostingDriver returns a new GithubCodeHostingDriver instance
func NewGithubCodeHostingDriver(fullRepository string) *GithubCodeHostingDriver {
	repositoryParts := strings.SplitN(fullRepository, "/", 2)
	return &GithubCodeHostingDriver{
		owner:      repositoryParts[0],
		repository: repositoryParts[1],
	}
}

// CanMergePullRequest returns whether or not MergePullRequest should be called when shipping
func (driver *GithubCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, error) {
	if os.Getenv("GIT_TOWN_GITHUB_TOKEN") == "" {
		return false, nil
	}
	driver.connect()
	pullRequestNumbers, err := driver.getPullRequestNumbers(branch, parentBranch)
	if err != nil {
		return false, err
	}
	return len(pullRequestNumbers) == 1, nil
}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Github
func (driver *GithubCodeHostingDriver) GetNewPullRequestURL(branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", driver.GetRepositoryURL(), toCompare)
}

// GetRepositoryURL returns the URL of the given repository on github.com
func (driver *GithubCodeHostingDriver) GetRepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", driver.owner, driver.repository)
}

// MergePullRequest merges the pull request through the Github API
func (driver *GithubCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	driver.connect()
	err := driver.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	return driver.mergePullRequest(options)
}

// Helper

func (driver *GithubCodeHostingDriver) connect() {
	if driver.client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GIT_TOWN_GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		driver.client = github.NewClient(tc)
	}
}

func (driver *GithubCodeHostingDriver) getPullRequestNumber(options MergePullRequestOptions) (int, error) {
	pullRequestNumbers, err := driver.getPullRequestNumbers(options.Branch, options.ParentBranch)
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

func (driver *GithubCodeHostingDriver) getPullRequestNumbers(branch, parentBranch string) ([]int, error) {
	pullRequests, _, err := driver.client.PullRequests.List(context.Background(), driver.owner, driver.repository, &github.PullRequestListOptions{
		Base:  parentBranch,
		Head:  driver.owner + ":" + branch,
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

func (driver *GithubCodeHostingDriver) mergePullRequest(options MergePullRequestOptions) (string, error) {
	pullRequestNumber, err := driver.getPullRequestNumber(options)
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
	result, _, err := driver.client.PullRequests.Merge(context.Background(), driver.owner, driver.repository, pullRequestNumber, githubCommitMessage, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: githubCommitTitle,
	})
	if err != nil {
		return "", err
	}
	return *result.SHA, nil
}

func (driver *GithubCodeHostingDriver) updatePullRequestsAgainst(options MergePullRequestOptions) error {
	pullRequests, _, err := driver.client.PullRequests.List(context.Background(), driver.owner, driver.repository, &github.PullRequestListOptions{
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
		_, _, err = driver.client.PullRequests.Edit(context.Background(), driver.owner, driver.repository, *pullRequest.Number, &github.PullRequest{
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
