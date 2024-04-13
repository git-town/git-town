package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestSlice(t *testing.T) {
	t.Parallel()
	tests := map[string][]string{
		"":                {},
		"single line":     {"single line"},
		"multiple\nlines": {"multiple", "lines"},
	}
	for give, want := range tests {
		have := stringslice.Lines(give)
		must.Eq(t, want, have)
	}
}
