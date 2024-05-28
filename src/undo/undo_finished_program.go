package undo

import (
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/undo/undostash"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// creates the program for undoing a program that finished
func CreateUndoForFinishedProgram(args CreateUndoProgramArgs) program.Program {
	result := program.Program{}
	result.AddProgram(args.RunState.AbortProgram)
	if !args.RunState.IsFinished() && args.HasOpenChanges {
		// Open changes in the middle of an unfinished command will be undone as well.
		// To achieve this, we commit them here so that they are gone when the branch is reset to the original SHA.
		result.Add(&opcodes.CommitOpenChanges{})
	}
	if endBranchesSnapshot, hasEndBranchesSnapshot := args.RunState.EndBranchesSnapshot.Get(); hasEndBranchesSnapshot {
		result.AddProgram(undobranches.DetermineUndoBranchesProgram(args.RunState.BeginBranchesSnapshot, endBranchesSnapshot, args.RunState.UndoablePerennialCommits, args.Config))
	}
	if endConfigSnapshot, hasEndConfigSnapshot := args.RunState.EndConfigSnapshot.Get(); hasEndConfigSnapshot {
		result.AddProgram(undoconfig.DetermineUndoConfigProgram(args.RunState.BeginConfigSnapshot, endConfigSnapshot))
	}
	if endStashSize, hasEndStashsize := args.RunState.EndStashSize.Get(); hasEndStashsize {
		result.AddProgram(undostash.DetermineUndoStashProgram(args.RunState.BeginStashSize, endStashSize))
	}
	result.AddProgram(args.RunState.FinalUndoProgram)
	initialBranch, hasInitialBranch := args.RunState.BeginBranchesSnapshot.Active.Get()
	if hasInitialBranch {
		result.Add(&opcodes.Checkout{Branch: initialBranch})
	}
	cmdhelpers.Wrap(&result, cmdhelpers.WrapOptions{
		DryRun:                   args.RunState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         args.RunState.IsFinished() && args.HasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{initialBranch},
	})
	return result
}
