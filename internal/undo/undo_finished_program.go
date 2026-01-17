package undo

import (
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/programs"
	"github.com/git-town/git-town/v22/internal/undo/undobranches"
	"github.com/git-town/git-town/v22/internal/undo/undoconfig"
	"github.com/git-town/git-town/v22/internal/undo/undostash"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// creates the program for undoing a program that finished
func CreateUndoForFinishedProgram(args CreateUndoProgramArgs) (undoProgram program.Program, changedBranches gitdomain.LocalBranchNames) {
	result := NewMutable(&program.Program{})
	result.Value.AddProgram(args.RunState.AbortProgram)
	if !args.RunState.IsFinished() && args.HasOpenChanges {
		// Open changes in the middle of an unfinished command will be undone as well.
		// To achieve this, we commit them here so that they are gone when the branch is reset to the original SHA.
		result.Value.Add(&opcodes.ChangesStage{})
		result.Value.Add(&opcodes.CommitWithMessage{
			AuthorOverride: None[gitdomain.Author](),
			CommitHook:     configdomain.CommitHookEnabled,
			Message:        "Committing open changes to undo them",
		})
	}
	// undo config changes
	if endConfigSnapshot, hasEndConfigSnapshot := args.RunState.EndConfigSnapshot.Get(); hasEndConfigSnapshot {
		result.Value.AddProgram(undoconfig.DetermineUndoConfigProgram(args.RunState.BeginConfigSnapshot, endConfigSnapshot))
	}
	// undo branch changes
	endBranchesSnapshot, hasEndBranchesSnapshot := args.RunState.EndBranchesSnapshot.Get()
	if hasEndBranchesSnapshot {
		undoProgram, changedBranches = undobranches.DetermineUndoBranchesProgram(args.RunState.BeginBranchesSnapshot, endBranchesSnapshot, args.RunState.UndoablePerennialCommits, args.Config, args.RunState.TouchedBranches, args.RunState.UndoAPIProgram, args.FinalMessages)
		result.Value.AddProgram(undoProgram)
	}
	// undo stash changes
	if endStashSize, hasEndStashsize := args.RunState.EndStashSize.Get(); hasEndStashsize {
		result.Value.AddProgram(undostash.DetermineUndoStashProgram(args.RunState.BeginStashSize, endStashSize))
	}
	updateProposalLineage := args.Config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI
	isOnline := args.Config.NormalConfig.Offline.IsOnline()
	if updateProposalLineage && isOnline {
		programs.AddSyncProposalsProgram(programs.AddSyncProposalsProgramArgs{
			ChangedBranches: changedBranches,
			Config:          args.Config,
			Program:         result,
		})
	}
	result.Value.AddProgram(args.RunState.FinalUndoProgram)
	initialBranchOpt := args.RunState.BeginBranchesSnapshot.Active
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{initialBranchOpt}
	if initialBranch, hasInitialBranch := initialBranchOpt.Get(); hasInitialBranch {
		result.Value.Add(&opcodes.CheckoutIfNeeded{Branch: initialBranch})
	}
	cmdhelpers.Wrap(result, cmdhelpers.WrapOptions{
		DryRun:                   args.RunState.DryRun,
		InitialStashSize:         args.RunState.BeginStashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         args.RunState.IsFinished() && args.HasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return result.Immutable(), changedBranches
}
