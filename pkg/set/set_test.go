package set_test

import (
	"testing"

	"github.com/git-town/git-town/v15/pkg/set"
	"github.com/shoenig/test/must"
)

func TestSet(t *testing.T) {
	t.Parallel()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()
		set := set.New[int]()
		set.Add(1)
		must.Eq(t, []int{1}, set.Values())
		set.Add(1)
		must.Eq(t, []int{1}, set.Values())
		set.Add(2)
		must.Eq(t, []int{1, 2}, set.Values())
	})

	t.Run("AddSet", func(t *testing.T) {
		t.Parallel()
		have := set.New(1)
		other := set.New(2, 3)
		have.AddSet(other)
		must.Eq(t, []int{1, 2, 3}, have.Values())
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		set := set.New(1, 2)
		must.True(t, set.Contains(1))
		must.True(t, set.Contains(2))
		must.False(t, set.Contains(3))
	})

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("no initial value", func(t *testing.T) {
			t.Parallel()
			set := set.New[int]()
			must.Eq(t, []int{}, set.Values())
		})

		t.Run("one initial value", func(t *testing.T) {
			t.Parallel()
			set := set.New(1)
			must.Eq(t, []int{1}, set.Values())
		})

		t.Run("multiple initial values", func(t *testing.T) {
			t.Parallel()
			set := set.New(1, 2)
			must.Eq(t, []int{1, 2}, set.Values())
		})
	})

	t.Run("Values", func(t *testing.T) {
		t.Parallel()

		t.Run("no values", func(t *testing.T) {
			set := set.New[int]()
			have := set.Values()
			want := []int{}
			must.Eq(t, want, have)
		})

		t.Run("with values", func(t *testing.T) {
			set := set.New(1, 2)
			have := set.Values()
			want := []int{1, 2}
			must.Eq(t, want, have)
		})
	})
}
