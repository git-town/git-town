package gh

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func (self Connector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody string) error {
	return self.Frontend.Run("gh", "pr", "edit", strconv.Itoa(proposalData.Data().Number), "--body="+updatedBody)
}

func (self Connector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	return self.Frontend.Run("gh", "pr", "edit", strconv.Itoa(proposalData.Data().Number), "--base="+target.String())
}
