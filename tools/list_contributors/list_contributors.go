package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

const org = "git-town"
const repo = "git-town"

func main() {
	client, context := githubClient()
	users := UserCollector{}

	// get the tag
	if len(os.Args) < 2 {
		fmt.Println("Usage: list_contributors <previous tag>")
		os.Exit(1)
	}
	tag := os.Args[1]

	// determine time of the given tag
	tagTime := timeOfTag(tag)
	fmt.Printf("release %s was made %s\n", tag, tagTime)

	// add users that created or commented on issues since the last tag
	query := fmt.Sprintf("repo:git-town/git-town closed:>=%s", tagTime.Format("2006-01-02"))
	issues, _, err := client.Search.Issues(context, query, &github.SearchOptions{
		Sort:  "closed",
		Order: "asc",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d issues and pull requests were closed since %s\n", len(issues.Issues), tagTime.Format("2006-01-02"))
	for _, issue := range issues.Issues {
		users.AddUser(*issue.User.Login)
		comments, _, err := client.Issues.ListComments(context, "git-town", "git-town", *issue.Number, nil)
		if err != nil {
			panic(err)
		}
		for _, comment := range comments {
			users.AddUser(*comment.User.Login)
		}
	}

	// Print unique usernames
	for username := range users.Users() {
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

func githubClient() (*github.Client, context.Context) {
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "YourGitHubToken"})
	httpClient := oauth2.NewClient(context, tokenSource)
	return github.NewClient(httpClient), context
}

func timeOfTag(tag string) time.Time {
	cmd := exec.Command("git", "show", "--format=%ci", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	result, err := time.Parse(time.RFC3339, string(output))
	if err != nil {
		panic(err.Error())
	}
	return result
}
