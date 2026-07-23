package stringss_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestZeroDelineated(t *testing.T) {
	t.Parallel()
	tests := map[stringss.ZeroDelineated][]string{
		"":                    {""},
		"single":              {"single"},
		"a\x00b":              {"a", "b"},
		"a\x00b\x00c":         {"a", "b", "c"},
		"\x00leading":         {"", "leading"},
		"trailing\x00":        {"trailing", ""},
		"a\x00\x00double":     {"a", "", "double"},
		"mixed\r\nchars\x001": {"mixed\r\nchars", "1"},
	}
	for give, want := range tests {
		t.Run(give.String(), func(t *testing.T) {
			t.Parallel()
			have := give.Lines()
			must.Eq(t, want, have)
		})
	}
}
