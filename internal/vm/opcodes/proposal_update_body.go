package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ProposalUpdateBody struct {
	Proposal    forgedomain.Proposal
	UpdatedBody gitdomain.ProposalBody
}

func (self *ProposalUpdateBody) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	proposalBodyUpdater, canUpdateProposalBody := connector.(forgedomain.ProposalBodyUpdater)
	if !canUpdateProposalBody {
		return errors.New(messages.UpdateProposalBodyUnsupported)
	}
	return proposalBodyUpdater.UpdateProposalBody(self.Proposal.Data, self.UpdatedBody)
}
