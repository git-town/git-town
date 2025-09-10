package github

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/pkg/colors"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/google/go-github/v58/github"
)

// type-check to ensure conformance to the Connector interface
var (
	githubAuthConnector APIConnector
	_                   forgedomain.Connector = githubAuthConnector
)

// APIConnector provides access to the GitHub API.
type APIConnector struct {
	WebConnector
	APIToken Option[forgedomain.GitHubToken]
	client   Mutable[github.Client]
	log      print.Logger
}

// FIND PROPOSALS

var _ forgedomain.ProposalFinder = githubAuthConnector

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	pullRequests, _, err := self.client.Value.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
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
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+strconv.Itoa(proposal.Number)), proposal.Title))
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitHub}), nil
}

// MERGE PROPOSALS

var _ forgedomain.ProposalMerger = githubAuthConnector

func (self APIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.ForgeGitHubMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
	commitMessageParts := message.Parts()
	_, _, err := self.client.Value.PullRequests.Merge(context.Background(), self.Organization, self.Repository, number, commitMessageParts.Text, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: commitMessageParts.Subject,
	})
	if err != nil {
		self.log.Ok()
	}
	return err
}

// SEARCH PROPOSALS

var _ forgedomain.ProposalSearcher = githubAuthConnector

func (self APIConnector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	pullRequests, _, err := self.client.Value.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
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
	self.log.Success(proposal.Target.String())
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitHub}), nil
}

// UPDATE PROPOSALS

var _ forgedomain.ProposalUpdater = githubAuthConnector

func (self APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody string) error {
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err := self.client.Value.PullRequests.Edit(context.Background(), self.Organization, self.Repository, data.Number, &github.PullRequest{
		Body: Ptr(updatedBody),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func (self APIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.Value.PullRequests.Edit(context.Background(), self.Organization, self.Repository, data.Number, &github.PullRequest{
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

// VERIFY LOGIN

var _ forgedomain.AuthVerifier = githubAuthConnector

func (self APIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, _, err := self.client.Value.Users.Get(context.Background(), "")
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.Value.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(*user.Login),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
