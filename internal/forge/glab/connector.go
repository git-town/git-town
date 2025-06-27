package glab

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
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

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(self.findProposal)
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	return runner.Run("glab", "repo", "view", "--web")
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
	output, err := self.Backend.Query("glab", "auth", "status")
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
	out, err := self.Backend.Query("glab", "mr", "list", "--source-branch="+branch.String(), "--target-branch="+target.String(), "--output=json")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []glabData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil || len(parsed) == 0 {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromToFound, len(parsed), branch, target)
	}
	return Some(parsed[0].ToProposal()), nil
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("glab", "--source-branch="+branch.String(), "--output=json")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []glabData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil || len(parsed) == 0 {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf("found more than one pull request: %d", len(parsed))
	}
	return Some(parsed[0].ToProposal()), nil
}

func (self Connector) squashMergeProposal(number int, message gitdomain.CommitMessage) (err error) {
	return self.Frontend.Run("glab", "mr", "merge", "--squash", "--body="+message.String(), strconv.Itoa(number))
}

func (self Connector) updateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	return self.Frontend.Run("glab", "edit", strconv.Itoa(proposalData.Data().Number), "--base="+target.String())
}

func ParsePermissionsOutput(output string) forgedomain.VerifyConnectionResult {
	result := forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
	lines := strings.Split(output, "\n")
	regex := regexp.MustCompile(`Logged in to \S+ as (\S+) `)
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
	return result
}

type glabData struct {
	Description  string `json:"description"`
	Mergeable    string `json:"detailed_merge_status"` //nolint:tagliatelle
	Number       int    `json:"iid"`                   //nolint:tagliatelle
	SourceBranch string `json:"source_branch"`         //nolint:tagliatelle
	TargetBranch string `json:"target_branch"`         //nolint:tagliatelle
	Title        string `json:"title"`
	URL          string `json:"web_url"` //nolint:tagliatelle
}

func (data glabData) ToProposal() forgedomain.Proposal {
	return forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         NewOption(data.Description),
			MergeWithAPI: data.Mergeable == "mergeable",
			Number:       data.Number,
			Source:       gitdomain.NewLocalBranchName(data.SourceBranch),
			Target:       gitdomain.NewLocalBranchName(data.TargetBranch),
			Title:        data.Title,
			URL:          data.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}
}
