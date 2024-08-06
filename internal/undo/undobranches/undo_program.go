package undobranches

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, fullConfig configdomain.ValidatedConfig, touchedBranches []gitdomain.BranchName) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnly(touchedBranches)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrDefault(),
		Config:                   fullConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrDefault(),
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
