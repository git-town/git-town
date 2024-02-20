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

func CreateUndoProgram(args CreateUndoProgramArgs) (program.Program, error) {
	undoProgram := args.RunState.AbortProgram
	undoConfigProgram, err := undoconfig.DetermineUndoConfigProgram(args.InitialConfigSnapshot, args.FinalConfigSnapshot)
	if err != nil {
		return program.Program{}, err
	}
	undoProgram.AddProgram(undoConfigProgram)
	undoBranchesProgram, err := undobranches.DetermineUndoBranchesProgram(args.InitialBranchesSnapshot, args.FinalBranchesSnapshot, args.UndoablePerennialCommits, &args.Run.FullConfig)
	if err != nil {
		return program.Program{}, err
	}
	undoProgram.AddProgram(undoBranchesProgram)
	undoStashProgram, err := undostash.DetermineUndoStashProgram(args.InitialStashSize, &args.Run.Backend)
	if err != nil {
		return program.Program{}, err
	}
	undoProgram.AddProgram(undoStashProgram)
	return undoProgram, nil
}

type CreateUndoProgramArgs struct {
	DryRun                   bool
	FinalBranchesSnapshot    gitdomain.BranchesSnapshot
	FinalConfigSnapshot      undoconfig.ConfigSnapshot
	InitialBranchesSnapshot  gitdomain.BranchesSnapshot
	InitialConfigSnapshot    undoconfig.ConfigSnapshot
	InitialStashSize         gitdomain.StashSize
	NoPushHook               configdomain.NoPushHook
	Run                      *git.ProdRunner
	RunState                 runstate.RunState
	UndoablePerennialCommits []gitdomain.SHA
}
