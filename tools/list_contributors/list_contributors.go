package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

const org = "git-town"
const repo = "git-town"

func main() {
	githubToken := loadAccessToken()
	client, context := githubClient(githubToken)
	users := NewUserCollector()

	// get the tag
	if len(os.Args) < 2 {
		fmt.Println("Usage: list_contributors <previous tag>")
		os.Exit(1)
	}
	tag := os.Args[1]
	fmt.Printf("Looking for contributors since %s\n", tag)

	// determine time of the given tag
	tagTime := timeOfTag(tag)
	fmt.Printf("release %s was made %s\n", tag, tagTime.Format("2006-01-02"))

	// load and categorize all closed issues
	issues := []*github.Issue{}
	pullRequests := []*github.Issue{}
	page := 0
	query := fmt.Sprintf("repo:git-town/git-town closed:>=%s", tagTime.Format("2006-01-02"))
	for {
		fmt.Printf("loading issues %d-%d ... ", page+1, page+100)
		results, _, err := client.Search.Issues(context, query, &github.SearchOptions{
			Sort:  "closed",
			Order: "asc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("ok")
		if len(results.Issues) == 0 {
			break
		}
		for _, issue := range results.Issues {
			if issue.IsPullRequest() {
				pullRequests = append(pullRequests, issue)
			} else {
				issues = append(issues, issue)
			}
		}
		page += 1
	}
	fmt.Printf("found %d issues and %d pull requests\n", len(issues), len(pullRequests))

	// register the creators of pull requests
	for _, pullRequest := range pullRequests {
		users.AddUser(*pullRequest.User.Login)
	}

	// register the users involved in the tickets
	for _, issue := range issues {
		fmt.Printf("%s submitted issue %d (%s)\n", *issue.User.Login, *issue.Number, *issue.Title)
		users.AddUser(*issue.User.Login)
		comments, _, err := client.Issues.ListComments(context, "git-town", "git-town", *issue.Number, nil)
		if err != nil {
			panic(err)
		}
		for _, comment := range comments {
			users.AddUser(*comment.User.Login)
		}
	}
	fmt.Println("\nUsers:")
	fmt.Println()
	for _, username := range users.Users() {
		fmt.Println("@" + username)
	}
}

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", "git-town.github-token")
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	result := strings.TrimSpace(string(output))
	fmt.Printf("using GitHub token %q\n", result)
	return result
}

func githubClient(token string) (*github.Client, context.Context) {
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context, tokenSource)
	return github.NewClient(httpClient), context
}

func timeOfTag(tag string) time.Time {
	cmd := exec.Command("git", "log", "-1", "--format=%cI", tag)
	outputData, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	output := strings.TrimSpace(string(outputData))
	result, err := time.Parse(time.RFC3339, output)
	if err != nil {
		panic(err.Error())
	}
	return result
}
