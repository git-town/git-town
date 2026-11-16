package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestIndentNonEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		giveIndent string
		giveLines  []string
		name       string
		wantLines  []string
	}{
		{
			giveIndent: "    ",
			giveLines:  []string{"one", "two"},
			name:       "add spaces",
			wantLines:  []string{"    one", "    two"},
		},
		{
			giveIndent: "\t\t",
			giveLines:  []string{"one", "two"},
			name:       "add tabs",
			wantLines:  []string{"\t\tone", "\t\ttwo"},
		},
		{
			giveIndent: "",
			giveLines:  []string{"one", "two"},
			name:       "no indentation",
			wantLines:  []string{"one", "two"},
		},
		{
			giveIndent: "  ",
			giveLines:  []string{"one", "", "two"},
			name:       "preserve empty lines",
			wantLines:  []string{"  one", "", "  two"},
		},
		{
			giveIndent: "\t",
			giveLines:  []string{"  one", "    two"},
			name:       "remove existing indentation and add new",
			wantLines:  []string{"\tone", "\ttwo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			haveLines := stringslice.ChangeIndentNonEmpty(tt.giveLines, tt.giveIndent)
			must.Eq(t, tt.wantLines, haveLines)
		})
	}
}
