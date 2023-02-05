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

func TestLast(t *testing.T) {
	t.Parallel()

	t.Run("empty list", func(t *testing.T) {
		give := []string{}
		have := stringslice.Last(give)
		want := (*string)(nil)
		assert.Equal(t, want, have)
	})

	t.Run("one element", func(t *testing.T) {
		one := "one"
		give := []string{one}
		have := stringslice.Last(give)
		assert.Equal(t, &one, have)
	})

	t.Run("many elements", func(t *testing.T) {
		one := "one"
		two := "two"
		give := []string{one, two}
		have := stringslice.Last(give)
		assert.Equal(t, &two, have)
	})
}

func TestMainFirst(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give []string
		want []string
	}{
		{give: []string{"main", "one", "two"}, want: []string{"main", "one", "two"}},
		{give: []string{"alpha", "main", "omega"}, want: []string{"main", "alpha", "omega"}},
		{give: []string{"main"}, want: []string{"main"}},
		{give: []string{}, want: []string{}},
	}
	for _, test := range tests {
		have := stringslice.MainFirst(test.give)
		assert.Equal(t, test.want, have)
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	give := []string{"one", "two", "three"}
	have := stringslice.Remove(give, "two")
	want := []string{"one", "three"}
	assert.Equal(t, have, want)
}

func TestRemoveAll(t *testing.T) {
	t.Run("remove no existing elements", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two"}
		have := stringslice.RemoveMany(give, []string{"zonk"})
		assert.Equal(t, give, have)
	})

	t.Run("remove some existing elements", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "three"}
		have := stringslice.RemoveMany(give, []string{"one", "three"})
		assert.Equal(t, []string{"two"}, have)
	})

	t.Run("remove all existing elements", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two"}
		have := stringslice.RemoveMany(give, []string{"two", "one"})
		assert.Equal(t, []string{}, have)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		give := []string{}
		have := stringslice.RemoveMany(give, []string{"zonk"})
		assert.Equal(t, []string{}, have)
	})

	t.Run("empty remove set", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two"}
		have := stringslice.RemoveMany(give, []string{})
		assert.Equal(t, give, have)
	})
}
