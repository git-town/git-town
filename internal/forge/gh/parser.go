package gh

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	data := parsed[0]
	return Some(forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Active:       data.State == "open",
			Body:         NewOption(data.Body),
			MergeWithAPI: data.Mergeable == "MERGEABLE",
			Number:       data.Number,
			Source:       gitdomain.NewLocalBranchName(data.HeadRefName),
			Target:       gitdomain.NewLocalBranchName(data.BaseRefName),
			Title:        data.Title,
			URL:          data.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}), nil
}

// data returned by glab in JSON mode
type jsonData struct {
	BaseRefName string `json:"baseRefName"`
	Body        string `json:"body"`
	HeadRefName string `json:"headRefName"`
	Mergeable   string `json:"mergeable"`
	Number      int    `json:"number"`
	State       string `json:"state"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}
