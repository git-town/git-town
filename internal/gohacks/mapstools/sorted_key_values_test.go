package mapstools_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/gohacks/mapstools"
	"github.com/shoenig/test/must"
)

func TestSortedKeyValues(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		data := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  4,
		}
		have := mapstools.SortedKeyValues(data)
		must.Len(t, 4, have)
		must.EqOp(t, "four", have[0].Key)
		must.EqOp(t, 4, have[0].Value)
		must.EqOp(t, "one", have[1].Key)
		must.EqOp(t, 1, have[1].Value)
		must.EqOp(t, "three", have[2].Key)
		must.EqOp(t, 3, have[2].Value)
		must.EqOp(t, "two", have[3].Key)
		must.EqOp(t, 2, have[3].Value)
	})
}
