package stringers_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/stringers"
	"github.com/shoenig/test/must"
)

func TestIndexOrStart(t *testing.T) {
	t.Parallel()

	t.Run("haystack contains the needle", func(t *testing.T) {
		t.Parallel()
		one := testEntry("one")
		two := testEntry("two")
		list := []testEntry{one, two}
		have := stringers.IndexOrStart(list, one)
		must.Eq(t, 0, have)
		have = stringers.IndexOrStart(list, two)
		must.Eq(t, 1, have)
	})

	t.Run("haystack does not contain the needle", func(t *testing.T) {
		t.Parallel()
		one := testEntry("one")
		two := testEntry("two")
		list := []testEntry{one}
		have := stringers.IndexOrStart(list, two)
		must.Eq(t, 0, have)
	})

	t.Run("empty haystack", func(t *testing.T) {
		t.Parallel()
		one := testEntry("one")
		list := []testEntry{}
		have := stringers.IndexOrStart(list, one)
		must.Eq(t, 0, have)
	})
}
