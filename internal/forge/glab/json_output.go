package glab

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func ParseJSONOutput(output string, branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	var parsed []jsonData
	err := json.Unmarshal([]byte(output), &parsed)
	if err != nil || len(parsed) == 0 {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(parsed), branch)
	}
	return Some(parsed[0].ToProposal()), nil
}

type ParseJSONOutputArgs struct {
	output             string
	multipleFoundError error
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

func (self jsonData) ToProposal() forgedomain.Proposal {
	return forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         NewOption(self.Description),
			MergeWithAPI: self.Mergeable == "mergeable",
			Number:       self.Number,
			Source:       gitdomain.NewLocalBranchName(self.SourceBranch),
			Target:       gitdomain.NewLocalBranchName(self.TargetBranch),
			Title:        self.Title,
			URL:          self.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}
}
