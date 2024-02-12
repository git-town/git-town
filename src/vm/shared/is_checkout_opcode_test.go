package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcode"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsCheckout(t *testing.T) {
	t.Parallel()
	tests := map[shared.Opcode]bool{
		&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("branch")}:         true,
		&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")}: true,
		&opcode.AbortMerge{}: false,
	}
	for give, want := range tests {
		have := shared.IsCheckoutOpcode(give)
		must.Eq(t, want, have)
	}
}
