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

// creates the program for undoing a program that finished
func CreateUndoFinishedProgram(args CreateUndoProgramArgs) program.Program {
	result := program.Program{}
	// if there is a pending operation --> abort it
	result.AddProgram(args.RunState.AbortProgram)
	if !args.RunState.IsFinished() && args.HasOpenChanges {
		// Open changes in the middle of an unfinished command will be undone as well.
		// To achieve this, we commit them here so that they are gone when the branch is reset to the original SHA.
		result.Add(&opcodes.CommitOpenChanges{})
	}
	// undo branch changes
	branchSpans := undobranches.NewBranchSpans(args.RunState.BeginBranchesSnapshot, args.RunState.EndBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	undoBranchesProgram := branchChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		BeginBranch:              args.RunState.BeginBranchesSnapshot.Active,
		Config:                   &args.Run.FullConfig,
		EndBranch:                args.RunState.EndBranchesSnapshot.Active,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	result.AddProgram(undoBranchesProgram)
	// undo config changes
	configSpans := undoconfig.NewConfigDiffs(args.RunState.BeginConfigSnapshot, args.RunState.EndConfigSnapshot)
	result.AddProgram(configSpans.UndoProgram())
	// undo stash changes
	stashDiff := undostash.NewStashDiff(args.RunState.BeginStashSize, args.RunState.EndStashSize)
	result.AddProgram(stashDiff.Program())
	// wrap up
	result.AddProgram(args.RunState.FinalUndoProgram)
	result.Add(&opcodes.Checkout{Branch: args.RunState.BeginBranchesSnapshot.Active})
	result.RemoveDuplicateCheckout()
	cmdhelpers.Wrap(&result, cmdhelpers.WrapOptions{
		DryRun:                   args.RunState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         args.RunState.IsFinished() && args.HasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{args.RunState.BeginBranchesSnapshot.Active},
	})
	return result
}
