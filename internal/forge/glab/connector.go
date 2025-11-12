package glab

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/gitlab"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
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
	out, err := self.Backend.Query("glab", "mr", "list", "--source-branch="+branch.String(), "--target-branch="+target.String(), "--output=json")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	proposals, err := ParseJSONOutput(out)
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	switch len(proposals) {
	case 0:
		return None[forgedomain.Proposal](), nil
	case 1:
		return Some(proposals[0]), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(proposals), branch)
	}
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = glabConnector // type check

func (self Connector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	out, err := self.Backend.Query("glab", "--source-branch="+branch.String(), "--output=json")
	if err != nil {
		return []forgedomain.Proposal{}, err
	}
	return ParseJSONOutput(out)
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = glabConnector // type check

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	return self.Frontend.Run("glab", "mr", "merge", "--squash", "--body="+message.String(), strconv.Itoa(number))
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = glabConnector // type check

func (self Connector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedDescription string) error {
	return self.Frontend.Run("glab", "mr", "update", strconv.Itoa(proposalData.Data().Number), "--description="+updatedDescription)
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = glabConnector // type check

func (self Connector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	return self.Frontend.Run("glab", "mr", "update", strconv.Itoa(proposalData.Data().Number), "--target-branch="+target.String())
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
