package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFirstElementOr(t *testing.T) {
	t.Parallel()

	t.Run("list contains an element", func(t *testing.T) {
		t.Parallel()
		list := []string{"one"}
		have := slice.FirstElementOr(list, "other")
		want := "one"
		must.EqOp(t, want, have)
	})

	t.Run("list is empty", func(t *testing.T) {
		t.Parallel()
		list := []string{}
		have := slice.FirstElementOr(list, "other")
		want := "other"
		must.EqOp(t, want, have)
	})
}
