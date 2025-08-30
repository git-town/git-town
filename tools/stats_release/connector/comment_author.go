package connector

import "github.com/google/go-github/v58/github"

func commentAuthor(comment github.IssueComment) string {
	return *comment.User.Login
}
