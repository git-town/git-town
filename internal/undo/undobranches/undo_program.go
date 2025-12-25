package undobranches

import (
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, validatedConfig config.ValidatedConfig, touchedBranches []gitdomain.BranchName, undoAPIProgram program.Program, finalMessages stringslice.Collector) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnly(touchedBranches)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrZero(),
		BranchInfos:              beginBranchesSnapshot.Branches,
		Config:                   validatedConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrZero(),
		FinalMessages:            finalMessages,
		UndoAPIProgram:           undoAPIProgram,
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
