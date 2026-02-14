package gh

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/github"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/pkg/colors"
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
	log      print.Logger
}

// ============================================================================
// browse the repo
// ============================================================================

func (self Connector) BrowseRepository(runner subshelldomain.Runner) error {
	return runner.Run("gh", "browse")
}

// ============================================================================
// create proposals
// ============================================================================

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
	self.log.Start(messages.APIProposalFindStart, branch, target)
	out, err := self.Backend.Query("gh", "pr", "list", "--head="+branch.String(), "--base="+target.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
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
		return Some(proposals[0]), nil
	default:
		self.log.Success("multiple")
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(proposals), branch, target)
	}
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = ghConnector // type-check

func (self Connector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	self.log.Start(messages.APIProposalSearchStart, branch.String())
	out, err := self.Backend.Query("gh", "pr", "list", "--head="+branch.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
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

var _ forgedomain.ProposalMerger = ghConnector // type-check

func (self Connector) SquashMergeProposal(number forgedomain.ProposalNumber, message gitdomain.CommitMessage) error {
	self.log.Start(messages.ForgeGithubMergingViaAPI, colors.BoldGreen().Styled("#"+number.String()))
	err := self.Frontend.Run("gh", "pr", "merge", "--squash", "--body="+message.String(), number.String())
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = ghConnector // type-check

func (self Connector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody gitdomain.ProposalBody) error {
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+proposalData.Data().Number.String()))
	err := self.Frontend.Run("gh", "pr", "edit", proposalData.Data().Number.String(), "--body="+updatedBody.String())
	self.log.Finished(err)
	return err
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = ghConnector // type-check

func (self Connector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	targetName := target.String()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+proposalData.Data().Number.String()), colors.BoldCyan().Styled(targetName))
	err := self.Frontend.Run("gh", "pr", "edit", proposalData.Data().Number.String(), "--base="+targetName)
	self.log.Finished(err)
	return err
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
	lines := stringslice.NonEmptyLines(output)
	parsePermissionsOnce.Do(func() {
		parsePermissionsLoggedInRegex = regexp.MustCompile(`Logged in to github.com account (\w+)`)
		parsePermissionsScopesRegex = regexp.MustCompile(`Token scopes: (.+)`)
	})
	for _, line := range lines {
		matches := parsePermissionsLoggedInRegex.FindStringSubmatch(line)
		if matches != nil {
			result.AuthenticatedUser = NewOption(matches[1])
			break
		}
	}
	if result.AuthenticatedUser.IsNone() {
		result.AuthenticationError = errors.New(messages.AuthenticationMissing)
	}
	for _, line := range lines {
		matches := parsePermissionsScopesRegex.FindStringSubmatch(line)
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

var (
	parsePermissionsOnce          sync.Once
	parsePermissionsLoggedInRegex *regexp.Regexp
	parsePermissionsScopesRegex   *regexp.Regexp
)
