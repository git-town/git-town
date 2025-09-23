package gitlab

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// type checks
var (
	apiConnector APIConnector
	_            forgedomain.Connector = apiConnector
)

// APIConnector provides access to the GitLab API.
type APIConnector struct {
	WebConnector
	client *gitlab.Client
	log    print.Logger
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = apiConnector

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
		TargetBranch: gitlab.Ptr(target.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	switch len(mergeRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parseMergeRequest(mergeRequests[0])
		self.log.Success(strconv.Itoa(proposal.Number))
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitLab}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(mergeRequests), branch, target)
	}
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector

func (self APIConnector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	switch len(mergeRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parseMergeRequest(mergeRequests[0])
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitLab}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(mergeRequests), branch)
	}
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = apiConnector

func (self APIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.ForgeGitLabMergingViaAPI, number)
	// the GitLab API wants the full commit message in the body
	_, _, err := self.client.MergeRequests.AcceptMergeRequest(self.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.Ptr(message.String()),
		Squash:              gitlab.Ptr(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Ptr(false),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = apiConnector

func (self APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedDescription string) error {
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), data.Number, &gitlab.UpdateMergeRequestOptions{
		Description: Ptr(updatedDescription),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = apiConnector

func (self APIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	self.log.Start(messages.ForgeGitLabUpdateMRViaAPI, data.Number, target)
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), data.Number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.Ptr(target.String()),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = apiConnector

func (self APIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, _, err := self.client.Users.CurrentUser()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.Username),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
