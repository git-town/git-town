package mapstools_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/gohacks/mapstools"
	"github.com/shoenig/test/must"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	t.Run("both have values", func(t *testing.T) {
		t.Parallel()
		map1 := map[string]int{"one": 1, "two": 2}
		map2 := map[string]int{"three": 3, "four": 4}
		have := mapstools.Merge(map1, map2)
		want := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
		must.Eq(t, want, have)
	})
}
