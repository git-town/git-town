package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/messages"
	"github.com/git-town/git-town/v24/internal/vm/shared"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// ConnectorProposalMerge squash merges the branch with the given name into the current branch.
type ConnectorProposalMerge struct {
	Branch                    gitdomain.LocalBranchName
	CommitMessage             Option[gitdomain.CommitMessage]
	EnterMessage              configdomain.ShipEnterMessage
	Proposal                  forgedomain.Proposal
	enteredEmptyCommitMessage bool
	mergeError                error
}

func (self *ConnectorProposalMerge) Abort() []shared.Opcode {
	if self.enteredEmptyCommitMessage {
		return []shared.Opcode{&ChangesDiscard{}}
	}
	return []shared.Opcode{}
}

func (self *ConnectorProposalMerge) AutomaticUndoError() error {
	if self.enteredEmptyCommitMessage {
		return errors.New(messages.ShipExitMergeError)
	}
	return self.mergeError
}

func (self *ConnectorProposalMerge) Run(args shared.RunArgs) error {
	proposalData := self.Proposal.Data.Data()
	commitMessage := self.CommitMessage
	if commitMessage.IsNone() && self.EnterMessage.ShouldEnterMessage() {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		self.enteredEmptyCommitMessage = true
		if err := args.Git.SquashMerge(args.Frontend, self.Branch); err != nil {
			return err
		}
		if err := args.Git.CommentOutSquashCommitMessage(args.Backend, Some(forgedomain.CommitBody(proposalData, proposalData.Title.String())+"\n\n")); err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
		if err := args.Git.CommitStart(args.Frontend); err != nil {
			return err
		}
		enteredMessage, err := args.Git.CommitMessage(args.Backend, "HEAD")
		if err != nil {
			return err
		}
		if err := args.Git.DeleteLastCommit(args.Frontend); err != nil {
			return err
		}
		self.enteredEmptyCommitMessage = false
		commitMessage = Some(enteredMessage)
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	proposalMerger, canMergeProposals := connector.(forgedomain.ProposalMerger)
	if !canMergeProposals {
		return errors.New(messages.ShipAPIConnectorUnsupported)
	}
	// When no commit message is given, the forge determines the squash commit message.
	self.mergeError = proposalMerger.SquashMergeProposal(proposalData.Number, commitMessage)
	return self.mergeError
}

// ShouldUndoOnError returns whether this opcode should cause the command to
// automatically undo if it errors.
func (self *ConnectorProposalMerge) ShouldUndoOnError() bool {
	return true
}
