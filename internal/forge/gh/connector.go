package gh

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	runner Runner
	log    print.Logger
}

// NewConnector provides a fully configured gh.Connector instance
// if the current repo is hosted on GitHub, otherwise nil.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	return Connector{
		runner: args.Runner,
		log:    args.Log,
	}, nil
}

type NewConnectorArgs struct {
	Log    print.Logger
	Runner Runner
}

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return github.DefaultProposalMessage(data)
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(self.findProposal)
}

func (self Connector) findProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.runner.Query("gh", "pr", "list", "--head="+branch.String(), "--base="+target.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []ghData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) == 0 {
		return None[forgedomain.Proposal](), nil
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

type ghData struct {
	BaseRefName string `json:"baseRefName"`
	Body        string `json:"body"`
	HeadRefName string `json:"headRefName"`
	Mergeable   string `json:"mergeable"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	return runner.Run("gh", "browse")
}

func (self Connector) SearchProposalFn() Option[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(self.searchProposal)
}

func (self Connector) UpdateProposalSourceFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error] {
	return None[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error]()
}

func (self Connector) UpdateProposalTargetFn() Option[func(forgedomain.ProposalInterface, gitdomain.LocalBranchName, stringslice.Collector) error] {
	return Some(self.updateProposalTarget)
}

func (self Connector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.runner.Query("gh", "--head="+branch.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []ghData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) == 0 {
		return None[forgedomain.Proposal](), nil
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

func (self Connector) updateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName, _ stringslice.Collector) error {
	return self.runner.Run("gh", "edit", strconv.Itoa(proposalData.Data().Number), "--base="+target.String())
}
