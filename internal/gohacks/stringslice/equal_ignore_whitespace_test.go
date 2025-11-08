package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

func TestEqualIgnoreWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fileLines  []string
		tableLines []string
		want       bool
	}{
		{
			name:       "exact match",
			fileLines:  []string{"one", "two"},
			tableLines: []string{"one", "two"},
			want:       true,
		},
		{
			name:       "match with different indentation",
			fileLines:  []string{"    one", "    two"},
			tableLines: []string{"  one", "  two"},
			want:       true,
		},
		{
			name:       "no match - different content",
			fileLines:  []string{"one", "three"},
			tableLines: []string{"one", "two"},
			want:       false,
		},
		{
			name:       "no match - different length",
			fileLines:  []string{"one"},
			tableLines: []string{"one", "two"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := stringslice.EqualIgnoreWhitespace(tt.fileLines, tt.tableLines)
			if result != tt.want {
				t.Errorf("matchesTable() = %v, expected %v", result, tt.want)
			}
		})
	}
}
