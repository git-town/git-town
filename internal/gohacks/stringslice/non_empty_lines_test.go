package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestNonEmptyLines(t *testing.T) {
	t.Parallel()
	tests := map[string][]string{
		"":                                {},
		"hello":                           {"hello"},
		"line1\nline2\nline3":             {"line1", "line2", "line3"},
		"line1\n\nline2\n\nline3":         {"line1", "line2", "line3"},
		"line1\n   \nline2\n\t\nline3":    {"line1", "line2", "line3"},
		"  line1  \n  line2  \n  line3  ": {"line1", "line2", "line3"},
		"line1\n\n   \nline2\n\t\t\n  line3  \n\n": {"line1", "line2", "line3"},
		"line1\nline2\n": {"line1", "line2"},
		"\n\n\n":         {},
		"   \n\t\n  ":    {},
	}
	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			have := stringslice.NonEmptyLines(input)
			must.Eq(t, expected, have)
		})
	}
}
