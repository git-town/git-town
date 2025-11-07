package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/shoenig/test/must"
)

func TestGetIndentation(t *testing.T) {
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
		have := gohacks.GetIndentation(give)
		must.EqOp(t, want, have)
	}
}

func TestEscapeNewLines(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"":                    "",
		"no newlines":         "no newlines",
		"one\nnewline":        "one\\nnewline",
		"two\nnew\nlines":     "two\\nnew\\nlines",
		"three\nnew\nlines\n": "three\\nnew\\nlines\\n",
	}
	for give, want := range tests {
		have := gohacks.EscapeNewLines(give)
		must.EqOp(t, want, have)
	}
}

func TestIndent(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		have := gohacks.IndentLines("", 2)
		must.EqOp(t, "  ", have)
	})

	t.Run("indent by 2 spaces", func(t *testing.T) {
		t.Parallel()
		give := "one\ntwo\nthree"
		have := gohacks.IndentLines(give, 2)
		want := "  one\n  two\n  three"
		must.EqOp(t, want, have)
	})

	t.Run("indent by 4 spaces", func(t *testing.T) {
		t.Parallel()
		give := "one\ntwo\nthree"
		have := gohacks.IndentLines(give, 4)
		want := "    one\n    two\n    three"
		must.EqOp(t, want, have)
	})
}
