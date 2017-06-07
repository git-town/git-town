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
	client *github.Client
}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Github
func (driver GithubCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("https://github.com/%s/compare/%s?expand=1", repository, toCompare)
}

// GetRepositoryURL returns the URL of the given repository on github.com
func (driver GithubCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://github.com/" + repository
}

// MergePullRequest merges the pull request through the Github API
func (driver GithubCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) error {
	driver.connect()
	err := driver.updatePullRequestsAgainst(options)
	if err != nil {
		return err
	}
	return driver.mergePullRequest(options)
}

// Helper

func (driver *GithubCodeHostingDriver) connect() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	driver.client = github.NewClient(tc)
}

func (driver *GithubCodeHostingDriver) getPullRequestNumber(options MergePullRequestOptions) (int, error) {
	pullRequests, _, err := driver.client.PullRequests.List(context.Background(), options.Owner, options.Repository, &github.PullRequestListOptions{
		Base:  options.ParentBranch,
		Head:  options.Owner + ":" + options.Branch,
		State: "open",
	})
	if err != nil {
		return -1, err
	}
	if len(pullRequests) == 0 {
		return -1, errors.New("No pull request found")
	}
	if len(pullRequests) > 1 {
		pullRequestNumbers := make([]string, len(pullRequests))
		for i, pullRequest := range pullRequests {
			pullRequestNumbers[i] = strconv.Itoa(*pullRequest.Number)
		}
		return -1, fmt.Errorf("Multiple pull requests found: %s", strings.Join(pullRequestNumbers, ", "))
	}
	return *pullRequests[0].Number, nil
}

func (driver *GithubCodeHostingDriver) mergePullRequest(options MergePullRequestOptions) error {
	pullRequestNumber, err := driver.getPullRequestNumber(options)
	if err != nil {
		return err
	}
	_, _, err = driver.client.PullRequests.Merge(context.Background(), options.Owner, options.Repository, pullRequestNumber, options.CommitMessage, &github.PullRequestOptions{
		MergeMethod: "squash",
	})
	return err
}

func (driver *GithubCodeHostingDriver) updatePullRequestsAgainst(options MergePullRequestOptions) error {
	pullRequests, _, err := driver.client.PullRequests.List(context.Background(), options.Owner, options.Repository, &github.PullRequestListOptions{
		Base:  options.Branch,
		State: "open",
	})
	if err != nil {
		return err
	}
	for _, pullRequest := range pullRequests {
		_, _, err = driver.client.PullRequests.Edit(context.Background(), options.Owner, options.Repository, *pullRequest.Number, &github.PullRequest{
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
