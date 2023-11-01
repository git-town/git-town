package undo

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/program"
)

func CreateUndoProgram(args CreateUndoProgramArgs) (program.Program, error) {
	undoConfigProgram, err := determineUndoConfigProgram(args.InitialConfigSnapshot, &args.Run.Backend)
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
	InitialBranchesSnapshot  domain.BranchesSnapshot
	InitialConfigSnapshot    ConfigSnapshot
	InitialStashSnapshot     domain.StashSnapshot
	NoPushHook               bool
	UndoablePerennialCommits []domain.SHA
}

func determineUndoBranchesProgram(initialBranchesSnapshot domain.BranchesSnapshot, undoablePerennialCommits []domain.SHA, noPushHook bool, runner *git.ProdRunner) (program.Program, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	branchSpans := NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		Lineage:                  runner.Config.Lineage(runner.Config.RemoveLocalConfigValue),
		BranchTypes:              runner.Config.BranchTypes(),
		InitialBranch:            initialBranchesSnapshot.Active,
		FinalBranch:              finalBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
		NoPushHook:               noPushHook,
	}), nil
}

func determineUndoConfigProgram(initialConfigSnapshot ConfigSnapshot, backend *git.BackendCommands) (program.Program, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return program.Program{}, errors.New(messages.DirCurrentProblem)
	}
	finalConfigSnapshot := ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(backend),
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram(), nil
}

func determineUndoStashProgram(initialStashSnapshot domain.StashSnapshot, backend *git.BackendCommands) (program.Program, error) {
	finalStashSnapshot, err := backend.StashSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Program(), nil
}
