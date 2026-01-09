package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/shoenig/test/must"
)

func TestAn(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"apple":    "an",
		"elephant": "an",
		"igloo":    "an",
		"orange":   "an",
		"umbrella": "an",
		"book":     "a",
		"car":      "a",
		"dog":      "a",
		"house":    "a",
		"tree":     "a",
		"zebra":    "a",
		"":         "a",
	}

	for give, want := range tests {
		have := gohacks.An(give)
		must.EqOp(t, want, have)
	}
}
