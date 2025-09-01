package undobranches

import (
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, validatedConfig config.ValidatedConfig, touchedBranches []gitdomain.BranchName, undoAPIProgram program.Program, finalMessages stringslice.Collector) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnly(touchedBranches)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrZero(),
		Config:                   validatedConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrZero(),
		FinalMessages:            finalMessages,
		UndoAPIProgram:           undoAPIProgram,
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
