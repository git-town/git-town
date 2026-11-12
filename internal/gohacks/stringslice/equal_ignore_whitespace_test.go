package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestEqualIgnoreWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give1 []string
		give2 []string
		name  string
		want  bool
	}{
		{
			give1: []string{"one", "two"},
			give2: []string{"one", "two"},
			name:  "exact match",
			want:  true,
		},
		{
			give1: []string{"    one", "    two"},
			give2: []string{"  one", "  two"},
			name:  "match with different indentation",
			want:  true,
		},
		{
			give1: []string{"one", "three"},
			give2: []string{"one", "two"},
			name:  "no match - different content",
			want:  false,
		},
		{
			give1: []string{"one"},
			give2: []string{"one", "two"},
			name:  "no match - different length",
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
