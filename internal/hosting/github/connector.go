package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	hostingdomain.Data
	APIToken Option[configdomain.GitHubToken]
	client   *github.Client
	log      print.Logger
}

func (self Connector) DefaultProposalMessage(proposal hostingdomain.Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	pullRequests, _, err := self.client.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		Base:  target.String(),
		State: "open",
	})
	if err != nil {
		return None[hostingdomain.Proposal](), err
	}
	if len(pullRequests) == 0 {
		return None[hostingdomain.Proposal](), nil
	}
	if len(pullRequests) > 1 {
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFound, len(pullRequests), branch, target)
	}
	proposal := parsePullRequest(pullRequests[0])
	return Some(proposal), nil
}

func (self Connector) NewProposalURL(branch, parentBranch, mainBranch gitdomain.LocalBranchName, proposalTitle gitdomain.ProposalTitle, proposalBody gitdomain.ProposalBody) (string, error) {
	toCompare := branch.String()
	if parentBranch != mainBranch {
		toCompare = parentBranch.String() + "..." + branch.String()
	}
	result := fmt.Sprintf("%s/compare/%s?expand=1", self.RepositoryURL(), url.PathEscape(toCompare))
	if proposalTitle != "" {
		result += "&title=" + url.QueryEscape(proposalTitle.String())
	}
	if proposalBody != "" {
		result += "&body=" + url.QueryEscape(proposalBody.String())
	}
	return result, nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) (err error) {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingGithubMergingViaAPI, number)
	commitMessageParts := message.Parts()
	_, _, err = self.client.PullRequests.Merge(context.Background(), self.Organization, self.Repository, number, commitMessageParts.Text, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: commitMessageParts.Subject,
	})
	self.log.Success()
	return err
}

func (self Connector) UpdateProposalTarget(number int, target gitdomain.LocalBranchName) error {
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
func GetAPIToken(gitConfigToken Option[configdomain.GitHubToken]) Option[configdomain.GitHubToken] {
	apiToken := os.ExpandEnv("$GITHUB_TOKEN")
	if apiToken != "" {
		return Some(configdomain.GitHubToken(apiToken))
	}
	apiToken = os.ExpandEnv("$GITHUB_AUTH_TOKEN")
	if apiToken != "" {
		return Some(configdomain.GitHubToken(apiToken))
	}
	return gitConfigToken
}

// NewConnector provides a fully configured GithubConnector instance
// if the current repo is hosted on GitHub, otherwise nil.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	githubClient := github.NewClient(httpClient)
	if args.RemoteURL.Host != "github.com" {
		url := "https://" + args.RemoteURL.Host
		var err error
		githubClient, err = githubClient.WithEnterpriseURLs(url, url)
		if err != nil {
			return Connector{}, fmt.Errorf(messages.GitHubEnterpriseInitializeError, err)
		}
	}
	return Connector{
		APIToken: args.APIToken,
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: githubClient,
		log:    args.Log,
	}, nil
}

type NewConnectorArgs struct {
	APIToken  Option[configdomain.GitHubToken]
	Log       print.Logger
	RemoteURL giturl.Parts
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
