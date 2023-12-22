package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

// Tests for SingleCacheDiff are in src/undo/config_test.go.

func TestSingleCache(t *testing.T) {
	t.Parallel()

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		alpha := configdomain.NewKey("alpha")
		beta := configdomain.NewKey("beta")
		original := configdomain.SingleCache{
			alpha: "A",
			beta:  "B",
		}
		cloned := original.Clone()
		cloned[alpha] = "new A"
		cloned[beta] = "new B"
		must.EqOp(t, "A", original[alpha])
		must.EqOp(t, "B", original[beta])
	})

	t.Run("KeysMatching", func(t *testing.T) {
		t.Parallel()
		cache := configdomain.SingleCache{
			configdomain.NewKey("key1"):  "A",
			configdomain.NewKey("key2"):  "B",
			configdomain.NewKey("other"): "other",
		}
		have := cache.KeysMatching("key")
		want := []configdomain.Key{
			configdomain.NewKey("key1"),
			configdomain.NewKey("key2"),
		}
		must.Eq(t, want, have)
	})
}
