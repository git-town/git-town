package undo

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
)

func CreateUndoList(args CreateUndoListArgs) (runstate.StepList, error) {
	// TODO: use a StepListBuilder here instead of creating so many separate StepLists and cut down on the redundant error checking.
	undoConfigSteps, err := determineUndoConfigSteps(args.InitialConfigSnapshot, &args.Run.Backend)
	if err != nil {
		return runstate.StepList{}, err
	}
	undoBranchesSteps, err := determineUndoBranchesSteps(args.InitialBranchesSnapshot, args.UndoablePerennialCommits, args.NoPushHook, args.Run)
	if err != nil {
		return runstate.StepList{}, err
	}
	undoStashSteps, err := determineUndoStashSteps(args.InitialStashSnapshot, &args.Run.Backend)
	if err != nil {
		return runstate.StepList{}, err
	}
	undoConfigSteps.AppendList(undoBranchesSteps)
	undoConfigSteps.AppendList(undoStashSteps)
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

func determineUndoConfigSteps(initialConfigSnapshot ConfigSnapshot, backend *git.BackendCommands) (runstate.StepList, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return runstate.StepList{}, errors.New(messages.DirCurrentProblem)
	}
	finalConfigSnapshot := ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(backend),
	}
	configDiff := NewConfigDiffs(initialConfigSnapshot, finalConfigSnapshot)
	return configDiff.UndoSteps(), nil
}

func determineUndoBranchesSteps(initialBranchesSnapshot domain.BranchesSnapshot, undoablePerennialCommits []domain.SHA, noPushHook bool, runner *git.ProdRunner) (runstate.StepList, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return runstate.StepList{}, err
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

func determineUndoStashSteps(initialStashSnapshot domain.StashSnapshot, backend *git.BackendCommands) (runstate.StepList, error) {
	finalStashSnapshot, err := backend.StashSnapshot()
	if err != nil {
		return runstate.StepList{}, err
	}
	stashDiff := NewStashDiff(initialStashSnapshot, finalStashSnapshot)
	return stashDiff.Steps(), nil
}
