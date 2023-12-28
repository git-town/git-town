package cmdhelpers

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// Wrap makes the given program perform housekeeping before and after it executes.
// TODO: only wrap if the program actually contains any opcodes.
func Wrap(program *program.Program, options WrapOptions) {
	if options.DryRun {
		return
	}
	program.Add(&opcode.PreserveCheckoutHistory{
		PreviousBranchCandidates: options.PreviousBranchCandidates,
	})
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
