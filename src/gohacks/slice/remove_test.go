package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
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
		list := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("222222"), domain.NewSHA("333333")}
		slice.Remove(&list, domain.NewSHA("222222"))
		want := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("333333")}
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
