package undo

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/steps"
)

func CreateUndoList(args CreateUndoListArgs) (steps.List, error) {
	undoConfigSteps, err := determineUndoConfigSteps(args.InitialConfigSnapshot, &args.Run.Backend)
	if err != nil {
		return steps.List{}, err
	}
	undoBranchesSteps, err := determineUndoBranchesSteps(args.InitialBranchesSnapshot, args.UndoablePerennialCommits, args.NoPushHook, args.Run)
	if err != nil {
		return steps.List{}, err
	}
	undoStashSteps, err := determineUndoStashSteps(args.InitialStashSnapshot, &args.Run.Backend)
	if err != nil {
		return steps.List{}, err
	}
	undoConfigSteps.AddList(undoBranchesSteps)
	undoConfigSteps.AddList(undoStashSteps)
	return undoConfigSteps, nil
}

type CreateUndoListArgs struct {
	Run                      *git.ProdRunner
	InitialBranchesSnapshot  domain.BranchesSnapshot
	InitialConfigSnapshot    ConfigSnapshot
	InitialStashSnapshot     domain.StashSnapshot
	NoPushHook               bool
	UndoablePerennialCommits []domain.SHA
}

func determineUndoConfigSteps(initialConfigSnapshot ConfigSnapshot, backend *git.BackendCommands) (steps.List, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return steps.List{}, errors.New(messages.DirCurrentProblem)
	}
	finalConfigSnapshot := ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(backend),
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoSteps(), nil
}

func determineUndoBranchesSteps(initialBranchesSnapshot domain.BranchesSnapshot, undoablePerennialCommits []domain.SHA, noPushHook bool, runner *git.ProdRunner) (steps.List, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return steps.List{}, err
	}
	branchSpans := NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoSteps(StepsArgs{
		Lineage:                  runner.Config.Lineage(),
		BranchTypes:              runner.Config.BranchTypes(),
		InitialBranch:            initialBranchesSnapshot.Active,
		FinalBranch:              finalBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
		NoPushHook:               noPushHook,
	}), nil
}

func determineUndoStashSteps(initialStashSnapshot domain.StashSnapshot, backend *git.BackendCommands) (steps.List, error) {
	finalStashSnapshot, err := backend.StashSnapshot()
	if err != nil {
		return steps.List{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Steps(), nil
}
