package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/vm/shared"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// ConnectorProposalMerge squash merges the branch with the given name into the current branch.
type ConnectorProposalMerge struct {
	Branch                    gitdomain.LocalBranchName
	CommitMessage             Option[gitdomain.CommitMessage]
	ProposalMessage           string
	ProposalNumber            int
	enteredEmptyCommitMessage bool
	mergeError                error
	undeclaredOpcodeMethods   `exhaustruct:"optional"`
}

func (self *ConnectorProposalMerge) AbortProgram() []shared.Opcode {
	if self.enteredEmptyCommitMessage {
		return []shared.Opcode{&ChangesDiscard{}}
	}
	return []shared.Opcode{}
}

func (self *ConnectorProposalMerge) AutomaticUndoError() error {
	if self.enteredEmptyCommitMessage {
		return errors.New(messages.ShipAbortedMergeError)
	}
	return self.mergeError
}

func (self *ConnectorProposalMerge) Run(args shared.RunArgs) error {
	commitMessage, hasCommitMessage := self.CommitMessage.Get()
	if !hasCommitMessage {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		self.enteredEmptyCommitMessage = true
		err := args.Git.SquashMerge(args.Frontend, self.Branch)
		if err != nil {
			return err
		}
		err = args.Git.CommentOutSquashCommitMessage(Some(self.ProposalMessage + "\n\n"))
		if err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		err = args.Git.CommitStart(args.Frontend)
		if err != nil {
			return err
		}
		commitMessage, err = args.Git.HeadCommitMessage(args.Backend)
		if err != nil {
			return err
		}
		err = args.Git.DeleteLastCommit(args.Frontend)
		if err != nil {
			return err
		}
		self.enteredEmptyCommitMessage = false
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	squashMergeProposal, canSquashMergeProposal := connector.SquashMergeProposalFn().Get()
	if !canSquashMergeProposal {
		return errors.New(messages.ShipAPIConnectorUnsupported)
	}
	self.mergeError = squashMergeProposal(self.ProposalNumber, commitMessage)
	return self.mergeError
}

// ShouldUndoOnError returns whether this opcode should cause the command to
// automatically undo if it errors.
func (self *ConnectorProposalMerge) ShouldUndoOnError() bool {
	return true
}
