package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalUpdateBody struct {
	Proposal                forgedomain.ProposalInterface
	UpdatedBody             Option[string]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateBody) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}

	updateProposalBodyFn, hasUpdateProposalBody := connector.UpdateProposalBodyFn().Get()
	if !hasUpdateProposalBody {
		return errors.New(messages.UpdateProposalBodyUnsupported)
	}
	return updateProposalBodyFn(self.Proposal, self.UpdatedBody.GetOrDefault())
}
