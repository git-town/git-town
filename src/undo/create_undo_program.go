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
	undoProgram := program.Program{}
	undoProgram.AddProgram(undoconfig.DetermineUndoConfigProgram(args.BeginConfigSnapshot, args.EndConfigSnapshot))
	undoProgram.AddProgram(undobranches.DetermineUndoBranchesProgram(args.BeginBranchesSnapshot, args.EndBranchesSnapshot, args.UndoablePerennialCommits, &args.Run.FullConfig))
	finalStashSize, err := args.Run.Backend.StashSize()
	if err != nil {
		return program.Program{}, err
	}
	undoProgram.AddProgram(undostash.DetermineUndoStashProgram(args.BeginStashSize, finalStashSize))
	return undoProgram, nil
}

type CreateUndoProgramArgs struct {
	BeginBranchesSnapshot    gitdomain.BranchesSnapshot
	BeginConfigSnapshot      undoconfig.ConfigSnapshot
	BeginStashSize           gitdomain.StashSize
	DryRun                   bool
	EndBranchesSnapshot      gitdomain.BranchesSnapshot
	EndConfigSnapshot        undoconfig.ConfigSnapshot
	NoPushHook               configdomain.NoPushHook
	Run                      *git.ProdRunner
	UndoablePerennialCommits []gitdomain.SHA
}
