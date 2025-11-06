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
	fmt.Println("FINDING PROPOSAL")
	self.log.Start(messages.APIProposalLookupStart)
	query := fmt.Sprintf(`source.branch.name = %q AND destination.branch.name = %q AND (state = "open" OR state = "new")`, branch, target)
	result1, err := self.client.Value.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    query,
		States:   []string{"open"},
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
	for _, proposalJSON := range proposals2 {
		proposal3, ok := proposalJSON.(map[string]any)
		if !ok {
			self.log.Failed(messages.APIUnexpectedResultDataStructure)
			return None[forgedomain.Proposal](), nil
		}
		proposal4, err := parsePullRequest(proposal3)
		if err != nil {
			self.log.Failed(err.Error())
			return None[forgedomain.Proposal](), nil
		}
		self.log.Success(fmt.Sprintf("#%d", proposal4.Number))
		fmt.Println("PROPOSAL", proposal4)
		return Some(forgedomain.Proposal{Data: proposal4, ForgeType: forgedomain.ForgeTypeBitbucket}), nil
	}
	self.log.Success("none")
	return None[forgedomain.Proposal](), nil
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector // type check

func (self APIConnector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	response1, err := self.client.Value.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    fmt.Sprintf(`source.branch.name = %q AND state = "open"`, branch),
		States:   []string{"open"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	response2, ok := response1.(map[string]any)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	size1, has := response2["size"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	size2, ok := size1.(float64)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	size3 := int(size2)
	if size3 == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	if size3 > 1 {
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromFound, size3, branch))
		return None[forgedomain.Proposal](), nil
	}
	values1, has := response2["values"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	values2, ok := values1.([]any)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	proposal1 := values2[0].(map[string]any)
	proposal2, err := parsePullRequest(proposal1)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), nil
	}
	self.log.Success(proposal2.Target.String())
	return Some(forgedomain.Proposal{Data: proposal2, ForgeType: forgedomain.ForgeTypeBitbucket}), nil
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
