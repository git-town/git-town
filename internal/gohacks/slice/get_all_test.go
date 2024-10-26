package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestGetAll(t *testing.T) {
	t.Parallel()

	t.Run("all are Some values", func(t *testing.T) {
		t.Parallel()
		give := []Option[int]{Some(1), Some(2), Some(3)}
		have := slice.GetAll(give)
		want := []int{1, 2, 3}
		must.Eq(t, want, have)
	})

	t.Run("all are None values", func(t *testing.T) {
		t.Parallel()
		give := []Option[int]{None[int](), None[int]()}
		have := slice.GetAll(give)
		must.Len(t, 0, have)
	})

	t.Run("mixed values", func(t *testing.T) {
		t.Parallel()
		give := []Option[int]{Some(1), None[int](), Some(3)}
		have := slice.GetAll(give)
		want := []int{1, 3}
		must.Eq(t, want, have)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		give := []Option[int]{}
		have := slice.GetAll(give)
		must.Len(t, 0, have)
	})
}
