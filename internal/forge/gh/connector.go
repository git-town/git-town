package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	Backend  subshelldomain.Querier
	Frontend subshelldomain.Runner
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	args := []string{"pr", "create", "--base=" + data.ParentBranch.String(), "--head=" + data.Branch.String()}
	if title, hasTitle := data.ProposalTitle.Get(); hasTitle {
		args = append(args, "--title="+title.String())
	}
	if body, hasBody := data.ProposalBody.Get(); hasBody {
		args = append(args, "--body="+body.String())
	}
	err := self.Frontend.Run("gh", args...)
	if err != nil {
		return err
	}
	return self.Frontend.Run("gh", "pr", "view", "--web")
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return github.DefaultProposalMessage(data)
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(self.findProposal)
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	return runner.Run("gh", "browse")
}

func (self Connector) SearchProposalFn() Option[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(self.searchProposal)
}

func (self Connector) SquashMergeProposalFn() Option[func(int, gitdomain.CommitMessage) (err error)] {
	return Some(self.squashMergeProposal)
}

func (self Connector) UpdateProposalSourceFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error] {
	return None[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error] {
	return Some(self.updateProposalTarget)
}

func (self Connector) VerifyConnection() forgedomain.VerifyConnectionResult {
	output, err := self.Backend.Query("gh", "auth", "status", "--active")
	if err != nil {
		return forgedomain.VerifyConnectionResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	return ParsePermissionsOutput(output)
}

func (self Connector) findProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("gh", "pr", "list", "--head="+branch.String(), "--base="+target.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []ghData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil || len(parsed) == 0 {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(parsed), branch, target)
	}
	pr := parsed[0]
	proposal := forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         NewOption(pr.Body),
			MergeWithAPI: pr.Mergeable == "MERGEABLE",
			Number:       pr.Number,
			Source:       gitdomain.NewLocalBranchName(pr.HeadRefName),
			Target:       gitdomain.NewLocalBranchName(pr.BaseRefName),
			Title:        pr.Title,
			URL:          pr.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}
	return Some(proposal), nil
}

type ghData struct {
	BaseRefName string `json:"baseRefName"`
	Body        string `json:"body"`
	HeadRefName string `json:"headRefName"`
	Mergeable   string `json:"mergeable"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("gh", "--head="+branch.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []ghData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil || len(parsed) == 0 {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf("found more than one pull request: %d", len(parsed))
	}
	pr := parsed[0]
	proposal := forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         NewOption(pr.Body),
			MergeWithAPI: pr.Mergeable == "MERGEABLE",
			Number:       pr.Number,
			Source:       gitdomain.NewLocalBranchName(pr.HeadRefName),
			Target:       gitdomain.NewLocalBranchName(pr.BaseRefName),
			Title:        pr.Title,
			URL:          pr.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}
	return Some(proposal), nil
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) (err error) {
	return self.Frontend.Run("gh", "pr", "merge", "--squash", "--body="+message.String(), strconv.Itoa(number))
}

func (self Connector) updateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	return self.Frontend.Run("gh", "edit", strconv.Itoa(proposalData.Data().Number), "--base="+target.String())
}

func ParsePermissionsOutput(output string) forgedomain.VerifyConnectionResult {
	result := forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
	lines := strings.Split(output, "\n")
	regex := regexp.MustCompile(`Logged in to github.com account (\w+)`)
	found := false
	for _, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if matches != nil {
			result.AuthenticatedUser = NewOption(matches[1])
			found = true
			break
		}
	}
	if !found {
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
