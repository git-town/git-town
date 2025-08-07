package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ProposalUpdateBody struct {
	Proposal                forgedomain.Proposal
	UpdatedBody             string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateBody) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}

	updateProposalBodyFn, canUpdateProposalBody := connector.UpdateProposalBodyFn().Get()
	if !canUpdateProposalBody {
		return errors.New(messages.UpdateProposalBodyUnsupported)
	}
	return updateProposalBodyFn(self.Proposal.Data, self.UpdatedBody)
}
