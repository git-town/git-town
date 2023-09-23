package undo

import (
	"errors"
	"fmt"
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
	undoBranchesSteps, err := determineUndoBranchesSteps(args.InitialBranchesSnapshot, args.UndoablePerennialCommits, args.Run)
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
	InitialBranchesSnapshot  BranchesSnapshot
	InitialConfigSnapshot    ConfigSnapshot
	InitialStashSnapshot     StashSnapshot
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
	configDiff := initialConfigSnapshot.Diff(finalConfigSnapshot)
	return configDiff.UndoSteps(), nil
}

func determineUndoBranchesSteps(initialBranchesSnapshot BranchesSnapshot, undoablePerennialCommits []domain.SHA, runner *git.ProdRunner) (runstate.StepList, error) {
	allBranches, active, err := runner.Backend.BranchInfos()
	if err != nil {
		return runstate.StepList{}, err
	}
	finalBranchesSnapshot := BranchesSnapshot{
		Branches: allBranches,
		Active:   active,
	}
	branchSpans := initialBranchesSnapshot.Span(finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	fmt.Println("22222222222222222222222222222")
	fmt.Println(branchChanges)
	return branchChanges.UndoSteps(StepsArgs{
		Lineage:                  runner.Config.Lineage(),
		BranchTypes:              runner.Config.BranchTypes(),
		InitialBranch:            initialBranchesSnapshot.Active,
		FinalBranch:              finalBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
	}), nil
}

func determineUndoStashSteps(initialStashSnapshot StashSnapshot, backend *git.BackendCommands) (runstate.StepList, error) {
	stashSize, err := backend.StashSize()
	if err != nil {
		return runstate.StepList{}, err
	}
	finalStashSnapshot := StashSnapshot{
		Amount: stashSize,
	}
	stashDiff := initialStashSnapshot.Diff(finalStashSnapshot)
	return stashDiff.Steps(), nil
}
