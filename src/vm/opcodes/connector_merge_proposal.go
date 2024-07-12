package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// ConnectorMergeProposal squash merges the branch with the given name into the current branch.
type ConnectorMergeProposal struct {
	Branch                    gitdomain.LocalBranchName
	CommitMessage             Option[gitdomain.CommitMessage]
	ProposalMessage           string
	ProposalNumber            int
	enteredEmptyCommitMessage bool
	mergeError                error
}

func (self *ConnectorMergeProposal) CreateAbortProgram() []shared.Opcode {
	if self.enteredEmptyCommitMessage {
		return []shared.Opcode{&DiscardOpenChanges{}}
	}
	return []shared.Opcode{}
}

func (self *ConnectorMergeProposal) CreateAutomaticUndoError() error {
	if self.enteredEmptyCommitMessage {
		return errors.New(messages.ShipAbortedMergeError)
	}
	return self.mergeError
}

func (self *ConnectorMergeProposal) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{self}
}

func (self *ConnectorMergeProposal) Run(args shared.RunArgs) error {
	commitMessage, hasCommitMessage := self.CommitMessage.Get()
	//nolint:nestif
	if !hasCommitMessage {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		self.enteredEmptyCommitMessage = true
		err := args.Git.SquashMerge(args.Frontend, self.Branch)
		if err != nil {
			return err
		}
		err = args.Git.CommentOutSquashCommitMessage(self.ProposalMessage + "\n\n")
		if err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		err = args.Git.StartCommit(args.Frontend)
		if err != nil {
			return err
		}
		commitMessage, err = args.Git.LastCommitMessage(args.Backend)
		if err != nil {
			return err
		}
		err = args.Git.DeleteLastCommit(args.Frontend)
		if err != nil {
			return err
		}
		self.enteredEmptyCommitMessage = false
	}
	if connector, hasConnector := args.Connector.Get(); hasConnector {
		self.mergeError = connector.SquashMergeProposal(self.ProposalNumber, commitMessage)
	} else {
		return hostingdomain.UnsupportedServiceError()
	}
	return self.mergeError
}

// ShouldAutomaticallyUndoOnError returns whether this opcode should cause the command to
// automatically undo if it errors.
func (self *ConnectorMergeProposal) ShouldAutomaticallyUndoOnError() bool {
	return true
}
