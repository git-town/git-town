package prelude_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/gohacks"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestMutable(t *testing.T) {
	t.Parallel()

	t.Run("remains mutable when called by value", func(t *testing.T) {
		t.Parallel()
		counter := gohacks.Counter(0)
		mutable := NewMutable(&counter)
		modify(mutable)
		must.EqOp(t, 1, mutable.Immutable())
	})
}

func modify(byValue Mutable[gohacks.Counter]) {
	byValue.Value.Inc()
}
