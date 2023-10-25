package undo

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/opcode"
)

func CreateUndoProgram(args CreateUndoProgramArgs) (opcode.Program, error) {
	undoConfigProgram, err := determineUndoConfigProgram(args.InitialConfigSnapshot, &args.Run.Backend)
	if err != nil {
		return opcode.Program{}, err
	}
	undoBranchesProgram, err := determineUndoBranchesProgram(args.InitialBranchesSnapshot, args.UndoablePerennialCommits, args.NoPushHook, args.Run)
	if err != nil {
		return opcode.Program{}, err
	}
	undoStashProgram, err := determineUndoStashProgram(args.InitialStashSnapshot, &args.Run.Backend)
	if err != nil {
		return opcode.Program{}, err
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

func determineUndoBranchesProgram(initialBranchesSnapshot domain.BranchesSnapshot, undoablePerennialCommits []domain.SHA, noPushHook bool, runner *git.ProdRunner) (opcode.Program, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return opcode.Program{}, err
	}
	branchSpans := NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		Lineage:                  runner.Config.Lineage(),
		BranchTypes:              runner.Config.BranchTypes(),
		InitialBranch:            initialBranchesSnapshot.Active,
		FinalBranch:              finalBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
		NoPushHook:               noPushHook,
	}), nil
}

func determineUndoConfigProgram(initialConfigSnapshot ConfigSnapshot, backend *git.BackendCommands) (opcode.Program, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return opcode.Program{}, errors.New(messages.DirCurrentProblem)
	}
	finalConfigSnapshot := ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(backend),
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoProgram(), nil
}

func determineUndoStashProgram(initialStashSnapshot domain.StashSnapshot, backend *git.BackendCommands) (opcode.Program, error) {
	finalStashSnapshot, err := backend.StashSnapshot()
	if err != nil {
		return opcode.Program{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Program(), nil
}
