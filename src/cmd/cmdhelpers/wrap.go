package cmdhelpers

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// Wrap makes the given program perform housekeeping before and after it executes.
func Wrap(program *program.Program, options WrapOptions) {
	if program.IsEmpty() {
		return
	}
	if !options.DryRun {
		previousBranchCandidates := options.PreviousBranchCandidates.LocalBranches().WithStatuses(gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusNotInSync)
		program.Add(&opcodes.PreserveCheckoutHistory{
			PreviousBranchCandidates: previousBranchCandidates.Names(),
		})
	}
	if options.StashOpenChanges {
		program.Prepend(&opcodes.StashOpenChanges{})
		program.Add(&opcodes.RestoreOpenChanges{})
	}
}

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	DryRun                   bool
	PreviousBranchCandidates gitdomain.BranchInfos
	RunInGitRoot             bool
	StashOpenChanges         bool
}
