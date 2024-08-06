package undo

import (
<<<<<<< HEAD:src/undo/undo_running_program.go
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/undo/undostash"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
=======
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/git"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/undo/undobranches"
	"github.com/git-town/git-town/v14/internal/undo/undoconfig"
	"github.com/git-town/git-town/v14/internal/undo/undostash"
	"github.com/git-town/git-town/v14/internal/vm/program"
	"github.com/git-town/git-town/v14/internal/vm/runstate"
>>>>>>> main:internal/undo/undo_running_program.go
)

// create the program to undo a currently running Git Town command
func CreateUndoForRunningProgram(args CreateUndoProgramArgs) (program.Program, error) {
	result := program.Program{}
	result.AddProgram(args.RunState.AbortProgram)
	if endConfigSnapshot, hasEndConfigSnapshot := args.RunState.EndConfigSnapshot.Get(); hasEndConfigSnapshot {
		result.AddProgram(undoconfig.DetermineUndoConfigProgram(args.RunState.BeginConfigSnapshot, endConfigSnapshot))
	}
	if endBranchesSnapshot, hasEndBranchesSnapshot := args.RunState.EndBranchesSnapshot.Get(); hasEndBranchesSnapshot {
		result.AddProgram(undobranches.DetermineUndoBranchesProgram(args.RunState.BeginBranchesSnapshot, endBranchesSnapshot, args.RunState.UndoablePerennialCommits, args.Config, args.RunState.TouchedBranches, args.Inputs))
	}
	finalStashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return program.Program{}, err
	}
	result.AddProgram(undostash.DetermineUndoStashProgram(args.RunState.BeginStashSize, finalStashSize))
	return result, nil
}

type CreateUndoProgramArgs struct {
	Backend        gitdomain.RunnerQuerier
	Config         configdomain.ValidatedConfig
	DryRun         configdomain.DryRun
	Git            git.Commands
	HasOpenChanges bool
	Inputs         components.TestInputs
	NoPushHook     configdomain.NoPushHook
	RunState       runstate.RunState
}
