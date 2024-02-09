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
	fmt.Printf("previous release %s was on %s\n", console.Cyan.Styled(lastRelease.Name), console.Cyan.Styled(lastRelease.ISOTime))
	contributors := data.NewUsers()
	gh := connector.NewConnector()

	// Add people who opened issues or pull requests since the last release.
	// openedIssuesOrPRs := gh.OpenedIssuesOrPRsSince(lastRelease.ISOTime)
	// contributors.AddUsers(gh.IssuesParticipants(openedIssuesOrPRs))

	// Add people who were involved with issues and pull requests that were resolved in this release.
	closedIssues, closedPullRequests := gh.ClosedIssues(lastRelease.ISOTime)
	contributors.AddUsers(gh.IssuesParticipants(closedIssues))
	contributors.AddUsers(gh.IssuesParticipants(closedPullRequests))

	// Add people who made any comment on any issue (old or new, open or closed) since the last release
	// newComments := gh.CommentsSince(lastRelease)
	// contributors.AddUsers(connector.CommentsAuthors(newComments))

	// Add people who added a reaction on any issue since the last release.
	// We only look at open issues + issues that were closed in the last release here.
	// openIssues := gh.OpenIssues()
	// allIssues := append(openIssues, closedIssues...)
	// allIssuesComments := gh.IssuesComments(allIssues)
	// contributors.AddUsers(connector.CommentsAuthors(allIssuesComments))

	// load all people who commented on pull requests since the last release

	// load all people who reacted since the last release
	// relevant issues = all open issues and issues closed since the last tag
	// for each issue:
	//   for each comment of the issue
	//     for each reaction to the comment
	//       if the reaction was made since the last tag: register the user

	// print statistics
	fmt.Printf("%d shipped PRs", len(closedPullRequests))
	fmt.Printf("%d resolved issues", len(closedIssues))
	users := contributors.Users()
	fmt.Printf("%d contributors:", len(users))
	userNames := []string{}
	for _, username := range users {
		userNames = append(userNames, "@"+username)
	}
	fmt.Println(strings.Join(userNames, ", "))
}
