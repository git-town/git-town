package main

import (
	"context"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

func createGithubClient(token string) (*github.Client, context.Context) {
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context, tokenSource)
	return github.NewClient(httpClient), context
}
