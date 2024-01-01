package cmdhelpers

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// Wrap makes the given program perform housekeeping before and after it executes.
func Wrap(program *program.Program, options WrapOptions) {
	if program.IsEmpty() {
		return
	}
	if !options.DryRun {
		program.Add(&opcode.PreserveCheckoutHistory{
			PreviousBranchCandidates: options.PreviousBranchCandidates,
		})
	}
	if options.StashOpenChanges {
		program.Prepend(&opcode.StashOpenChanges{})
		program.Add(&opcode.RestoreOpenChanges{})
	}
}

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	DryRun                   bool
	RunInGitRoot             bool
	StashOpenChanges         bool
	PreviousBranchCandidates gitdomain.LocalBranchNames
}
