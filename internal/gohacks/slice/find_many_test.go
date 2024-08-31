package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFindMany(t *testing.T) {
	t.Parallel()

	t.Run("haystack contains all needles", func(t *testing.T) {
		t.Parallel()
		haystack := []string{"one", "two", "three"}
		needles := []string{"two", "three"}
		have := slice.FindMany(haystack, needles)
		want := []int{1, 2}
		must.Eq(t, want, have)
	})

	t.Run("haystack is missing some needles", func(t *testing.T) {
		t.Parallel()
		haystack := []string{"one", "two", "three"}
		needles := []string{"two", "four"}
		have := slice.FindMany(haystack, needles)
		want := []int{1}
		must.Eq(t, want, have)
	})

	t.Run("haystack is empty", func(t *testing.T) {
		t.Parallel()
		haystack := []string(nil)
		needles := []string{"one", "two"}
		have := slice.FindMany(haystack, needles)
		must.Len(t, 0, have)
	})

	t.Run("no needles given", func(t *testing.T) {
		t.Parallel()
		haystack := []string{"one", "two", "three"}
		needles := []string(nil)
		have := slice.FindMany(haystack, needles)
		must.Len(t, 0, have)
	})
}
