package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/git-town/git-town/v24/internal/cli/print"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/messages"
	"github.com/git-town/git-town/v24/pkg/colors"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/google/go-github/v58/github"
)

// type checks
var (
	apiConnector APIConnector
	_            forgedomain.Connector = apiConnector
)

// APIConnector provides access to the GitHub API.
type APIConnector struct {
	client Mutable[github.Client]
	log    print.Logger
	WebConnector
}

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalFindStart, branch, target)
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
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+proposalData.Number.String()), proposalData.Title))
	proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeGithub}
	return Some(proposal), nil
}

func (self APIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	self.log.Start(messages.APIProposalSearchStart, branch.String())
	pullRequests, _, err := self.client.Value.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		Head:  self.Organization + ":" + branch.String(),
		State: "open",
	})
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	result := make([]forgedomain.Proposal, len(pullRequests))
	ids := make([]string, len(pullRequests))
	for p, pullRequest := range pullRequests {
		proposalData := parsePullRequest(pullRequest)
		proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeGithub}
		result[p] = proposal
		ids[p] = colors.BoldGreen().Styled(fmt.Sprintf("#%d", proposalData.Number))
	}
	if len(pullRequests) == 0 {
		self.log.Success("none")
	} else {
		self.log.Log(strings.Join(ids, ", "))
	}
	return result, nil
}

func (self APIConnector) SquashMergeProposal(number forgedomain.ProposalNumber, message Option[gitdomain.CommitMessage]) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.ForgeGithubMergingViaAPI, colors.BoldGreen().Styled("#"+number.String()))
	// When no commit message is given, GitHub determines the squash commit message.
	commitMessageParts := message.GetOrZero().Parts()
	_, _, err := self.client.Value.PullRequests.Merge(context.Background(), self.Organization, self.Repository, number.Int(), commitMessageParts.Body, &github.PullRequestOptions{
		CommitTitle: commitMessageParts.Title.String(),
		MergeMethod: "squash",
	})
	self.log.Finished(err)
	return err
}

func (self APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody gitdomain.ProposalBody) error {
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+data.Number.String()))
	_, _, err := self.client.Value.PullRequests.Edit(context.Background(), self.Organization, self.Repository, data.Number.Int(), &github.PullRequest{
		Body: new(updatedBody.String()),
	})
	self.log.Finished(err)
	return err
}

func (self APIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+data.Number.String()), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.Value.PullRequests.Edit(context.Background(), self.Organization, self.Repository, data.Number.Int(), &github.PullRequest{
		Base: &github.PullRequestBranch{
			Ref: &targetName,
		},
	})
	self.log.Finished(err)
	return err
}

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

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = apiConnector // type check

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector // type check

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = apiConnector // type check

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = apiConnector // type check

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = apiConnector // type check

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = apiConnector // type check
