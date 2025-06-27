package gh

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

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
