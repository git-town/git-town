package glab

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/gitlab"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	glabConnector Connector
	_             forgedomain.Connector = glabConnector
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	Backend  subshelldomain.Querier
	Frontend subshelldomain.Runner
	log      print.Logger
}

// ============================================================================
// browse the repo
// ============================================================================

func (self Connector) BrowseRepository(runner subshelldomain.Runner) error {
	return runner.Run("glab", "repo", "view", "--web")
}

// ============================================================================
// create proposals
// ============================================================================

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	args := []string{"mr", "create", "--source-branch=" + data.Branch.String(), "--target-branch=" + data.ParentBranch.String()}
	title, hasTitle := data.ProposalTitle.Get()
	if hasTitle {
		args = append(args, "--title="+title.String())
	}
	body, hasBody := data.ProposalBody.Get()
	if hasBody {
		args = append(args, "--description="+body.String())
	}
	if !hasTitle || !hasBody {
		args = append(args, "--fill")
	}
	args = append(args, "--web")
	return self.Frontend.Run("glab", args...)
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return gitlab.DefaultProposalMessage(data)
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = glabConnector // type check

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalFindStart, branch, target)
	out, err := self.Backend.Query("glab", "mr", "list", "--source-branch="+branch.String(), "--target-branch="+target.String(), "--output=json")
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	proposals, err := ParseJSONOutput(out)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	switch len(proposals) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := proposals[0]
		self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+proposal.Data.Data().Number.String()), proposal.Data.Data().Title))
		return Some(proposal), nil
	default:
		self.log.Success("multiple")
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(proposals), branch)
	}
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = glabConnector // type check

func (self Connector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	self.log.Start(messages.APIProposalSearchStart, branch.String())
	out, err := self.Backend.Query("glab", "mr", "list", "--source-branch="+branch.String(), "--output=json")
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	proposals, err := ParseJSONOutput(out)
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	ids := make([]string, len(proposals))
	for p, proposal := range proposals {
		ids[p] = colors.BoldGreen().Styled(fmt.Sprintf("#%d", proposal.Data.Data().Number))
	}
	self.log.Log(strings.Join(ids, ", "))
	return proposals, err
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = glabConnector // type check

func (self Connector) SquashMergeProposal(number forgedomain.ProposalNumber, message gitdomain.CommitMessage) error {
	self.log.Start(messages.ForgeGithubMergingViaAPI, colors.BoldGreen().Styled("#"+number.String()))
	err := self.Frontend.Run("glab", "mr", "merge", "--squash", "--body="+message.String(), number.String())
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = glabConnector // type check

func (self Connector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedDescription gitdomain.ProposalBody) error {
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+proposalData.Data().Number.String()))
	err := self.Frontend.Run("glab", "mr", "update", proposalData.Data().Number.String(), "--description="+updatedDescription.String())
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = glabConnector // type check

func (self Connector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+proposalData.Data().Number.String()), colors.BoldCyan().Styled(targetName))
	err := self.Frontend.Run("glab", "mr", "update", proposalData.Data().Number.String(), "--target-branch="+target.String())
	self.log.Finished(err)
	return err
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = glabConnector

func (self Connector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	output, err := self.Backend.Query("glab", "auth", "status")
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	return ParsePermissionsOutput(output)
}
