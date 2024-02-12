package undobranches

import (
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

func DetermineUndoBranchesProgram(initialBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, runner *git.ProdRunner) (program.Program, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	branchSpans := NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		Config:                   &runner.FullConfig,
		FinalBranch:              finalBranchesSnapshot.Active,
		InitialBranch:            initialBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
	}), nil
}
