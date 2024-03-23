package undo

import (
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git"
	"github.com/git-town/git-town/v13/src/undo/undobranches"
	"github.com/git-town/git-town/v13/src/undo/undoconfig"
	"github.com/git-town/git-town/v13/src/undo/undostash"
	"github.com/git-town/git-town/v13/src/vm/program"
	"github.com/git-town/git-town/v13/src/vm/runstate"
)

// create the program to undo a currently running Git Town command
func CreateUndoForRunningProgram(args CreateUndoProgramArgs) (program.Program, error) {
	result := program.Program{}
	result.AddProgram(args.RunState.AbortProgram)
	result.AddProgram(undoconfig.DetermineUndoConfigProgram(args.RunState.BeginConfigSnapshot, args.RunState.EndConfigSnapshot))
	result.AddProgram(undobranches.DetermineUndoBranchesProgram(args.RunState.BeginBranchesSnapshot, args.RunState.EndBranchesSnapshot, args.RunState.UndoablePerennialCommits, &args.Run.Config.FullConfig))
	finalStashSize, err := args.Run.Backend.StashSize()
	if err != nil {
		return program.Program{}, err
	}
	result.AddProgram(undostash.DetermineUndoStashProgram(args.RunState.BeginStashSize, finalStashSize))
	return result, nil
}

type CreateUndoProgramArgs struct {
	DryRun         bool
	HasOpenChanges bool
	NoPushHook     configdomain.NoPushHook
	Run            *git.ProdRunner
	RunState       runstate.RunState
}
