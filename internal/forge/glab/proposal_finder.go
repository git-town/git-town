package glab

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

var _ forgedomain.ProposalFinder = glabConnector

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("glab", "mr", "list", "--source-branch="+branch.String(), "--target-branch="+target.String(), "--output=json")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	return ParseJSONOutput(out, branch)
}

func ParseJSONOutput(output string, branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	var parsed []jsonData
	err := json.Unmarshal([]byte(output), &parsed)
	if err != nil || len(parsed) == 0 {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(parsed), branch)
	}
	return Some(createProposal(parsed[0])), nil
}

type jsonData struct {
	Description  string `json:"description"`
	Mergeable    string `json:"detailed_merge_status"` //nolint:tagliatelle
	Number       int    `json:"iid"`                   //nolint:tagliatelle
	SourceBranch string `json:"source_branch"`         //nolint:tagliatelle
	TargetBranch string `json:"target_branch"`         //nolint:tagliatelle
	Title        string `json:"title"`
	URL          string `json:"web_url"` //nolint:tagliatelle
}

func createProposal(data jsonData) forgedomain.Proposal {
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
		ForgeType: forgedomain.ForgeTypeGitLab,
	}
}
