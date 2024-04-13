package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestNaturalSort(t *testing.T) {
	t.Parallel()
	tests := map[*[]stringer]*[]stringer{
		{}:                       {},                       // empty
		{"1"}:                    {"1"},                    // single element
		{"a100a", "a20b", "a3c"}: {"a3c", "a20b", "a100a"}, // ordering by numeric value
		{"a10b10", "a10b2"}:      {"a10b2", "a10b10"},      // multiple parts of numbers and characters
	}
	for give, want := range tests {
		have := slice.NaturalSort(*give)
		must.Eq(t, want, &have)
	}
}

type stringer string

func (self stringer) String() string {
	return string(self)
}
