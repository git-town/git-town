package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestSlice(t *testing.T) {
	t.Parallel()
	tests := map[string][]string{
		"":                {},
		"multiple\nlines": {"multiple", "lines"},
		"single line":     {"single line"},
	}
	for give, want := range tests {
		have := stringslice.Lines(give)
		must.Eq(t, want, have)
	}
}
