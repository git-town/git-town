package connector

import "github.com/google/go-github/v58/github"

func reactionAuthor(reaction github.Reaction) string {
	return *reaction.User.Login
}
