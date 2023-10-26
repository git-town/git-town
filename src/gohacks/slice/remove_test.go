package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestRemove(t *testing.T) {
	t.Parallel()

	t.Run("slice type", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "three"}
		have := slice.Remove(give, "two")
		want := []string{"one", "three"}
		must.Eq(t, want, have)
	})

	t.Run("slice alias type", func(t *testing.T) {
		t.Parallel()
		give := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("222222"), domain.NewSHA("333333")}
		have := slice.Remove(give, domain.NewSHA("222222"))
		want := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("333333")}
		must.Eq(t, want, have)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		give := []string{}
		have := slice.Remove(give, "foo")
		want := []string{}
		must.Eq(t, want, have)
	})
}
