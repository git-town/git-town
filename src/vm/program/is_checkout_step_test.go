package program_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestIsCheckoutStep(t *testing.T) {
	t.Parallel()

	t.Run("given a step.Checkout", func(t *testing.T) {
		t.Parallel()
		give := &opcode.Checkout{Branch: domain.NewLocalBranchName("branch")}
		must.True(t, program.IsCheckoutStep(give))
	})

	t.Run("given a step.CheckoutIfExists", func(t *testing.T) {
		t.Parallel()
		give := &opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")}
		must.True(t, program.IsCheckoutStep(give))
	})

	t.Run("given another step", func(t *testing.T) {
		t.Parallel()
		give := &opcode.AbortMerge{}
		must.False(t, program.IsCheckoutStep(give))
	})
}
