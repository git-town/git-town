package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestLongest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give []string
		want int
	}{
		{
			give: []string{"one", "two", "three"},
			want: 5,
		},
		{
			give: []string{""},
			want: 0,
		},
		{
			give: []string{},
			want: 0,
		},
	}
	for _, test := range tests {
		have := stringslice.Longest(test.give)
		must.EqOp(t, test.want, have)
	}
}
