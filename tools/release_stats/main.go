package main

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
)

const (
	org      = "git-town" // name of the GitHub organization
	repo     = "git-town" // name of the GitHub repo
	pageSize = 100        // how many results to load per page from the API
)

var cyan = termenv.String().Foreground(termenv.ANSICyan)

func main() {
	contributors := NewUsers()
	gh := newGithubConnector()
	lastRelease := loadPreviousRelease()

	// Add people who opened issues or pull requests since the last release.
	openedIssuesOrPRs := gh.openedIssuesOrPRsSince(lastRelease.ISOTime)
	contributors.AddUsers(gh.issuesParticipants(openedIssuesOrPRs))

	// Add people who were involved with issues and pull requests that were closed in this release.
	// People are counted as contributors to this release even if their interaction was a long time ago,
	// as long as what they contributed to is a part of this release.
	closedIssues, closedPullRequests := gh.loadClosedIssues(lastRelease.ISOTime)
	contributors.AddUsers(gh.issuesParticipants(closedIssues))
	contributors.AddUsers(gh.issuesParticipants(closedPullRequests))

	// people who made any comment on any issue (old or new, open or closed) since the last release

	// people who added a reaction on anything issue since the last release

	// load all people who commented on pull requests since the last release

	// load all people who reacted since the last release
	// relevant issues = all open issues and issues closed since the last tag
	// for each issue:
	//   for each comment of the issue
	//     for each reaction to the comment
	//       if the reaction was made since the last tag: register the user

	// print statistics
	// shipped pull requests
	// closed issues
	// contributors

	// register the users involved in the tickets
	userNames := []string{}
	for _, username := range contributors.Users() {
		userNames = append(userNames, "@"+username)
	}
	fmt.Println(strings.Join(userNames, ", "))
}
