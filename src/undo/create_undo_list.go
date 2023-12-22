package undo

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/undo/undobranches"
	"github.com/git-town/git-town/v11/src/undo/undoconfig"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func CreateUndoProgram(args CreateUndoProgramArgs) (program.Program, error) {
	undoConfigProgram, err := determineUndoConfigProgram(args.InitialConfigSnapshot, &args.Run.GitTown.Access)
	if err != nil {
		return program.Program{}, err
	}
	undoBranchesProgram, err := determineUndoBranchesProgram(args.InitialBranchesSnapshot, args.UndoablePerennialCommits, args.NoPushHook, args.Run)
	if err != nil {
		return program.Program{}, err
	}
	undoStashProgram, err := determineUndoStashProgram(args.InitialStashSnapshot, &args.Run.Backend)
	if err != nil {
		return program.Program{}, err
	}
	undoConfigProgram.AddProgram(undoBranchesProgram)
	undoConfigProgram.AddProgram(undoStashProgram)
	return undoConfigProgram, nil
}

type CreateUndoProgramArgs struct {
	Run                      *git.ProdRunner
	InitialBranchesSnapshot  undodomain.BranchesSnapshot
	InitialConfigSnapshot    undodomain.ConfigSnapshot
	InitialStashSnapshot     undodomain.StashSnapshot
	NoPushHook               configdomain.NoPushHook
	UndoablePerennialCommits []gitdomain.SHA
}

func determineUndoBranchesProgram(initialBranchesSnapshot undodomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, noPushHook configdomain.NoPushHook, runner *git.ProdRunner) (program.Program, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	branchSpans := undobranches.NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		Lineage:                  runner.GitTown.Lineage(runner.GitTown.RemoveLocalConfigValue),
		BranchTypes:              runner.GitTown.BranchTypes(),
		InitialBranch:            initialBranchesSnapshot.Active,
		FinalBranch:              finalBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
		NoPushHook:               noPushHook,
	}), nil
}

func determineUndoConfigProgram(initialConfigSnapshot undodomain.ConfigSnapshot, configGit *gitconfig.Access) (program.Program, error) {
	fullCache, err := gitconfig.LoadFullCache(configGit)
	if err != nil {
		return program.Program{}, err
	}
	finalConfigSnapshot := undodomain.ConfigSnapshot{
		GitConfig: fullCache,
	}
	configDiff := undoconfig.NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram(), nil
}

func determineUndoStashProgram(initialStashSnapshot undodomain.StashSnapshot, backend *git.BackendCommands) (program.Program, error) {
	finalStashSnapshot, err := backend.StashSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Program(), nil
}
