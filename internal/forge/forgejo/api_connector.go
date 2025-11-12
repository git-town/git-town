package forgejo

import (
	"errors"
	"fmt"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	apiConnector APIConnector
	_            forgedomain.Connector = apiConnector
)

// APIConnector provides access to the Forgejo API.
type APIConnector struct {
	WebConnector
	APIToken  Option[forgedomain.ForgejoToken]
	cache     forgedomain.ProposalCache
	_client   OptionalMutable[forgejo.Client] // don't use directly, call .getClient()
	log       print.Logger
	remoteURL giturl.Parts
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &apiConnector // type check

func (self *APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if proposal := self.cache.BySourceTarget(branch, target); proposal.IsSome() {
		return proposal, nil
	}
	result, err := self.findProposalAtForge(branch, target)
	self.cache.SetOption(result)
	return result, err
}

func (self *APIConnector) findProposalAtForge(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	client, err := self.getClient()
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	openPullRequests, _, err := client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
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
	proposalDatas := parsePullRequests(pullRequests)
	switch len(proposalDatas) {
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
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &apiConnector // type check

func (self *APIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	client, err := self.getClient()
	if err != nil {
		return []forgedomain.Proposal{}, err
	}
	openPullRequests, _, err := client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 50,
		},
		State: forgejo.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	pullRequests := filterPullRequests2(openPullRequests, branch)
	result := make([]forgedomain.Proposal, len(pullRequests))
	for p, pullRequest := range pullRequests {
		proposalData := parsePullRequest(pullRequest)
		self.log.Success(proposalData.Target.String())
		proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeForgejo}
		result[p] = proposal
	}
	if len(result) == 0 {
		self.log.Success("none")
	}
	return result, nil
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = &apiConnector // type check

func (self *APIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.ForgeForgejoMergingViaAPI, colors.BoldGreen().Styled(strconv.Itoa(number)))
	client, err := self.getClient()
	if err != nil {
		return err
	}
	_, _, err = client.MergePullRequest(self.Organization, self.Repository, int64(number), forgejo.MergePullRequestOption{
		Style:   forgejo.MergeStyleSquash,
		Title:   commitMessageParts.Title.String(),
		Message: commitMessageParts.Body,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	self.log.Start(messages.APIProposalLookupStart)
	_, _, err = client.GetPullRequest(self.Organization, self.Repository, int64(number))
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = &apiConnector // type check

func (self *APIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody string) error {
	data := proposalData.Data()
	client, err := self.getClient()
	if err != nil {
		return err
	}
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err = client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), forgejo.EditPullRequestOption{
		Body: newBody,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &apiConnector // type check

func (self *APIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	client, err := self.getClient()
	if err != nil {
		return err
	}
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(targetName))
	_, _, err = client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), forgejo.EditPullRequestOption{
		Base: targetName,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = &apiConnector // type check

func (self *APIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	client, err := self.getClient()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	user, _, err := client.GetMyUserInfo()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
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

func (self *APIConnector) getClient() (*forgejo.Client, error) {
	if client, hasClient := self._client.Get(); hasClient {
		return client, nil
	}
	forgejoClient, err := forgejo.NewClient("https://"+self.remoteURL.Host, forgejo.SetToken(self.APIToken.GetOrZero().String()))
	if err != nil {
		return nil, err
	}
	self._client = MutableSome(forgejoClient)
	return forgejoClient, nil
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

func filterPullRequests2(pullRequests []*forgejo.PullRequest, branch gitdomain.LocalBranchName) []*forgejo.PullRequest {
	result := []*forgejo.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Head.Name == branch.String() {
			result = append(result, pullRequest)
		}
	}
	return result
}

func parsePullRequest(pullRequest *forgejo.PullRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		Active:       pullRequest.State == forgejo.StateOpen,
		MergeWithAPI: pullRequest.Mergeable,
		Number:       int(pullRequest.Index),
		Source:       gitdomain.NewLocalBranchName(pullRequest.Head.Ref),
		Target:       gitdomain.NewLocalBranchName(pullRequest.Base.Ref),
		Title:        pullRequest.Title,
		Body:         NewOption(pullRequest.Body),
		URL:          pullRequest.HTMLURL,
	}
}

func parsePullRequests(pullRequests []*forgejo.PullRequest) []forgedomain.Proposal {
	result := []forgedomain.Proposal{}
	for _, pullRequest := range pullRequests {
		proposalData := parsePullRequest(pullRequest)
		proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeForgejo}
		result = append(result, proposal)
	}
	return result
}
