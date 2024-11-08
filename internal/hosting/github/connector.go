package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/git-town/git-town/v16/internal/cli/colors"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	if len(hostingdomain.ReadProposalOverride()) == 0 {
		return Some(self.findProposalViaOverride)
	}
	if self.APIToken.IsSome() {
		return Some(self.findProposalViaAPI)
	}
	return None[func(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)]()
}

func (self Connector) NewProposalURL(branch, parentBranch, mainBranch gitdomain.LocalBranchName, proposalTitle gitdomain.ProposalTitle, proposalBody gitdomain.ProposalBody) (string, error) {
	toCompare := branch.String()
	if parentBranch != mainBranch {
		toCompare = parentBranch.String() + "..." + branch.String()
	}
	result := fmt.Sprintf("%s/compare/%s?expand=1", self.RepositoryURL(), url.PathEscape(toCompare))
	if len(proposalTitle) > 0 {
		result += "&title=" + url.QueryEscape(proposalTitle.String())
	}
	if len(proposalBody) > 0 {
		result += "&body=" + url.QueryEscape(proposalBody.String())
	}
	return result, nil
}

func (self Connector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", self.HostnameWithStandardPort(), self.Organization, self.Repository)
}

func (self Connector) SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)] {
	if self.APIToken.IsNone() {
		return None[func(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error)]()
	}
	return Some(self.searchProposal)
}

func (self Connector) SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) (err error)] {
	if self.APIToken.IsNone() {
		return None[func(number int, message gitdomain.CommitMessage) (err error)]()
	}
	return Some(self.squashMergeProposal)
}

func (self Connector) UpdateProposalSourceFn() Option[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error] {
	return None[func(number int, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	if self.APIToken.IsNone() {
		return None[func(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error]()
	}
	return Some(self.updateProposalTarget)
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	pullRequests, _, err := self.client.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		Base:  target.String(),
		State: "open",
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	if len(pullRequests) == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	if len(pullRequests) > 1 {
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(pullRequests), branch, target)
	}
	proposal := parsePullRequest(pullRequests[0])
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+strconv.Itoa(proposal.Number)), proposal.Title))
	return Some(proposal), nil
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	proposalURLOverride := hostingdomain.ReadProposalOverride()
	self.log.Ok()
	if proposalURLOverride == hostingdomain.OverrideNoProposal {
		return None[hostingdomain.Proposal](), nil
	}
	return Some(hostingdomain.Proposal{
		MergeWithAPI: true,
		Number:       123,
		Source:       branch,
		Target:       target,
		Title:        "title",
		URL:          proposalURLOverride,
	}), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	pullRequests, _, err := self.client.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		State: "open",
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[hostingdomain.Proposal](), err
	}
	if len(pullRequests) == 0 {
		self.log.Success("none")
		return None[hostingdomain.Proposal](), nil
	}
	if len(pullRequests) > 1 {
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
	proposal := parsePullRequest(pullRequests[0])
	self.log.Success(proposal.Target.String())
	return Some(proposal), nil
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) (err error) {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingGithubMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
	commitMessageParts := message.Parts()
	_, _, err = self.client.PullRequests.Merge(context.Background(), self.Organization, self.Repository, number, commitMessageParts.Text, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: commitMessageParts.Subject,
	})
	if err != nil {
		self.log.Ok()
	}
	return err
}

func (self Connector) updateProposalTarget(number int, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.PullRequests.Edit(context.Background(), self.Organization, self.Repository, number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &(targetName),
		},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

// getGitHubApiToken returns the GitHub API token to use.
// It first checks the GITHUB_TOKEN environment variable.
// If that is not set, it checks the GITHUB_AUTH_TOKEN environment variable.
// If that is not set, it checks the git config.
func GetAPIToken(gitConfigToken Option[configdomain.GitHubToken]) Option[configdomain.GitHubToken] {
	apiToken := os.ExpandEnv("$GITHUB_TOKEN")
	if len(apiToken) > 0 {
		return Some(configdomain.GitHubToken(apiToken))
	}
	apiToken = os.ExpandEnv("$GITHUB_AUTH_TOKEN")
	if len(apiToken) > 0 {
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
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.GetRef()),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.GetRef()),
		Title:        pullRequest.GetTitle(),
		MergeWithAPI: pullRequest.GetMergeableState() == "clean",
		URL:          *pullRequest.HTMLURL,
	}
}
