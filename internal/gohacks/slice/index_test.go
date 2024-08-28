package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/gohacks/slice"
	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	t.Run("contains element", func(t *testing.T) {
		t.Parallel()
		haystack := []int{1, 2, 3}
		have := slice.Index(haystack, 2)
		want := Some(1)
		must.Eq(t, want, have)
	})
}
