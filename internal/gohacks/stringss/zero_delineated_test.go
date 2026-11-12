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
		"\x00leading":         {"", "leading"},
		"a\x00\x00double":     {"a", "", "double"},
		"a\x00b":              {"a", "b"},
		"a\x00b\x00c":         {"a", "b", "c"},
		"mixed\r\nchars\x001": {"mixed\r\nchars", "1"},
		"single":              {"single"},
		"trailing\x00":        {"trailing", ""},
	}
	for give, want := range tests {
		t.Run(give.String(), func(t *testing.T) {
			t.Parallel()
			have := give.Lines()
			must.Eq(t, want, have)
		})
	}
}
