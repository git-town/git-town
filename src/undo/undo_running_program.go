package undo

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/undo/undostash"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
)

// create the program to undo a currently running Git Town command
func CreateUndoForRunningProgram(args CreateUndoProgramArgs) (program.Program, error) {
	result := program.Program{}
	result.AddProgram(args.RunState.AbortProgram)
	result.AddProgram(undoconfig.DetermineUndoConfigProgram(args.BeginConfigSnapshot, args.EndConfigSnapshot))
	result.AddProgram(undobranches.DetermineUndoBranchesProgram(args.BeginBranchesSnapshot, args.EndBranchesSnapshot, args.RunState.UndoablePerennialCommits, &args.Run.FullConfig))
	finalStashSize, err := args.Run.Backend.StashSize()
	if err != nil {
		return program.Program{}, err
	}
	result.AddProgram(undostash.DetermineUndoStashProgram(args.BeginStashSize, finalStashSize))
	return result, nil
}

type CreateUndoProgramArgs struct {
	BeginBranchesSnapshot gitdomain.BranchesSnapshot
	BeginConfigSnapshot   undoconfig.ConfigSnapshot
	BeginStashSize        gitdomain.StashSize
	DryRun                bool
	EndBranchesSnapshot   gitdomain.BranchesSnapshot
	EndConfigSnapshot     undoconfig.ConfigSnapshot
	HasOpenChanges        bool
	NoPushHook            configdomain.NoPushHook
	Run                   *git.ProdRunner
	RunState              runstate.RunState
}
