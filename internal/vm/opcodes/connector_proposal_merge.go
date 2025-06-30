package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// ConnectorProposalMerge squash merges the branch with the given name into the current branch.
type ConnectorProposalMerge struct {
	Branch                    gitdomain.LocalBranchName
	CommitMessage             Option[gitdomain.CommitMessage]
	Proposal                  forgedomain.Proposal
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
		return errors.New(messages.ShipAbortMergeError)
	}
	return self.mergeError
}

func (self *ConnectorProposalMerge) Run(args shared.RunArgs) error {
	commitMessage, hasCommitMessage := self.CommitMessage.Get()
	proposalData := self.Proposal.Data.Data()
	if !hasCommitMessage {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		self.enteredEmptyCommitMessage = true
		err := args.Git.SquashMerge(args.Frontend, self.Branch)
		if err != nil {
			return err
		}
		err = args.Git.CommentOutSquashCommitMessage(Some(forgedomain.CommitBody(proposalData, proposalData.Title) + "\n\n"))
		if err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		err = args.Git.CommitStart(args.Frontend)
		if err != nil {
			return err
		}
		commitMessage, err = args.Git.CommitMessage(args.Backend, "HEAD")
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
	self.mergeError = squashMergeProposal(proposalData.Number, commitMessage)
	return self.mergeError
}

// ShouldUndoOnError returns whether this opcode should cause the command to
// automatically undo if it errors.
func (self *ConnectorProposalMerge) ShouldUndoOnError() bool {
	return true
}
