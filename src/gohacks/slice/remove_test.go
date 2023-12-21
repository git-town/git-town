package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestRemove(t *testing.T) {
	t.Parallel()

	t.Run("slice type", func(t *testing.T) {
		t.Parallel()
		list := []string{"one", "two", "three"}
		slice.Remove(&list, "two")
		want := []string{"one", "three"}
		must.Eq(t, want, list)
	})

	t.Run("slice alias type", func(t *testing.T) {
		t.Parallel()
		list := gitdomain.SHAs{gitdomain.NewSHA("111111"), gitdomain.NewSHA("222222"), gitdomain.NewSHA("333333")}
		slice.Remove(&list, gitdomain.NewSHA("222222"))
		want := gitdomain.SHAs{gitdomain.NewSHA("111111"), gitdomain.NewSHA("333333")}
		must.Eq(t, want, list)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		list := []string{}
		slice.Remove(&list, "foo")
		want := []string{}
		must.Eq(t, want, list)
	})
}
