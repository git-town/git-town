package undo

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/undo/undobranches"
	"github.com/git-town/git-town/v11/src/undo/undoconfig"
	"github.com/git-town/git-town/v11/src/undo/undostash"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func CreateUndoProgram(args CreateUndoProgramArgs) (program.Program, error) {
	undoConfigProgram, err := undoconfig.DetermineUndoConfigProgram(args.InitialConfigSnapshot, &args.Run.Config.GitConfig)
	if err != nil {
		return program.Program{}, err
	}
	undoBranchesProgram, err := undobranches.DetermineUndoBranchesProgram(args.InitialBranchesSnapshot, args.UndoablePerennialCommits, args.Run)
	if err != nil {
		return program.Program{}, err
	}
	undoStashProgram, err := undostash.DetermineUndoStashProgram(args.InitialStashSnapshot, &args.Run.Backend)
	if err != nil {
		return program.Program{}, err
	}
	undoConfigProgram.AddProgram(undoBranchesProgram)
	undoConfigProgram.AddProgram(undoStashProgram)
	return undoConfigProgram, nil
}

type CreateUndoProgramArgs struct {
	DryRun                   bool
	Run                      *git.ProdRunner
	InitialBranchesSnapshot  gitdomain.BranchesStatus
	InitialConfigSnapshot    undoconfig.ConfigSnapshot
	InitialStashSnapshot     gitdomain.StashSize
	NoPushHook               configdomain.NoPushHook
	UndoablePerennialCommits []gitdomain.SHA
}
