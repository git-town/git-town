package undobranches

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, fullConfig configdomain.ValidatedConfig, touchedBranches []gitdomain.BranchName) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnlyTheseBranches(touchedBranches)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrDefault(),
		Config:                   fullConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrDefault(),
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
