package cmdhelpers

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Wrap makes the given program perform housekeeping before and after it executes.
func Wrap(program Mutable[program.Program], options WrapOptions) {
	if program.Value.IsEmpty() {
		return
	}
	if !options.DryRun {
		program.Value.Add(&opcodes.CheckoutHistoryPreserve{
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
	PreviousBranchCandidates []Option[gitdomain.LocalBranchName]
	RunInGitRoot             bool
	StashOpenChanges         bool
}
