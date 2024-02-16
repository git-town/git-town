package cmdhelpers

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// Wrap makes the given program perform housekeeping before and after it executes.
func Wrap(program *program.Program, options WrapOptions) {
	fmt.Println("11111111111111", program)
	if program.IsEmpty() {
		return
	}
	if !options.DryRun {
		program.Add(&opcodes.PreserveCheckoutHistory{
			PreviousBranchCandidates: options.PreviousBranchCandidates,
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
	PreviousBranchCandidates gitdomain.LocalBranchNames
	RunInGitRoot             bool
	StashOpenChanges         bool
}
