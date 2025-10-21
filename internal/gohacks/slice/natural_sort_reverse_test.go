package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestNaturalSortReverse(t *testing.T) {
	t.Parallel()
	tests := map[*[]stringer]*[]stringer{
		{}:                       {},                       // empty
		{"1"}:                    {"1"},                    // single element
		{"a3c", "a20b", "a100a"}: {"a100a", "a20b", "a3c"}, // ordering by numeric value (reversed)
		{"a10b2", "a10b10"}:      {"a10b10", "a10b2"},      // multiple parts of numbers and characters (reversed)
	}
	for give, want := range tests {
		slice.NaturalSortReverse(*give)
		must.Eq(t, want, give)
	}
}
