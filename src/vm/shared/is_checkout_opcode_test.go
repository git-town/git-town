package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsCheckout(t *testing.T) {
	t.Parallel()

	t.Run("given an opcode.Checkout", func(t *testing.T) {
		t.Parallel()
		give := &opcode.Checkout{Branch: gitdomain.NewLocalBranchName("branch")}
		must.True(t, shared.IsCheckoutOpcode(give))
	})

	t.Run("given an opcode.CheckoutIfExists", func(t *testing.T) {
		t.Parallel()
		give := &opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")}
		must.True(t, shared.IsCheckoutOpcode(give))
	})

	t.Run("given another opcode", func(t *testing.T) {
		t.Parallel()
		give := &opcode.AbortMerge{}
		must.False(t, shared.IsCheckoutOpcode(give))
	})
}
