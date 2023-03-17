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
		give := []string{"initial", "one", "two"}
		want := []string{"initial", "one", "two"}
		have := stringslice.Hoist(give, "initial")
		assert.Equal(t, want, have)
	})

	t.Run("contains the element to hoist", func(t *testing.T) {
		t.Parallel()
		give := []string{"alpha", "initial", "omega"}
		want := []string{"initial", "alpha", "omega"}
		have := stringslice.Hoist(give, "initial")
		assert.Equal(t, want, have)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		give := []string{}
		want := []string{}
		have := stringslice.Hoist(give, "initial")
		assert.Equal(t, want, have)
	})
}

func TestRemove(t *testing.T) {
	t.Parallel()
	give := []string{"one", "two", "three"}
	have := stringslice.Remove(give, "two")
	want := []string{"one", "three"}
	assert.Equal(t, have, want)
}
