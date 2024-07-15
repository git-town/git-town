package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestDryRun(t *testing.T) {
	t.Parallel()

	t.Run("IsTrue", func(t *testing.T) {
		t.Parallel()
		dryRun := configdomain.DryRun(true)
		must.True(t, dryRun.IsTrue())
		must.False(t, dryRun.IsFalse())
	})

	t.Run("IsFalse", func(t *testing.T) {
		t.Parallel()
		dryRun := configdomain.DryRun(false)
		must.False(t, dryRun.IsTrue())
		must.True(t, dryRun.IsFalse())
	})
}
