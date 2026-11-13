package stringss_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestEscapeNewLines(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"":                    "",
		"no newlines":         "no newlines",
		"one\nnewline":        "one\\nnewline",
		"three\nnew\nlines\n": "three\\nnew\\nlines\\n",
		"two\nnew\nlines":     "two\\nnew\\nlines",
	}
	for give, want := range tests {
		have := stringss.EscapeNewLines(give)
		must.EqOp(t, want, have)
	}
}
