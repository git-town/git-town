package undo

import (
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/undo/undobranches"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/undo/undostash"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
)

// creates the program for undoing a program that finished
func CreateUndoForFinishedProgram(args CreateUndoProgramArgs) program.Program {
	result := NewMutable(&program.Program{})
	result.Value.AddProgram(args.RunState.AbortProgram)
	if !args.RunState.IsFinished() && args.HasOpenChanges {
		// Open changes in the middle of an unfinished command will be undone as well.
		// To achieve this, we commit them here so that they are gone when the branch is reset to the original SHA.
		result.Value.Add(&opcodes.CommitOpenChanges{
			AddAll:  true,
			Message: "Committing WIP for git town undo",
		})
	}
	if endBranchesSnapshot, hasEndBranchesSnapshot := args.RunState.EndBranchesSnapshot.Get(); hasEndBranchesSnapshot {
		result.Value.AddProgram(undobranches.DetermineUndoBranchesProgram(args.RunState.BeginBranchesSnapshot, endBranchesSnapshot, args.RunState.UndoablePerennialCommits, args.Config, args.RunState.TouchedBranches))
	}
	if endConfigSnapshot, hasEndConfigSnapshot := args.RunState.EndConfigSnapshot.Get(); hasEndConfigSnapshot {
		result.Value.AddProgram(undoconfig.DetermineUndoConfigProgram(args.RunState.BeginConfigSnapshot, endConfigSnapshot))
	}
	if endStashSize, hasEndStashsize := args.RunState.EndStashSize.Get(); hasEndStashsize {
		result.Value.AddProgram(undostash.DetermineUndoStashProgram(args.RunState.BeginStashSize, endStashSize))
	}
	result.Value.AddProgram(args.RunState.FinalUndoProgram)
	previousBranchCandidates := gitdomain.LocalBranchNames{}
	if initialBranch, hasInitialBranch := args.RunState.BeginBranchesSnapshot.Active.Get(); hasInitialBranch {
		result.Value.Add(&opcodes.Checkout{Branch: initialBranch})
		previousBranchCandidates = append(previousBranchCandidates, initialBranch)
	}
	cmdhelpers.Wrap(result, cmdhelpers.WrapOptions{
		DryRun:                   args.RunState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         args.RunState.IsFinished() && args.HasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return result.Get()
}
