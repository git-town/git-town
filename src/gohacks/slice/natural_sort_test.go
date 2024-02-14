package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestNaturalSort(t *testing.T) {
	t.Parallel()
	tests := map[*[]stringer]*[]stringer{
		newStringers():                    newStringers(),                    // empty
		newStringers("a"):                 newStringers("a"),                 // single element
		newStringers("a100", "a20", "a3"): newStringers("a3", "a20", "a100"), // multiple elements
		newStringers("a10b", "a10a"):      newStringers("a10a", "a10b"),      // multiple elements
	}
	for give, want := range tests {
		have := slice.NaturalSort(*give)
		must.Eq(t, want, &have)
	}
}

type stringer struct {
	s string
}

func (s stringer) String() string {
	return s.s
}

func newStringers(names ...string) *[]stringer {
	result := make([]stringer, len(names))
	for n, name := range names {
		result[n] = stringer{name}
	}
	return &result
}
