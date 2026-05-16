package stringss_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestLeadingWhitespace(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"    text":   "    ",   // spaces
		"\t\ttext":   "\t\t",   // tabs
		"  \t  text": "  \t  ", // mixed spaces and tabs
		"text":       "",       // no indentation
		"":           "",       // empty string
		"    ":       "    ",   // only whitespace
	}
	for give, want := range tests {
		have := stringss.LeadingWhitespace(give)
		must.EqOp(t, want, have)
	}
}
