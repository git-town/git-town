package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ProposalUpdateBody struct {
	Proposal    forgedomain.Proposal
	UpdatedBody string
}

func (self *ProposalUpdateBody) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	apiConnector, isAPIConnector := connector.(forgedomain.APIConnector)
	if !isAPIConnector {
		return errors.New(messages.UpdateProposalBodyUnsupported)
	}
	return apiConnector.UpdateProposalBody(self.Proposal.Data, self.UpdatedBody)
}
