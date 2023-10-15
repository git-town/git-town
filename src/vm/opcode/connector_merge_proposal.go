package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// ConnectorMergeProposal squash merges the branch with the given name into the current branch.
type ConnectorMergeProposal struct {
	Branch                    domain.LocalBranchName
	CommitMessage             string
	ProposalMessage           string
	enteredEmptyCommitMessage bool
	mergeError                error
	ProposalNumber            int
	undeclaredOpcodeMethods
}

func (op *ConnectorMergeProposal) CreateAbortProgram() []shared.Opcode {
	if op.enteredEmptyCommitMessage {
		return []shared.Opcode{&DiscardOpenChanges{}}
	}
	return []shared.Opcode{}
}

func (op *ConnectorMergeProposal) CreateAutomaticAbortError() error {
	if op.enteredEmptyCommitMessage {
		return fmt.Errorf(messages.ShipAbortedMergeError)
	}
	return op.mergeError
}

func (op *ConnectorMergeProposal) Run(args shared.RunArgs) error {
	commitMessage := op.CommitMessage
	//nolint:nestif
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		op.enteredEmptyCommitMessage = true
		err := args.Runner.Frontend.SquashMerge(op.Branch)
		if err != nil {
			return err
		}
		err = args.Runner.Backend.CommentOutSquashCommitMessage(op.ProposalMessage + "\n\n")
		if err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		err = args.Runner.Frontend.StartCommit()
		if err != nil {
			return err
		}
		commitMessage, err = args.Runner.Backend.LastCommitMessage()
		if err != nil {
			return err
		}
		err = args.Runner.Frontend.DeleteLastCommit()
		if err != nil {
			return err
		}
		op.enteredEmptyCommitMessage = false
	}
	_, op.mergeError = args.Connector.SquashMergeProposal(op.ProposalNumber, commitMessage)
	return op.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this opcode should cause the command to
// automatically abort if it errors.
func (op *ConnectorMergeProposal) ShouldAutomaticallyAbortOnError() bool {
	return true
}
