package mapstools_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/mapstools"
	"github.com/shoenig/test/must"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	t.Run("both are nil", func(t *testing.T) {
		t.Parallel()
		var map1 map[string]int
		have := mapstools.Merge(map1, nil)
		want := map[string]int{}
		must.Eq(t, want, have)
	})

	t.Run("both have values", func(t *testing.T) {
		t.Parallel()
		map1 := map[string]int{"one": 1, "two": 2}
		map2 := map[string]int{"four": 4, "three": 3}
		have := mapstools.Merge(map1, map2)
		want := map[string]int{"four": 4, "one": 1, "three": 3, "two": 2}
		must.Eq(t, want, have)
	})

	t.Run("first is nil", func(t *testing.T) {
		t.Parallel()
		map1 := map[string]int{"one": 1, "two": 2}
		have := mapstools.Merge(nil, map1)
		want := map[string]int{"one": 1, "two": 2}
		must.Eq(t, want, have)
	})

	t.Run("second is nil", func(t *testing.T) {
		t.Parallel()
		map1 := map[string]int{"one": 1, "two": 2}
		have := mapstools.Merge(map1, nil)
		want := map[string]int{"one": 1, "two": 2}
		must.Eq(t, want, have)
	})
}
