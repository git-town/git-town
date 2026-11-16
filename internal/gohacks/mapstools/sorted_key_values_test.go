package mapstools_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/mapstools"
	"github.com/shoenig/test/must"
)

func TestSortedKeyValues(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		data := map[string]int{
			"four":  4,
			"one":   1,
			"three": 3,
			"two":   2,
		}
		keys := []string{}
		values := []int{}
		for key, value := range mapstools.SortedKeyValues(data) {
			keys = append(keys, key)
			values = append(values, value)
		}
		must.Eq(t, []string{"four", "one", "three", "two"}, keys)
		must.Eq(t, []int{4, 1, 3, 2}, values)
	})

	t.Run("zero content", func(t *testing.T) {
		t.Parallel()
		data := map[string]int{}
		keys := []string{}
		for key := range mapstools.SortedKeyValues(data) {
			keys = append(keys, key)
		}
		must.Len(t, 0, keys)
	})
}
