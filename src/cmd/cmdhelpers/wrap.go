package cmdhelpers

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// wrap wraps the given list with opcodes that change the Git root directory or stash away open changes.
// TODO: only wrap if the list actually contains any opcodes.
func Wrap(program *program.Program, options WrapOptions) {
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
	RunInGitRoot             bool
	StashOpenChanges         bool
	PreviousBranchCandidates gitdomain.LocalBranchNames
}
