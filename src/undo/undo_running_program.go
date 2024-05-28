package undo

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/undo/undostash"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
)

// create the program to undo a currently running Git Town command
func CreateUndoForRunningProgram(args CreateUndoProgramArgs) (program.Program, error) {
	result := program.Program{}
	result.AddProgram(args.RunState.AbortProgram)
	if endConfigSnapshot, hasEndConfigSnapshot := args.RunState.EndConfigSnapshot.Get(); hasEndConfigSnapshot {
		result.AddProgram(undoconfig.DetermineUndoConfigProgram(args.RunState.BeginConfigSnapshot, endConfigSnapshot))
	}
	if endBranchesSnapshot, hasEndBranchesSnapshot := args.RunState.EndBranchesSnapshot.Get(); hasEndBranchesSnapshot {
		result.AddProgram(undobranches.DetermineUndoBranchesProgram(args.RunState.BeginBranchesSnapshot, endBranchesSnapshot, args.RunState.UndoablePerennialCommits, args.Config))
	}
	finalStashSize, err := args.Backend.StashSize()
	if err != nil {
		return program.Program{}, err
	}
	result.AddProgram(undostash.DetermineUndoStashProgram(args.RunState.BeginStashSize, finalStashSize))
	return result, nil
}

type CreateUndoProgramArgs struct {
	Backend        git.BackendCommands
	Config         configdomain.ValidatedConfig
	DryRun         bool
	HasOpenChanges bool
	NoPushHook     configdomain.NoPushHook
	RunState       runstate.RunState
}
