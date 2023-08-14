package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/slice"
	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	t.Parallel()

	t.Run("AppendAllMissing", func(t *testing.T) {
		t.Parallel()
		list := []string{"one", "two", "three"}
		give := []string{"two", "four", "five"}
		want := []string{"one", "two", "three", "four", "five"}
		have := slice.AppendAllMissing(list, give)
		assert.Equal(t, want, have)
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two"}
		assert.True(t, slice.Contains(give, "one"))
		assert.True(t, slice.Contains(give, "two"))
		assert.False(t, slice.Contains(give, "three"))
	})

	t.Run("FirstElementOr", func(t *testing.T) {
		t.Parallel()
		t.Run("list contains an element", func(t *testing.T) {
			t.Parallel()
			list := []string{"one"}
			have := slice.FirstElementOr(list, "other")
			want := "one"
			assert.Equal(t, want, have)
		})
		t.Run("list is empty", func(t *testing.T) {
			t.Parallel()
			list := []string{}
			have := slice.FirstElementOr(list, "other")
			want := "other"
			assert.Equal(t, want, have)
		})
	})

	t.Run("Hoist", func(t *testing.T) {
		t.Parallel()

		t.Run("already hoisted", func(t *testing.T) {
			t.Parallel()
			give := []string{"initial", "one", "two"}
			want := []string{"initial", "one", "two"}
			have := slice.Hoist(give, "initial")
			assert.Equal(t, want, have)
		})

		t.Run("contains the element to hoist", func(t *testing.T) {
			t.Parallel()
			give := []string{"alpha", "initial", "omega"}
			want := []string{"initial", "alpha", "omega"}
			have := slice.Hoist(give, "initial")
			assert.Equal(t, want, have)
		})

		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := []string{}
			want := []string{}
			have := slice.Hoist(give, "initial")
			assert.Equal(t, want, have)
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "three"}
		have := slice.Remove(give, "two")
		want := []string{"one", "three"}
		assert.Equal(t, have, want)
	})
}
