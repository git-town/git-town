package undo

import (
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/undo/undostash"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// provides the program to undo the given runstate
func createProgram(args ExecuteArgs) program.Program {
	undoProgram := program.Program{}
	undoProgram.AddProgram(args.RunState.AbortProgram)
	if !args.RunState.IsFinished() && args.HasOpenChanges {
		// Open changes in the middle of an unfinished command will be undone as well.
		// To achieve this, we commit them here so that they are gone when the branch is reset to the original SHA.
		undoProgram.Add(&opcodes.CommitOpenChanges{})
	}

	// undo branch changes
	branchSpans := undobranches.NewBranchSpans(args.RunState.BeginBranchesSnapshot, args.RunState.EndBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	undoBranchesProgram := branchChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		Config:                   args.FullConfig,
		EndBranch:                args.RunState.EndBranchesSnapshot.Active,
		BeginBranch:              args.RunState.BeginBranchesSnapshot.Active,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	undoProgram.AddProgram(undoBranchesProgram)

	// undo config changes
	configSpans := undoconfig.NewConfigDiffs(args.RunState.BeginConfigSnapshot, args.RunState.EndConfigSnapshot)
	configUndoProgram := configSpans.UndoProgram()
	undoProgram.AddProgram(configUndoProgram)

	// undo stash changes
	stashDiff := undostash.NewStashDiff(args.RunState.BeginStashSize, args.InitialStashSize)
	undoStashProgram := stashDiff.Program()
	undoProgram.AddProgram(undoStashProgram)

	undoProgram.AddProgram(args.RunState.FinalUndoProgram)
	undoProgram.Add(&opcodes.Checkout{Branch: args.RunState.BeginBranchesSnapshot.Active})
	undoProgram.RemoveDuplicateCheckout()

	cmdhelpers.Wrap(&undoProgram, cmdhelpers.WrapOptions{
		DryRun:                   args.RunState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         args.RunState.IsFinished() && args.HasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{args.RunState.BeginBranchesSnapshot.Active},
	})
	return undoProgram
}
