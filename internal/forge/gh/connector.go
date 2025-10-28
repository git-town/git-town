package gh

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/github"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	ghConnector Connector
	_           forgedomain.Connector = ghConnector
)

// Connector talks to the GitHub API through the "gh" executable.
type Connector struct {
	Backend  subshelldomain.Querier
	Frontend subshelldomain.Runner
}

func (self Connector) BrowseRepository(runner subshelldomain.Runner) error {
	return runner.Run("gh", "browse")
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	args := []string{"pr", "create", "--head=" + data.Branch.String(), "--base=" + data.ParentBranch.String()}
	if title, hasTitle := data.ProposalTitle.Get(); hasTitle {
		args = append(args, "--title="+title.String())
	}
	if body, hasBody := data.ProposalBody.Get(); hasBody {
		args = append(args, "--body="+body.String())
	}
	if err := self.Frontend.Run("gh", args...); err != nil {
		return err
	}
	// check if the proposal exists
	proposal, err := self.FindProposal(data.Branch, data.ParentBranch)
	if err != nil {
		return err
	}
	if proposal.IsNone() {
		return nil
	}
	return self.Frontend.Run("gh", "pr", "view", "--web")
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return github.DefaultProposalMessage(data)
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = ghConnector // type-check

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("gh", "pr", "list", "--head="+branch.String(), "--base="+target.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	return ParseJSONOutput(out, branch)
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = ghConnector // type-check

func (self Connector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("gh", "--head="+branch.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	return ParseJSONOutput(out, branch)
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = ghConnector // type-check

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	return self.Frontend.Run("gh", "pr", "merge", "--squash", "--body="+message.String(), strconv.Itoa(number))
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = ghConnector // type-check

func (self Connector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody string) error {
	return self.Frontend.Run("gh", "pr", "edit", strconv.Itoa(proposalData.Data().Number), "--body="+updatedBody)
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = ghConnector // type-check

func (self Connector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	return self.Frontend.Run("gh", "pr", "edit", strconv.Itoa(proposalData.Data().Number), "--base="+target.String())
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = ghConnector // type check

func (self Connector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	output, err := self.Backend.Query("gh", "auth", "status", "--active")
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	return ParsePermissionsOutput(output)
}

func ParsePermissionsOutput(output string) forgedomain.VerifyCredentialsResult {
	result := forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
	lines := strings.Split(output, "\n")
	regex := regexp.MustCompile(`Logged in to github.com account (\w+)`)
	for _, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if matches != nil {
			result.AuthenticatedUser = NewOption(matches[1])
			break
		}
	}
	if result.AuthenticatedUser.IsNone() {
		result.AuthenticationError = errors.New(messages.AuthenticationMissing)
	}
	regex = regexp.MustCompile(`Token scopes: (.+)`)
	for _, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if matches != nil {
			parts := strings.Split(matches[1], ", ")
			if slices.Contains(parts, "'repo'") {
				break
			}
			result.AuthorizationError = fmt.Errorf(messages.AuthorizationMissing, parts)
		}
	}
	return result
}
