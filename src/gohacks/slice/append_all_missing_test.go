package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestAppendAllMissing(t *testing.T) {
	t.Parallel()

	t.Run("slice type", func(t *testing.T) {
		t.Parallel()
		list := []string{"one", "two", "three"}
		give := []string{"two", "four", "five"}
		have := slice.AppendAllMissing(list, give)
		want := []string{"one", "two", "three", "four", "five"}
		must.Eq(t, want, have)
	})

	t.Run("aliased slice type", func(t *testing.T) {
		t.Parallel()
		list := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("222222")}
		give := domain.SHAs{domain.NewSHA("333333"), domain.NewSHA("444444")}
		have := slice.AppendAllMissing(list, give)
		want := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("222222"), domain.NewSHA("333333"), domain.NewSHA("444444")}
		must.Eq(t, want, have)
	})
}
