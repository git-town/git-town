package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestConnect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give []string
		want string
	}{
		{
			give: []string{},
			want: "",
		},
		{
			give: []string{"one"},
			want: `"one"`,
		},
		{
			give: []string{"one", "two"},
			want: `"one" and "two"`,
		},
		{
			give: []string{"one", "two", "three"},
			want: `"one", "two", and "three"`,
		},
		{
			give: []string{"one", "two", "three", "four"},
			want: `"one", "two", "three", and "four"`,
		},
	}
	for _, test := range tests {
		have := stringslice.Connect(test.give)
		must.EqOp(t, test.want, have)
	}
}
