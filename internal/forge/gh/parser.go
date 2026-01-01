package gh

import (
	"encoding/json"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func ParseJSONOutput(output string) ([]forgedomain.Proposal, error) {
	var parsed []jsonData
	err := json.Unmarshal([]byte(output), &parsed)
	if err != nil || len(parsed) == 0 {
		return []forgedomain.Proposal{}, err
	}
	result := make([]forgedomain.Proposal, len(parsed))
	for d, data := range parsed {
		result[d] = forgedomain.Proposal{
			Data: forgedomain.ProposalData{
				Active:       data.State == "open",
				Body:         NewOption(data.Body),
				MergeWithAPI: data.Mergeable == "MERGEABLE",
				Number:       data.Number,
				Source:       data.HeadRefName,
				Target:       data.BaseRefName,
				Title:        data.Title,
				URL:          data.URL,
			},
			ForgeType: forgedomain.ForgeTypeGithub,
		}
	}
	return result, nil
}

// data returned by glab in JSON mode
type jsonData struct {
	BaseRefName gitdomain.LocalBranchName `json:"baseRefName"`
	Body        gitdomain.ProposalBody    `json:"body"`
	HeadRefName gitdomain.LocalBranchName `json:"headRefName"`
	Mergeable   string                    `json:"mergeable"`
	Number      int                       `json:"number"`
	State       string                    `json:"state"`
	Title       gitdomain.ProposalTitle   `json:"title"`
	URL         string                    `json:"url"`
}
