package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ProposalUpdateBody struct {
	Proposal    forgedomain.Proposal
	UpdatedBody string
}

func (self *ProposalUpdateBody) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalUpdateBodyProblem, self.Proposal.Data.Data().Number)
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

func (self *ProposalUpdateBody) ShouldUndoOnError() bool {
	return true
}

func (self *ProposalUpdateBody) UndoExternalChanges() []shared.Opcode {
	return []shared.Opcode{
		&ProposalUpdateBody{
			Proposal:    self.Proposal,
			UpdatedBody: self.Proposal.Data.Data().Body.GetOrZero(),
		},
	}
}
