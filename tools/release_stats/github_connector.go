package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
)

// wrapper around the low-level GitHub connector, provides higher-level functions needed by this program
type githubConnector struct {
	client  *github.Client
	context context.Context
}

func newGithubConnector() githubConnector {
	githubToken := loadAccessToken()
	fmt.Printf("using GitHub token %s\n", cyan.Styled(githubToken))
	client, context := createGithubClient(githubToken)
	return githubConnector{
		client:  client,
		context: context,
	}
}
