package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestSortStringers(t *testing.T) {
	t.Parallel()
	tests := map[*[]stringer]*[]stringer{
		// empty
		newStringers(): newStringers(),
		// single element
		newStringers("a"): newStringers("a"),
		// multiple elements
		newStringers("b20", "b1", "a2"): newStringers("a2", "b1", "b20"),
	}
	for give, want := range tests {
		have := slice.NatSort(*give)
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
