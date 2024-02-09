package connector

import (
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/google/go-github/v58/github"
)

func CommentsAuthors(comments []*github.IssueComment) data.Users {
	result := data.NewUsers()
	for _, comment := range comments {
		result.AddUser(CommentAuthor(comment))
	}
	return result
}

func CommentAuthor(comment *github.IssueComment) string {
	return *comment.User.Login
}
