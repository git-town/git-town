package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/shared"
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

func (self *ConnectorMergeProposal) CreateAbortProgram() []shared.Opcode {
	if self.enteredEmptyCommitMessage {
		return []shared.Opcode{&DiscardOpenChanges{}}
	}
	return []shared.Opcode{}
}

func (self *ConnectorMergeProposal) CreateAutomaticAbortError() error {
	if self.enteredEmptyCommitMessage {
		return fmt.Errorf(messages.ShipAbortedMergeError)
	}
	return self.mergeError
}

func (self *ConnectorMergeProposal) Run(args shared.RunArgs) error {
	commitMessage := self.CommitMessage
	//nolint:nestif
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		self.enteredEmptyCommitMessage = true
		err := args.Runner.Frontend.SquashMerge(self.Branch)
		if err != nil {
			return err
		}
		err = args.Runner.Backend.CommentOutSquashCommitMessage(self.ProposalMessage + "\n\n")
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
		self.enteredEmptyCommitMessage = false
	}
	_, self.mergeError = args.Connector.SquashMergeProposal(self.ProposalNumber, commitMessage)
	return self.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this opcode should cause the command to
// automatically abort if it errors.
func (self *ConnectorMergeProposal) ShouldAutomaticallyAbortOnError() bool {
	return true
}
