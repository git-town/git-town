package undobranches

import (
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, validatedConfig config.ValidatedConfig, touchedBranches []gitdomain.BranchName, undoAPIProgram program.Program) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnly(touchedBranches)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrDefault(),
		Config:                   validatedConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrDefault(),
		UndoAPIProgram:           undoAPIProgram,
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
