package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

func TestIndentNonEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		giveLines  []string
		giveIndent string
		wantLines  []string
	}{
		{
			name:       "add spaces",
			giveLines:  []string{"one", "two"},
			giveIndent: "    ",
			wantLines:  []string{"    one", "    two"},
		},
		{
			name:       "add tabs",
			giveLines:  []string{"one", "two"},
			giveIndent: "\t\t",
			wantLines:  []string{"\t\tone", "\t\ttwo"},
		},
		{
			name:       "no indentation",
			giveLines:  []string{"one", "two"},
			giveIndent: "",
			wantLines:  []string{"one", "two"},
		},
		{
			name:       "preserve empty lines",
			giveLines:  []string{"one", "", "two"},
			giveIndent: "  ",
			wantLines:  []string{"  one", "", "  two"},
		},
		{
			name:       "remove existing indentation and add new",
			giveLines:  []string{"  one", "    two"},
			giveIndent: "\t",
			wantLines:  []string{"\tone", "\ttwo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := stringslice.ChangeIndentNonEmpty(tt.giveLines, tt.giveIndent)
			if len(result) != len(tt.wantLines) {
				t.Errorf("indentTableLines() returned %d lines, expected %d", len(result), len(tt.wantLines))
				return
			}
			for i, line := range result {
				if line != tt.wantLines[i] {
					t.Errorf("indentTableLines()[%d] = %q, expected %q", i, line, tt.wantLines[i])
				}
			}
		})
	}
}
