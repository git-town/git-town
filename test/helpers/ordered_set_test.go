package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/shoenig/test"
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
			test.Eq(t, want, have)
		})
		t.Run("element already exists in set", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two")
			set = set.Add("two")
			have := set.Elements()
			want := []string{"one", "two"}
			test.Eq(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		set := helpers.NewOrderedSet("one", "two")
		test.True(t, set.Contains("one"))
		test.True(t, set.Contains("two"))
		test.False(t, set.Contains("zonk"))
	})

	t.Run("Elements", func(t *testing.T) {
		t.Parallel()
		set := helpers.NewOrderedSet("one", "two")
		have := set.Elements()
		want := []string{"one", "two"}
		test.Eq(t, want, have)
	})

	t.Run("Join", func(t *testing.T) {
		t.Parallel()
		t.Run("strings", func(t *testing.T) {
			set := helpers.NewOrderedSet("one", "two", "three")
			have := set.Join(", ")
			want := "one, two, three"
			test.EqOp(t, want, have)
		})
		t.Run("ints", func(t *testing.T) {
			set := helpers.NewOrderedSet(1, 2, 3)
			have := set.Join(", ")
			want := "1, 2, 3"
			test.EqOp(t, want, have)
		})
		t.Run("SHAs", func(t *testing.T) {
			set := helpers.NewOrderedSet(
				domain.NewSHA("111111"),
				domain.NewSHA("222222"),
				domain.NewSHA("333333"),
			)
			have := set.Join(", ")
			want := "111111, 222222, 333333"
			test.EqOp(t, want, have)
		})
	})
}
