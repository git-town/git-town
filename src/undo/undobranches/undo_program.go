package undobranches

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, fullConfig *configdomain.FullConfig) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active,
		Config:                   fullConfig,
		EndBranch:                endBranchesSnapshot.Active,
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
