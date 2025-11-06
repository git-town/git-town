package gitea

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"golang.org/x/oauth2"
)

// type checks
var (
	apiConnector AuthConnector
	_            forgedomain.Connector = apiConnector
)

// AuthConnector provides access to the gitea API.
type AuthConnector struct {
	WebConnector
	APIToken  forgedomain.GiteaToken
	RemoteURL giturl.Parts
	_client   OptionalMutable[gitea.Client] // don't use directly, call .getClient()
	log       print.Logger
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &apiConnector // type check

func (self *AuthConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	client, err := self.getClient()
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	self.log.Start(messages.APIProposalLookupStart)
	openPullRequests, _, err := client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
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
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitea}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(pullRequests), branch, target)
	}
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &apiConnector // type check

func (self *AuthConnector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	client, err := self.getClient()
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	openPullRequests, _, err := client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	pullRequests := FilterPullRequests2(openPullRequests, branch)
	switch len(pullRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		pullRequest := pullRequests[0]
		proposal := parsePullRequest(pullRequest)
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitea}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(pullRequests), branch)
	}
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = &apiConnector // type check

func (self *AuthConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	client, err := self.getClient()
	if err != nil {
		return err
	}
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	commitMessageParts := message.Parts()
	self.log.Start(messages.ForgeGitHubMergingViaAPI, colors.BoldGreen().Styled(strconv.Itoa(number)))
	_, _, err = client.MergePullRequest(self.Organization, self.Repository, int64(number), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
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

func (self *AuthConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody string) error {
	client, err := self.getClient()
	if err != nil {
		return err
	}
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err = client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), gitea.EditPullRequestOption{
		Body: &updatedBody,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &apiConnector // type check

func (self *AuthConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	client, err := self.getClient()
	if err != nil {
		return err
	}
	data := proposalData.Data()
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(targetName))
	_, _, err = client.EditPullRequest(self.Organization, self.Repository, int64(data.Number), gitea.EditPullRequestOption{
		Base: targetName,
	})
	self.log.Finished(err)
	return err
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = &apiConnector

func (self *AuthConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
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
	_, _, err = client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 1,
		},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.UserName),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}

func (self *AuthConnector) getClient() (*gitea.Client, error) {
	if client, hasClient := self._client.Get(); hasClient {
		return client, nil
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: self.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	giteaClient, err := gitea.NewClient("https://"+self.RemoteURL.Host, gitea.SetHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	self._client = MutableSome(giteaClient)
	return giteaClient, err
}
