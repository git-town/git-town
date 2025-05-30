package opcodes_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsCheckout(t *testing.T) {
	t.Parallel()
	branch := gitdomain.NewLocalBranchName("foo")
	tests := map[shared.Opcode]bool{
		&opcodes.CheckoutIfNeeded{Branch: branch}: true,  // Checkout is (obviously) a checkout opcode
		&opcodes.CheckoutIfExists{Branch: branch}: true,  // CheckoutIfExists is also a checkout opcode
		&opcodes.MergeAbort{}:                     false, // any other opcode doesn't match
	}
	for give, want := range tests {
		have := opcodes.IsCheckoutOpcode(give)
		must.Eq(t, want, have)
	}
}

func TestIsEndOfBranchProgramOpcode(t *testing.T) {
	t.Parallel()
	tests := map[shared.Opcode]bool{
		&opcodes.ProgramEndOfBranch{}: true,
		&opcodes.MergeAbort{}:         false,
	}
	for give, want := range tests {
		have := opcodes.IsEndOfBranchProgramOpcode(give)
		must.Eq(t, want, have)
	}
}
