package undobranches

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

func DetermineUndoBranchesProgram(initialBranchesSnapshot gitdomain.BranchesStatus, undoablePerennialCommits []gitdomain.SHA, noPushHook configdomain.NoPushHook, runner *git.ProdRunner) (program.Program, error) {
	finalBranchesSnapshot, err := runner.Backend.BranchesSnapshot()
	if err != nil {
		return program.Program{}, err
	}
	branchSpans := NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		Lineage:                  runner.GitTown.Lineage(runner.GitTown.RemoveLocalConfigValue),
		BranchTypes:              runner.GitTown.BranchTypes(),
		InitialBranch:            initialBranchesSnapshot.Active,
		FinalBranch:              finalBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
		NoPushHook:               noPushHook,
	}), nil
}
