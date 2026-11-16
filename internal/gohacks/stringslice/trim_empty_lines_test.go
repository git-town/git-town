package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestTrimEmptyLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give []string
		name string
		want []string
	}{
		{
			give: []string{"one", "two"},
			name: "no empty lines",
			want: []string{"one", "two"},
		},
		{
			give: []string{"one", "two", ""},
			name: "trailing empty line",
			want: []string{"one", "two"},
		},
		{
			give: []string{"one", "two", "", "", ""},
			name: "multiple trailing empty lines",
			want: []string{"one", "two"},
		},
		{
			give: []string{"", "", "one", "two"},
			name: "leading empty lines",
			want: []string{"one", "two"},
		},
		{
			give: []string{"", "one", "", "two", ""},
			name: "empty lines before, middle, and after",
			want: []string{"one", "", "two"},
		},
		{
			give: []string{},
			name: "empty string",
			want: []string{},
		},
		{
			give: []string{"", "", ""},
			name: "only empty lines",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have := stringslice.TrimEmptyLines(tt.give)
			must.Eq(t, tt.want, have)
		})
	}
}
