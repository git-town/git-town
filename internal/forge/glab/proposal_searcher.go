package glab

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func (self Connector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.Backend.Query("glab", "--source-branch="+branch.String(), "--output=json")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	return ParseJSONOutput(out, branch)
}
