package connector

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

const (
	org      = "git-town" // name of the GitHub organization
	repo     = "git-town" // name of the GitHub repo
	pageSize = 100        // how many results to load per page from the API
)

// wrapper around the low-level GitHub connector, provides higher-level functions needed by this program
type Connector struct {
	client  *github.Client
	context context.Context //nolint:containedctx // we are sure there is always only one context here, this is just a little script
}

func NewConnector() Connector {
	githubToken := loadAccessToken()
	fmt.Printf("using GitHub token %s\n", githubToken)
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	httpClient := oauth2.NewClient(context, tokenSource)
	return Connector{
		client:  github.NewClient(httpClient),
		context: context,
	}
}
