package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/stringslice"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Parallel()
	give := []string{"one", "two"}
	assert.True(t, stringslice.Contains(give, "one"))
	assert.True(t, stringslice.Contains(give, "two"))
	assert.False(t, stringslice.Contains(give, "three"))
}

func TestHoist(t *testing.T) {
	t.Parallel()

	t.Run("already hoisted", func(t *testing.T) {
		t.Parallel()
		give := []string{"main", "one", "two"}
		want := []string{"main", "one", "two"}
		have := stringslice.Hoist(give, "main")
		assert.Equal(t, want, have)
	})

	t.Run("contains the element to hoist", func(t *testing.T) {
		t.Parallel()
		give := []string{"alpha", "main", "omega"}
		want := []string{"main", "alpha", "omega"}
		have := stringslice.Hoist(give, "main")
		assert.Equal(t, want, have)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		give := []string{}
		want := []string{}
		have := stringslice.Hoist(give, "main")
		assert.Equal(t, want, have)
	})
}

func TestLongest(t *testing.T) {
	t.Parallel()
	give := []string{"1", "333", "22"}
	want := 3
	have := stringslice.Longest(give)
	assert.Equal(t, want, have)
}

func TestRemove(t *testing.T) {
	t.Parallel()
	give := []string{"one", "two", "three"}
	want := []string{"one", "three"}
	have := stringslice.Remove(give, "two")
	assert.Equal(t, want, have)
}
