package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/tools/release_stats/connector"
	"github.com/git-town/git-town/tools/release_stats/console"
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/git-town/git-town/tools/release_stats/git"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: list_contributors <previous tag>")
		os.Exit(1)
	}
	lastRelease := git.LoadTag(os.Args[1])
	fmt.Printf("previous release %s was on %s\n", console.Green.Styled(lastRelease.Name), console.Cyan.Styled(lastRelease.ISOTime))
	contributors := data.NewUsers()
	gh := connector.NewConnector()

	closedIssues, closedPullRequests := gh.ClosedIssues(lastRelease.ISOTime)

	fmt.Println("\n\nDETERMINING PARTICIPANTS IN CLOSED ISSUES")
	fmt.Println()
	contributors.AddUsers(gh.IssuesParticipants(closedIssues, "issue"))

	fmt.Println("\n\nDETERMINING PARTICIPANTS IN CLOSED PULL REQUESTS")
	fmt.Println()
	contributors.AddUsers(gh.IssuesParticipants(closedPullRequests, "PR"))

	// print statistics
	fmt.Printf("\n%d shipped PRs\n", len(closedPullRequests))
	fmt.Printf("%d resolved issues\n", len(closedIssues))
	users := contributors.Users()
	fmt.Printf("%d contributors:\n", len(users))
	userNames := []string{}
	for _, username := range users {
		userNames = append(userNames, "@"+username)
	}
	fmt.Println(strings.Join(userNames, ", "))
}
