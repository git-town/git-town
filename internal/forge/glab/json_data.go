package glab

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type jsonData struct {
	Description  string `json:"description"`
	Mergeable    string `json:"detailed_merge_status"` //nolint:tagliatelle
	Number       int    `json:"iid"`                   //nolint:tagliatelle
	SourceBranch string `json:"source_branch"`         //nolint:tagliatelle
	TargetBranch string `json:"target_branch"`         //nolint:tagliatelle
	Title        string `json:"title"`
	URL          string `json:"web_url"` //nolint:tagliatelle
}

func (data jsonData) ToProposal() forgedomain.Proposal {
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
