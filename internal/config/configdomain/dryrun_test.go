package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

func TestDryRun(t *testing.T) {
	t.Parallel()

	t.Run("IsFalse", func(t *testing.T) {
		t.Parallel()
		dryRun := configdomain.DryRun(false)
		if dryRun {
			t.Fail()
		}
	})

	t.Run("IsTrue", func(t *testing.T) {
		t.Parallel()
		dryRun := configdomain.DryRun(true)
		if !dryRun {
			t.Fail()
		}
	})
}
