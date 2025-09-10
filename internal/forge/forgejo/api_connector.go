package forgejo

import (
	"errors"
	"fmt"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/pkg/colors"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to ensure conformance to the Connector interface
var (
	apiConnector APIConnector
	_            forgedomain.CredentialVerifier = apiConnector
	_            forgedomain.Connector          = apiConnector
)

// APIConnector provides access to the Forgejo API.
type APIConnector struct {
	WebConnector
	APIToken Option[forgedomain.ForgejoToken]
	client   *forgejo.Client
	log      print.Logger
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = apiConnector

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 50,
		},
		State: forgejo.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	pullRequests := FilterPullRequests(openPullRequests, branch, target)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parsePullRequest(pullRequests[0])
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeForgejo}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(pullRequests), branch, target)
	}
}

// ============================================================================
// merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = apiConnector

func (self APIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.ForgeForgejoMergingViaAPI, colors.BoldGreen().Styled(strconv.Itoa(number)))
	_, _, err := self.client.MergePullRequest(self.Organization, self.Repository, int64(number), forgejo.MergePullRequestOption{
		Style:   forgejo.MergeStyleSquash,
		Title:   commitMessageParts.Subject,
		Message: commitMessageParts.Text,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	self.log.Start(messages.APIProposalLookupStart)
	_, _, err = self.client.GetPullRequest(self.Organization, self.Repository, int64(number))
	self.log.Ok()
	return err
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector

func (self APIConnector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	openPullRequests, _, err := self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 50,
		},
		State: forgejo.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	pullRequests := filterPullRequests2(openPullRequests, branch)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		pullRequest := pullRequests[0]
		proposal := parsePullRequest(pullRequest)
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeForgejo}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = apiConnector

func (self APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody string) error {
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err := self.client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), forgejo.EditPullRequestOption{
		Body: newBody,
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
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(targetName))
	_, _, err := self.client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), forgejo.EditPullRequestOption{
		Base: targetName,
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
	user, _, err := self.client.GetMyUserInfo()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 1,
		},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.UserName),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}

func FilterPullRequests(pullRequests []*forgejo.PullRequest, branch, target gitdomain.LocalBranchName) []*forgejo.PullRequest {
	result := []*forgejo.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() && pullRequest.Base.Name == target.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

func parsePullRequest(pullRequest *forgejo.PullRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		MergeWithAPI: pullRequest.Mergeable,
		Number:       int(pullRequest.Index),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.Ref),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        pullRequest.Title,
		Body:         NewOption(pullRequest.Body),
		URL:          pullRequest.HTMLURL,
	}
}

func filterPullRequests2(pullRequests []*forgejo.PullRequest, branch gitdomain.LocalBranchName) []*forgejo.PullRequest {
	result := []*forgejo.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}
