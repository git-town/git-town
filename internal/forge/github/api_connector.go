package github

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/google/go-github/v58/github"
)

// type checks
var (
	apiConnector APIConnector
	_            forgedomain.Connector = apiConnector
)

// APIConnector provides access to the GitHub API.
type APIConnector struct {
	WebConnector
	cache  forgedomain.ProposalCache
	client Mutable[github.Client]
	log    print.Logger
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = apiConnector // type check

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if proposal := self.cache.BySourceTarget(branch, target); proposal.IsSome() {
		return proposal, nil
	}
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
	proposalData := parsePullRequest(pullRequests[0])
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+strconv.Itoa(proposalData.Number)), proposalData.Title))
	proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeGitHub}
	self.cache.Set(proposal)
	return Some(proposal), nil
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector // type check

func (self APIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if proposals := self.cache.BySource(branch); len(proposals) > 0 {
		return proposals, nil
	}
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	pullRequests, _, err := self.client.Value.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		State: "open",
	})
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	result := make([]forgedomain.Proposal, len(pullRequests))
	for p, pullRequest := range pullRequests {
		proposalData := parsePullRequest(pullRequest)
		self.log.Success(proposalData.Target.String())
		proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeGitHub}
		result[p] = proposal
	}
	if len(pullRequests) == 0 {
		self.log.Success("none")
	}
	self.cache.SetMany(result)
	return result, nil
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = apiConnector // type check

func (self APIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.ForgeGitHubMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
	commitMessageParts := message.Parts()
	_, _, err := self.client.Value.PullRequests.Merge(context.Background(), self.Organization, self.Repository, number, commitMessageParts.Body, &github.PullRequestOptions{
		MergeMethod: "squash",
		CommitTitle: commitMessageParts.Title.String(),
	})
	if err != nil {
		self.log.Ok()
	}
	return err
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = apiConnector // type check

func (self APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody string) error {
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err := self.client.Value.PullRequests.Edit(context.Background(), self.Organization, self.Repository, data.Number, &github.PullRequest{
		Body: Ptr(updatedBody),
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = apiConnector // type check

func (self APIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.Value.PullRequests.Edit(context.Background(), self.Organization, self.Repository, data.Number, &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &(targetName),
		},
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = apiConnector // type check

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
