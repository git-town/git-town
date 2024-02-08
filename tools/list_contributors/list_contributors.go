package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/muesli/termenv"
	"golang.org/x/oauth2"
)

const (
	org      = "git-town"
	repo     = "git-town"
	pagesize = 100
)

func main() {
	cyan := termenv.String().Foreground(termenv.ANSICyan)
	githubToken := loadAccessToken()
	fmt.Printf("using GitHub token %s\n", cyan.Styled(githubToken))
	client, context := githubClient(githubToken)
	users := NewUserCollector()

	// get the tag
	if len(os.Args) < 2 {
		fmt.Println("Usage: list_contributors <previous tag>")
		os.Exit(1)
	}
	tag := os.Args[1]

	// determine time of the given tag
	tagTime := timeOfTag(tag)
	fmt.Printf("previous release %s was made %s\n", cyan.Styled(tag), cyan.Styled(tagTime.Format("2006-01-02")))

	// load and categorize all closed issues
	issues := []*github.Issue{}
	pullRequests := []*github.Issue{}
	page := 0
	query := fmt.Sprintf("repo:git-town/git-town closed:>=%s", tagTime.Format("2006-01-02"))
	fmt.Print("loading issues ")
	for {
		results, _, err := client.Search.Issues(context, query, &github.SearchOptions{
			Sort:  "closed",
			Order: "asc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: pagesize,
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Print(".")
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
	fmt.Printf(" %s issues and %s pull requests\n", cyan.Styled(strconv.Itoa(len(issues))), cyan.Styled(strconv.Itoa(len(pullRequests))))

	// register the creators of pull requests
	for _, pullRequest := range pullRequests {
		users.AddUser(*pullRequest.User.Login)
	}

	// register the users involved in the tickets
	for _, issue := range issues {
		issueUsers := NewUserCollector()
		issueUsers.AddUser(*issue.User.Login)
		comments, _, err := client.Issues.ListComments(context, "git-town", "git-town", *issue.Number, nil)
		if err != nil {
			panic(err)
		}
		for _, comment := range comments {
			issueUsers.AddUser(*comment.User.Login)
		}
		users.AddUsers(issueUsers)
		fmt.Printf("#%d (%s): %s\n", *issue.Number, *issue.Title, cyan.Styled(strings.Join(issueUsers.Users(), ", ")))
	}
	fmt.Println("\nUsers:")
	fmt.Println()
	userNames := []string{}
	for _, username := range users.Users() {
		userNames = append(userNames, "@"+username)
	}
	fmt.Println(strings.Join(userNames, ", "))
}

func loadAccessToken() string {
	process := exec.Command("git", "config", "--get", "git-town.github-token")
	output, err := process.Output()
	if err != nil {
		panic(err.Error())
	}
	result := strings.TrimSpace(string(output))
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
