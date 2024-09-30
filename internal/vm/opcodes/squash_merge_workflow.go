package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// SquashMergeWorkflow squash merges the branch with the given name into the current branch.
type SquashMergeWorkflow struct {
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	Parent        gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SquashMergeWorkflow) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&DiscardOpenChanges{},
	}
}

func (self *SquashMergeWorkflow) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *SquashMergeWorkflow) Run(args shared.RunArgs) error {
	// TODO: extract into separate opcodes for Git resilience
	// Possible create a SquashMergeProgram function that returns these opcodes
	branchAuthors, err := args.Git.BranchAuthors(args.Backend, self.Branch, self.Parent)
	if err != nil {
		return err
	}
	author, aborted, err := dialog.SelectSquashCommitAuthor(self.Branch, branchAuthors, args.DialogTestInputs.Next())
	if err != nil {
		return fmt.Errorf(messages.SquashCommitAuthorProblem, err)
	}
	if aborted {
		return errors.New("aborted by user")
	}
	repoAuthor := args.Config.Author()
	if repoAuthor == author {
		author = ""
	}
	args.PrependOpcodes(&SquashMerge{Branch: self.Branch})
	if !args.Config.DryRun {
		args.PrependOpcodes(&CommentOutSquashCommitMessage{})
	}
	args.PrependOpcodes(&Commit{
		Message:                        self.CommitMessage,
		FallbackToDefaultCommitMessage: false,
		undeclaredOpcodeMethods:        undeclaredOpcodeMethods{},
	})
	squashedCommitSHA, err := args.Git.SHAForBranch(args.Backend, self.Parent.BranchName())
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}

func (self *SquashMergeWorkflow) ShouldUndoOnError() bool {
	return true
}
