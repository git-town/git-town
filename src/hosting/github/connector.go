package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/git-town/git-town/v13/src/cli/print"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/git/giturl"
	"github.com/git-town/git-town/v13/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	hostingdomain.Config
	APIToken   configdomain.GitHubToken
	MainBranch gitdomain.LocalBranchName
	client     *github.Client
	log        print.Logger
}

func (self *Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self *Connector) FindProposal(branch, target gitdomain.LocalBranchName) (*hostingdomain.Proposal, error) {
	pullRequests, _, err := self.client.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		Base:  target.String(),
		State: "open",
	})
	if err != nil {
		return nil, err
	}
	if len(pullRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(pullRequests) > 1 {
		return nil, fmt.Errorf(messages.ProposalMultipleFound, len(pullRequests), branch, target)
	}
	proposal := parsePullRequest(pullRequests[0])
	return &proposal, nil
}

func (self *Connector) NewProposalURL(branch, parentBranch gitdomain.LocalBranchName) (string, error) {
	toCompare := branch.String()
	if parentBranch != self.MainBranch {
		toCompare = parentBranch.String() + "..." + branch.String()
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", self.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (self *Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self *Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) (err error) {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingGithubMergingViaAPI, number)
	commitMessageParts := message.Split()
	_, _, err = self.client.PullRequests.Merge(context.Background(), self.Organization, self.Repository, number, commitMessageParts.Body, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: commitMessageParts.Title,
	})
	self.log.Success()
	return err
}

func (self *Connector) UpdateProposalTarget(number int, target gitdomain.LocalBranchName) error {
	self.log.Start(messages.HostingGithubUpdatePRViaAPI, number)
	targetName := target.String()
	_, _, err := self.client.PullRequests.Edit(context.Background(), self.Organization, self.Repository, number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &(targetName),
		},
	})
	if err != nil {
		self.log.Failed(err)
		return err
	}
	self.log.Success()
	return nil
}

// getGitHubApiToken returns the GitHub API token to use.
// It first checks the GITHUB_TOKEN environment variable.
// If that is not set, it checks the GITHUB_AUTH_TOKEN environment variable.
// If that is not set, it checks the git config.
func GetAPIToken(gitConfigToken configdomain.GitHubToken) configdomain.GitHubToken {
	apiToken := os.ExpandEnv("$GITHUB_TOKEN")
	if apiToken != "" {
		return configdomain.GitHubToken(apiToken)
	}
	apiToken = os.ExpandEnv("$GITHUB_AUTH_TOKEN")
	if apiToken != "" {
		return configdomain.GitHubToken(apiToken)
	}
	return gitConfigToken
}

// NewConnector provides a fully configured GithubConnector instance
// if the current repo is hosted on Github, otherwise nil.
func NewConnector(args NewConnectorArgs) (*Connector, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return &Connector{
		APIToken: args.APIToken,
		Config: hostingdomain.Config{
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		MainBranch: args.MainBranch,
		client:     github.NewClient(httpClient),
		log:        args.Log,
	}, nil
}

type NewConnectorArgs struct {
	APIToken        configdomain.GitHubToken
	HostingPlatform configdomain.HostingPlatform
	Log             print.Logger
	MainBranch      gitdomain.LocalBranchName
	OriginURL       *giturl.Parts
}

// parsePullRequest extracts standardized proposal data from the given GitHub pull-request.
func parsePullRequest(pullRequest *github.PullRequest) hostingdomain.Proposal {
	return hostingdomain.Proposal{
		Number:       pullRequest.GetNumber(),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.GetRef()),
		Title:        pullRequest.GetTitle(),
		MergeWithAPI: pullRequest.GetMergeableState() == "clean",
	}
}
