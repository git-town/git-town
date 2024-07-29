package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/git-town/git-town/tools/stats_release/connector"
	"github.com/git-town/git-town/tools/stats_release/console"
	"github.com/git-town/git-town/tools/stats_release/data"
	"github.com/git-town/git-town/tools/stats_release/git"
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
	fmt.Println()
	fmt.Printf("%s shipped PRs\n", console.Green.Styled(strconv.Itoa(len(closedPullRequests))))
	fmt.Printf("%s resolved issues\n", console.Green.Styled(strconv.Itoa(len(closedIssues))))
	users := contributors.Values()
	fmt.Printf("%s contributors:\n", console.Cyan.Styled(strconv.Itoa(len(users))))
	userNames := []string{}
	for _, username := range users {
		userNames = append(userNames, "@"+username)
	}
	fmt.Println(strings.Join(userNames, ", "))
}
