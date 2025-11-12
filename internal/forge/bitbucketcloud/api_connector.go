package bitbucketcloud

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
	"github.com/ktrysmt/go-bitbucket"
)

// type checks
var (
	apiConnector APIConnector
	_            forgedomain.Connector = apiConnector
)

// APIConnector provides access to the Bitbucket Cloud API.
type APIConnector struct {
	WebConnector
	client Mutable[bitbucket.Client]
	log    print.Logger
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = apiConnector // type check

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	query := fmt.Sprintf("source.branch.name = %q AND destination.branch.name = %q", branch, target)
	result1, err := self.client.Value.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    query,
		States:   []string{"open", "new"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	if result1 == nil {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	result2, ok := result1.(map[string]any)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	proposals1, has := result2["values"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	proposals2, ok := proposals1.([]any)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	for _, proposal1 := range proposals2 {
		proposal2, ok := proposal1.(map[string]any)
		if !ok {
			self.log.Failed(messages.APIUnexpectedResultDataStructure)
			return None[forgedomain.Proposal](), nil
		}
		proposal3, err := parsePullRequest(proposal2)
		if err != nil {
			self.log.Failed(err.Error())
			return None[forgedomain.Proposal](), nil
		}
		if !proposal3.Active {
			continue
		}
		self.log.Success(fmt.Sprintf("#%d", proposal3.Number))
		return Some(forgedomain.Proposal{Data: proposal3, ForgeType: forgedomain.ForgeTypeBitbucket}), nil
	}
	self.log.Success("none")
	return None[forgedomain.Proposal](), nil
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector // type check

func (self APIConnector) SearchProposal(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	response1, err := self.client.Value.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    fmt.Sprintf("source.branch.name = %q", branch),
		States:   []string{"open", "new"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	response2, ok := response1.(map[string]any)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return []forgedomain.Proposal{}, nil
	}
	proposals1, has := response2["values"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return []forgedomain.Proposal{}, nil
	}
	proposals2, ok := proposals1.([]any)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return []forgedomain.Proposal{}, nil
	}
	result := make([]forgedomain.Proposal, len(proposals2))
	for p, proposal1 := range proposals2 {
		proposal2, ok := proposal1.(map[string]any)
		if !ok {
			self.log.Failed(messages.APIUnexpectedResultDataStructure)
			return []forgedomain.Proposal{}, nil
		}
		proposal3, err := parsePullRequest(proposal2)
		if err != nil {
			self.log.Failed(err.Error())
			return []forgedomain.Proposal{}, nil
		}
		if !proposal3.Active {
			continue
		}
		self.log.Success(fmt.Sprintf("#%d ", proposal3.Number))
		result[p] = forgedomain.Proposal{Data: proposal3, ForgeType: forgedomain.ForgeTypeBitbucket}
	}
	if len(result) == 0 {
		self.log.Success("none")
	}
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
	self.log.Start(messages.ForgeBitbucketMergingViaAPI, colors.BoldGreen().Styled("#"+strconv.Itoa(number)))
	_, err := self.client.Value.Repositories.PullRequests.Merge(&bitbucket.PullRequestsOptions{
		ID:       strconv.Itoa(number),
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Message:  message.String(),
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = apiConnector // type check

func (self APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody string) error {
	data := proposalData.(forgedomain.BitbucketCloudProposalData)
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, err := self.client.Value.Repositories.PullRequests.Update(&bitbucket.PullRequestsOptions{
		ID:                strconv.Itoa(data.Number),
		Owner:             self.Organization,
		RepoSlug:          self.Repository,
		SourceBranch:      data.Source.String(),
		DestinationBranch: data.Target.String(),
		Title:             data.Title,
		Description:       newBody,
		Draft:             data.Draft,
		CloseSourceBranch: data.CloseSourceBranch,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// udpate proposal source
// ============================================================================

var _ forgedomain.ProposalSourceUpdater = apiConnector // type check

func (self APIConnector) UpdateProposalSource(proposalData forgedomain.ProposalInterface, source gitdomain.LocalBranchName) error {
	data := proposalData.(forgedomain.BitbucketCloudProposalData)
	self.log.Start(messages.APIUpdateProposalSource, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(source.String()))
	_, err := self.client.Value.Repositories.PullRequests.Update(&bitbucket.PullRequestsOptions{
		ID:                strconv.Itoa(data.Number),
		Owner:             self.Organization,
		RepoSlug:          self.Repository,
		SourceBranch:      source.String(),
		DestinationBranch: data.Target.String(),
		Title:             data.Title,
		Description:       data.Body.GetOrZero(),
		Draft:             data.Draft,
		CloseSourceBranch: data.CloseSourceBranch,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = apiConnector // type check

func (self APIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.(forgedomain.BitbucketCloudProposalData)
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(target.String()))
	_, err := self.client.Value.Repositories.PullRequests.Update(&bitbucket.PullRequestsOptions{
		ID:                strconv.Itoa(data.Number),
		Owner:             self.Organization,
		RepoSlug:          self.Repository,
		SourceBranch:      data.Source.String(),
		DestinationBranch: target.String(),
		Title:             data.Title,
		Description:       data.Body.GetOrZero(),
		Draft:             data.Draft,
		CloseSourceBranch: data.CloseSourceBranch,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = apiConnector // type check

func (self APIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, err := self.client.Value.User.Profile()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, err = self.client.Value.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    "",
		States:   []string{},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.Username),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
