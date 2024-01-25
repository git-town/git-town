package stringers_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/gohacks/stringers"
	"github.com/shoenig/test/must"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	t.Run("haystack contains the needle", func(t *testing.T) {
		t.Parallel()
		one := TestEntry("one")
		two := TestEntry("two")
		list := []TestEntry{one, two}
		have := stringers.Index(list, one)
		must.Eq(t, 0, have)
		have = stringers.Index(list, two)
		must.Eq(t, 1, have)
	})
	t.Run("haystack does not contain the needle", func(t *testing.T) {
		t.Parallel()
		one := TestEntry("one")
		two := TestEntry("two")
		list := []TestEntry{one}
		have := stringers.Index(list, two)
		must.Eq(t, -1, have)
	})
	t.Run("empty haystack", func(t *testing.T) {
		t.Parallel()
		one := TestEntry("one")
		list := []TestEntry{}
		have := stringers.Index(list, one)
		must.Eq(t, -1, have)
	})
}

type TestEntry string

func (self TestEntry) String() string {
	return string(self)
}
