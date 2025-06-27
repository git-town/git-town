package gh

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

// data returned by glab in JSON mode
type jsonData struct {
	BaseRefName string `json:"baseRefName"`
	Body        string `json:"body"`
	HeadRefName string `json:"headRefName"`
	Mergeable   string `json:"mergeable"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

func (self jsonData) ToProposal() forgedomain.Proposal {
	return forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         NewOption(self.Body),
			MergeWithAPI: self.Mergeable == "MERGEABLE",
			Number:       self.Number,
			Source:       gitdomain.NewLocalBranchName(self.HeadRefName),
			Target:       gitdomain.NewLocalBranchName(self.BaseRefName),
			Title:        self.Title,
			URL:          self.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}
}
