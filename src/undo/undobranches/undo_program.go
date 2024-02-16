package undobranches

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

func DetermineUndoBranchesProgram(initialBranchesSnapshot, finalBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, fullConfig *configdomain.FullConfig) (program.Program, error) {
	branchSpans := NewBranchSpans(initialBranchesSnapshot, finalBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		Config:                   fullConfig,
		FinalBranch:              finalBranchesSnapshot.Active,
		InitialBranch:            initialBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
	}), nil
}
