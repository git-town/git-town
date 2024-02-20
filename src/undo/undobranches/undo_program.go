package undobranches

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, fullConfig *configdomain.FullConfig) (program.Program, error) {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		Config:                   fullConfig,
		FinalBranch:              endBranchesSnapshot.Active,
		InitialBranch:            beginBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
	}), nil
}
