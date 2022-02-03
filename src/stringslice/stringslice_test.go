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
