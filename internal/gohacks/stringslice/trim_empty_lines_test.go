package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestTrimEmptyLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give []string
		want []string
	}{
		{
			name: "no empty lines",
			give: []string{"one", "two"},
			want: []string{"one", "two"},
		},
		{
			name: "trailing empty line",
			give: []string{"one", "two", ""},
			want: []string{"one", "two"},
		},
		{
			name: "multiple trailing empty lines",
			give: []string{"one", "two", "", "", ""},
			want: []string{"one", "two"},
		},
		{
			name: "leading empty lines",
			give: []string{"", "", "one", "two"},
			want: []string{"one", "two"},
		},
		{
			name: "empty string",
			give: []string{},
			want: []string{},
		},
		{
			name: "only empty lines",
			give: []string{"", "", ""},
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
