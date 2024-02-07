package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

const org = "git-town"
const repo = "git-town"

func main() {
	// load access token
	token := loadAccessToken()
	ctx := context.Background()
	client := github.NewClient(nil).WithAuth(token, "")

	// Get all repositories in the organization
	repos, _, err := client.Repositories.List(ctx, orgName, nil)
	if err != nil {
		panic(err)
	}

	// Set to store unique usernames
	usernames := map[string]struct{}{}

	// Loop through each repository
	for _, repo := range repos {
		// Define date 3 months ago
		threeMonthsAgo := time.Now().AddDate(0, -3, 0)

		// Search for issues created in the last 3 months
		issuesOpts := &github.SearchOptions{
			Q:    fmt.Sprintf("repo:%s is:issue created:>=%s", *repo.Name, threeMonthsAgo.Format("2006-01-02")),
			Sort: "created",
			Dir:  "desc",
		}
		issues, _, err := client.Search.Issues(ctx, issuesOpts)
		if err != nil {
			panic(err)
		}

		// Loop through each issue
		for _, issue := range issues.Items {
			// Get comments on the issue
			comments, _, err := client.Issues.ListComments(ctx, *repo.Owner.Login, *repo.Name, *issue.Number, nil)
			if err != nil {
				panic(err)
			}

			// Extract usernames from comments
			for _, comment := range comments {
				usernames[*comment.User.Login] = struct{}{}
			}
		}

		// Search for pull requests created in the last 3 months
		pullRequestsOpts := &github.SearchOptions{
			Q:    fmt.Sprintf("repo:%s is:pr created:>=%s", *repo.Name, threeMonthsAgo.Format("2006-01-02")),
			Sort: "created",
			Dir:  "desc",
		}
		pullRequests, _, err := client.Search.PullRequests(ctx, pullRequestsOpts)
		if err != nil {
			panic(err)
		}

		// Loop through each pull request
		for _, pullRequest := range pullRequests.Items {
			// Get comments on the pull request
			comments, _, err := client.PullRequests.ListComments(ctx, *repo.Owner.Login, *repo.Name, *pullRequest.Number, nil)
			if err != nil {
				panic(err)
			}

			// Extract usernames from comments
			for _, comment := range comments {
				usernames[*comment.User.Login] = struct{}{}
			}
		}
	}

	// Print unique usernames
	for username := range usernames {
		fmt.Println(username)
	}
}

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", "git-town.github-token")
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	return string(output)
}

func githubClient() {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "YourGitHubToken"},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return github.NewClient(httpClient)
}
