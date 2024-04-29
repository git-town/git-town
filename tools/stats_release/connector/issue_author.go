package connector

import (
	"github.com/google/go-github/v58/github"
)

func issueAuthor(issue github.Issue) string {
	return *issue.User.Login
}
