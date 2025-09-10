package glab

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func (self Connector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedDescription string) error {
	return self.Frontend.Run("glab", "mr", "update", strconv.Itoa(proposalData.Data().Number), "--description="+updatedDescription)
}

func (self Connector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	return self.Frontend.Run("glab", "mr", "update", strconv.Itoa(proposalData.Data().Number), "--target-branch="+target.String())
}
