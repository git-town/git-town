package undo

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/undo/undostash"
	"github.com/git-town/git-town/v12/src/vm/program"
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
	undoStashProgram, err := undostash.DetermineUndoStashProgram(args.InitialStashSize, &args.Run.Backend)
	if err != nil {
		return program.Program{}, err
	}
	undoConfigProgram.AddProgram(undoBranchesProgram)
	undoConfigProgram.AddProgram(undoStashProgram)
	return undoConfigProgram, nil
}

type CreateUndoProgramArgs struct {
	DryRun                   bool
	InitialBranchesSnapshot  gitdomain.BranchesSnapshot
	InitialConfigSnapshot    undoconfig.ConfigSnapshot
	InitialStashSize         gitdomain.StashSize
	NoPushHook               configdomain.NoPushHook
	Run                      *git.ProdRunner
	UndoablePerennialCommits []gitdomain.SHA
}
