package undo

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
)

func CreateUndoList(args CreateUndoListArgs) (runstate.StepList, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return runstate.StepList{}, errors.New(messages.DirCurrentProblem)
	}
	finalConfigSnapshot := ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(args.Run.Backend),
	}
	configDiff := finalConfigSnapshot.Diff(args.InitialConfigSnapshot)
	undoConfigSteps := configDiff.UndoSteps()
	allBranches, active, err := args.Run.Backend.BranchInfos()
	if err != nil {
		return runstate.StepList{}, err
	}
	finalBranchesSnapshot := BranchesSnapshot{
		Branches: allBranches,
		Active:   active,
	}
	bba := args.InitialBranchesSnapshot.Changes(finalBranchesSnapshot)
	branchesDiff := bba.Diff()
	undoBranchesSteps := branchesDiff.Steps(args.Run.Config.Lineage(), args.Run.Config.BranchTypes(), args.InitialBranchesSnapshot.Active, finalBranchesSnapshot.Active)
	undoConfigSteps.AppendList(undoBranchesSteps)
	return undoConfigSteps, nil
}

type CreateUndoListArgs struct {
	Run                     *git.ProdRunner
	InitialBranchesSnapshot BranchesSnapshot
	InitialConfigSnapshot   ConfigSnapshot
}
