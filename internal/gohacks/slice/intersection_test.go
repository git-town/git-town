package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFindAllMissing(t *testing.T) {
	t.Parallel()

	t.Run("aliased slice type", func(t *testing.T) {
		t.Parallel()
		existing := gitdomain.SHAs{"111111", "222222"}
		additional := gitdomain.SHAs{"222222", "333333", "444444"}
		have := slice.Intersection(existing, additional...)
		want := gitdomain.SHAs{"333333", "444444"}
		must.Eq(t, want, have)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		list := make([]string, 0)
		have := slice.Intersection(list, "one", "two", "three")
		want := []string{"one", "two", "three"}
		must.Eq(t, want, have)
	})

	t.Run("zero slice", func(t *testing.T) {
		t.Parallel()
		var list []string
		have := slice.Intersection(list, "one", "two", "three")
		want := []string{"one", "two", "three"}
		must.Eq(t, want, have)
	})
}
