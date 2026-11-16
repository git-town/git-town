package stringss_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestAn(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"":         "a",
		"apple":    "an",
		"book":     "a",
		"car":      "a",
		"dog":      "a",
		"elephant": "an",
		"house":    "a",
		"igloo":    "an",
		"orange":   "an",
		"tree":     "a",
		"umbrella": "an",
		"zebra":    "a",
	}

	for give, want := range tests {
		have := stringss.An(give)
		must.EqOp(t, want, have)
	}
}
