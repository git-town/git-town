package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFirstNonEmpty(t *testing.T) {
	t.Parallel()
	t.Run("one element", func(t *testing.T) {
		t.Parallel()
		one := gitdomain.CommitMessage("one")
		have := slice.FirstNonEmpty(one)
		want := one
		must.EqOp(t, want, have)
	})
	t.Run("first element is non-empty", func(t *testing.T) {
		t.Parallel()
		one := gitdomain.CommitMessage("one")
		empty := gitdomain.CommitMessage("")
		have := slice.FirstNonEmpty(one, empty)
		want := one
		must.EqOp(t, want, have)
	})
	t.Run("second element is non-empty", func(t *testing.T) {
		t.Parallel()
		empty := gitdomain.CommitMessage("")
		two := gitdomain.CommitMessage("two")
		have := slice.FirstNonEmpty(empty, two)
		want := two
		must.EqOp(t, want, have)
	})
	t.Run("third element is non-empty", func(t *testing.T) {
		t.Parallel()
		empty := gitdomain.CommitMessage("")
		three := gitdomain.CommitMessage("three")
		have := slice.FirstNonEmpty(empty, empty, three)
		want := three
		must.EqOp(t, want, have)
	})
}
