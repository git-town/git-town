package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestAppendAllMissing(t *testing.T) {
	t.Parallel()

	t.Run("slice type", func(t *testing.T) {
		t.Parallel()
		list := []string{"one", "two", "three"}
		additional := []string{"two", "four", "five"}
		slice.AppendAllMissing(&list, additional)
		want := []string{"one", "two", "three", "four", "five"}
		must.Eq(t, want, list)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		list := make([]string, 0)
		additional := []string{"one", "two", "three"}
		slice.AppendAllMissing(&list, additional)
		want := []string{"one", "two", "three"}
		must.Eq(t, want, list)
	})

	t.Run("zero slice", func(t *testing.T) {
		t.Parallel()
		var list []string
		additional := []string{"one", "two", "three"}
		slice.AppendAllMissing(&list, additional)
		want := []string{"one", "two", "three"}
		must.Eq(t, want, list)
	})

	t.Run("aliased slice type", func(t *testing.T) {
		t.Parallel()
		list := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("222222")}
		give := domain.SHAs{domain.NewSHA("333333"), domain.NewSHA("444444")}
		slice.AppendAllMissing(&list, give)
		want := domain.SHAs{domain.NewSHA("111111"), domain.NewSHA("222222"), domain.NewSHA("333333"), domain.NewSHA("444444")}
		must.Eq(t, want, list)
	})
}
