package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/git-town/git-town/v20/internal/cli/colors"
	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/git/giturl"
	"github.com/git-town/git-town/v20/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	forgedomain.Data
	APIToken Option[configdomain.GitHubToken]
	client   *github.Client
	log      print.Logger
}

func (self Connector) DefaultProposalMessage(proposal forgedomain.Proposal) string {
	return forgedomain.CommitBody(proposal, fmt.Sprintf("%s (#%d)", proposal.Title(), proposal.Number()))
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if len(forgedomain.ReadProposalOverride()) > 0 {
		return Some(self.findProposalViaOverride)
	}
	if self.APIToken.IsSome() {
		return Some(self.findProposalViaAPI)
	}
	return None[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
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

func (self Connector) SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if self.APIToken.IsNone() {
		return None[func(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
	}
	return Some(self.searchProposal)
}

func (self Connector) SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) (err error)] {
	if self.APIToken.IsNone() {
		return None[func(number int, message gitdomain.CommitMessage) (err error)]()
	}
	return Some(self.squashMergeProposal)
}

func (self Connector) UpdateProposalSourceFn() Option[func(proposal forgedomain.Proposal, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error] {
	return None[func(proposal forgedomain.Proposal, _ gitdomain.LocalBranchName, finalMessages stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(proposal forgedomain.Proposal, target gitdomain.LocalBranchName, _ stringslice.Collector) error] {
	if self.APIToken.IsNone() {
		return None[func(proposal forgedomain.Proposal, target gitdomain.LocalBranchName, _ stringslice.Collector) error]()
	}
	return Some(self.updateProposalTarget)
}

func (self Connector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	pullRequests, _, err := self.client.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		Base:  target.String(),
		State: "open",
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	if len(pullRequests) == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	if len(pullRequests) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(pullRequests), branch, target)
	}
	proposal := parsePullRequest(pullRequests[0])
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+strconv.Itoa(proposal.number)), proposal.title))
	return Some(forgedomain.Proposal(proposal)), nil
}

func (self Connector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := forgedomain.ReadProposalOverride()
	self.log.Ok()
	if proposalURLOverride == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	return Some(forgedomain.Proposal(
		Proposal{
			body:         None[string](),
			mergeWithAPI: true,
			number:       123,
			source:       branch,
			target:       target,
			title:        "title",
			url:          proposalURLOverride,
		})), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	pullRequests, _, err := self.client.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		State: "open",
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	if len(pullRequests) == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	if len(pullRequests) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
	proposal := parsePullRequest(pullRequests[0])
	self.log.Success(proposal.target.String())
	return Some(forgedomain.Proposal(proposal)), nil
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) (err error) {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.ForgeGithubMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
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

func (self Connector) updateProposalTarget(forgeProposal forgedomain.Proposal, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	proposal := forgeProposal.(Proposal)
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(proposal.number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.PullRequests.Edit(context.Background(), self.Organization, self.Repository, proposal.number, &github.PullRequest{
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
		Data: forgedomain.Data{
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
func parsePullRequest(pullRequest *github.PullRequest) Proposal {
	return Proposal{
		number:       pullRequest.GetNumber(),
		source:       gitdomain.NewLocalBranchName(pullRequest.Head.GetRef()),
		target:       gitdomain.NewLocalBranchName(pullRequest.Base.GetRef()),
		title:        pullRequest.GetTitle(),
		body:         NewOption(pullRequest.GetBody()),
		mergeWithAPI: pullRequest.GetMergeableState() == "clean",
		url:          *pullRequest.HTMLURL,
	}
}
