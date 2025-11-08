package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestEqualIgnoreWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		give1 []string
		give2 []string
		want  bool
	}{
		{
			name:  "exact match",
			give1: []string{"one", "two"},
			give2: []string{"one", "two"},
			want:  true,
		},
		{
			name:  "match with different indentation",
			give1: []string{"    one", "    two"},
			give2: []string{"  one", "  two"},
			want:  true,
		},
		{
			name:  "no match - different content",
			give1: []string{"one", "three"},
			give2: []string{"one", "two"},
			want:  false,
		},
		{
			name:  "no match - different length",
			give1: []string{"one"},
			give2: []string{"one", "two"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have := stringslice.EqualIgnoreWhitespace(tt.give1, tt.give2)
			must.EqOp(t, tt.want, have)
		})
	}
}
