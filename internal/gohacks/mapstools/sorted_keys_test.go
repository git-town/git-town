package mapstools_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/shoenig/test/must"
)

func TestSortedKeys(t *testing.T) {
	t.Parallel()

	t.Run("normal map", func(t *testing.T) {
		t.Parallel()
		give := map[string]string{
			"one":   "1",
			"two":   "2",
			"three": "3",
		}
		keys := []string{}
		for key := range mapstools.SortedKeys(give) {
			keys = append(keys, key)
		}
		must.Eq(t, []string{"one", "three", "two"}, keys)
	})
}
