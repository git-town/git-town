package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestOrderedSet(t *testing.T) {
	t.Parallel()
	t.Run("Add", func(t *testing.T) {
		t.Parallel()
		t.Run("element doesn't exist in set", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two")
			set = set.Add("three")
			have := set.Elements()
			want := []string{"one", "two", "three"}
			assert.Equal(t, want, have)
		})
		t.Run("element already exists in set", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two")
			set = set.Add("two")
			have := set.Elements()
			want := []string{"one", "two"}
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		set := helpers.NewOrderedSet("one", "two")
		assert.True(t, set.Contains("one"))
		assert.True(t, set.Contains("two"))
		assert.False(t, set.Contains("zonk"))
	})

	t.Run("Elements", func(t *testing.T) {
		set := helpers.NewOrderedSet("one", "two")
		have := set.Elements()
		want := []string{"one", "two"}
		assert.Equal(t, want, have)
	})

	t.Run("Join", func(t *testing.T) {
		set := helpers.NewOrderedSet("one", "two", "three")
		have := set.Join(", ")
		want := "one, two, three"
		assert.Equal(t, want, have)

	})
}
