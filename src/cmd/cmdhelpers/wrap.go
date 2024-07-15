package cmdhelpers

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// Wrap makes the given program perform housekeeping before and after it executes.
func Wrap(program Mutable[program.Program], options WrapOptions) {
	if program.Value.IsEmpty() {
		return
	}
	if !options.DryRun {
		program.Value.Add(&opcodes.PreserveCheckoutHistory{
			PreviousBranchCandidates: options.PreviousBranchCandidates,
		})
	}
	if options.StashOpenChanges {
		program.Value.Prepend(&opcodes.StashOpenChanges{})
		program.Value.Add(&opcodes.RestoreOpenChanges{})
	}
}

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	DryRun                   configdomain.DryRun
	PreviousBranchCandidates gitdomain.LocalBranchNames
	RunInGitRoot             bool
	StashOpenChanges         bool
}
