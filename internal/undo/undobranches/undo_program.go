package undobranches

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/config"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v19/internal/vm/program"
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, validatedConfig config.ValidatedConfig, touchedBranches []gitdomain.BranchName, undoAPIProgram program.Program, finalMessages stringslice.Collector) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnly(touchedBranches)
	fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB BRANCH SPANS TOUCHED ONLY", branchSpans)
	branchChanges := branchSpans.Changes()
	fmt.Println("CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC CHANGES", branchChanges)
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrDefault(),
		Config:                   validatedConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrDefault(),
		FinalMessages:            finalMessages,
		UndoAPIProgram:           undoAPIProgram,
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
