package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// UpdateProposalTarget updates the target of the proposal with the given number at the code hosting platform.
type UpdateProposalTarget struct {
	NewTarget      gitdomain.LocalBranchName
	ProposalNumber int
	undeclaredOpcodeMethods
}

func (self *UpdateProposalTarget) CreateAutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *UpdateProposalTarget) Run(args shared.RunArgs) error {
	return args.Connector.UpdateProposalTarget(self.ProposalNumber, self.NewTarget)
}

func (self *UpdateProposalTarget) ShouldAutomaticallyUndoOnError() bool {
	return true
}
