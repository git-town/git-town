package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/test/helpers"
	"github.com/shoenig/test/must"
)

func TestOrderedSet(t *testing.T) {
	t.Parallel()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()
		t.Run("element doesn't exist in set", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two")
			set = set.Add("three")
			have := set.Elements()
			want := []string{"one", "two", "three"}
			must.Eq(t, want, have)
		})
		t.Run("element already exists in set", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two")
			set = set.Add("two")
			have := set.Elements()
			want := []string{"one", "two"}
			must.Eq(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		set := helpers.NewOrderedSet("one", "two")
		must.True(t, set.Contains("one"))
		must.True(t, set.Contains("two"))
		must.False(t, set.Contains("zonk"))
	})

	t.Run("Elements", func(t *testing.T) {
		t.Parallel()
		set := helpers.NewOrderedSet("one", "two")
		have := set.Elements()
		want := []string{"one", "two"}
		must.Eq(t, want, have)
	})

	t.Run("Join", func(t *testing.T) {
		t.Parallel()
		t.Run("strings", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two", "three")
			have := set.Join(", ")
			want := "one, two, three"
			must.EqOp(t, want, have)
		})
		t.Run("ints", func(t *testing.T) {
			set := helpers.NewOrderedSet(1, 2, 3)
			have := set.Join(", ")
			want := "1, 2, 3"
			must.EqOp(t, want, have)
		})
		t.Run("SHAs", func(t *testing.T) {
			set := helpers.NewOrderedSet(
				gitdomain.NewSHA("111111"),
				gitdomain.NewSHA("222222"),
				gitdomain.NewSHA("333333"),
			)
			have := set.Join(", ")
			want := "111111, 222222, 333333"
			must.EqOp(t, want, have)
		})
	})
}
