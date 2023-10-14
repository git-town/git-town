package steps_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/vm/steps"
	"github.com/shoenig/test/must"
)

func TestIsCheckoutStep(t *testing.T) {
	t.Parallel()

	t.Run("given a step.Checkout", func(t *testing.T) {
		t.Parallel()
		give := &step.Checkout{Branch: domain.NewLocalBranchName("branch")}
		must.True(t, steps.IsCheckoutStep(give))
	})

	t.Run("given a step.CheckoutIfExists", func(t *testing.T) {
		t.Parallel()
		give := &step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")}
		must.True(t, steps.IsCheckoutStep(give))
	})

	t.Run("given another step", func(t *testing.T) {
		t.Parallel()
		give := &step.AbortMerge{}
		must.False(t, steps.IsCheckoutStep(give))
	})
}
