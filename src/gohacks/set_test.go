package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestSet(t *testing.T) {
	t.Parallel()

	t.Run("Add, Values", func(t *testing.T) {
		t.Parallel()
		set := gohacks.NewSet[int]()
		must.False(t, set.Contains(1))
		set.Add(1)
		must.Eq(t, []int{1}, set.Values())
		set.Add(1)
		must.Eq(t, []int{1}, set.Values())
	})

	t.Run("Contains", func(t *testing.T) {
		set := gohacks.NewSet(1, 2)
		must.True(t, set.Contains(1))
		must.True(t, set.Contains(2))
		must.False(t, set.Contains(3))
	})
}
