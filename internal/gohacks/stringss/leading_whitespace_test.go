package stringss_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestLeadingWhitespace(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"":           "",       // empty string
		"\t\ttext":   "\t\t",   // tabs
		"  \t  text": "  \t  ", // mixed spaces and tabs
		"    ":       "    ",   // only whitespace
		"    text":   "    ",   // spaces
		"text":       "",       // no indentation
	}
	for give, want := range tests {
		have := stringss.LeadingWhitespace(give)
		must.EqOp(t, want, have)
	}
}
