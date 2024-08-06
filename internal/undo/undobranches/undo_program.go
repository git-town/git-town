package undobranches

import (
<<<<<<< HEAD:src/undo/undobranches/undo_program.go
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/program"
=======
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/program"
>>>>>>> main:internal/undo/undobranches/undo_program.go
)

func DetermineUndoBranchesProgram(beginBranchesSnapshot, endBranchesSnapshot gitdomain.BranchesSnapshot, undoablePerennialCommits []gitdomain.SHA, fullConfig configdomain.ValidatedConfig, touchedBranches []gitdomain.BranchName, inputs components.TestInputs) program.Program {
	branchSpans := NewBranchSpans(beginBranchesSnapshot, endBranchesSnapshot)
	branchSpans = branchSpans.KeepOnly(touchedBranches)
	branchChanges := branchSpans.Changes()
	return branchChanges.UndoProgram(BranchChangesUndoProgramArgs{
		BeginBranch:              beginBranchesSnapshot.Active.GetOrDefault(),
		Config:                   fullConfig,
		EndBranch:                endBranchesSnapshot.Active.GetOrDefault(),
		Inputs:                   inputs,
		UndoablePerennialCommits: undoablePerennialCommits,
	})
}
