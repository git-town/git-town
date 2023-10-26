package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestLowerAll(t *testing.T) {
	t.Parallel()
	t.Run("list contains element at the last position", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "last"}
		have := slice.LowerAll(give, "last")
		want := []string{"one", "two", "last"}
		must.Eq(t, want, have)
	})

	t.Run("list contains element in the middle", func(t *testing.T) {
		t.Parallel()
		give := []shared.Opcode{
			&opcode.AbortMerge{},
			&opcode.RestoreOpenChanges{},
			&opcode.AbortRebase{},
		}
		have := slice.LowerAll[shared.Opcode](give, &opcode.RestoreOpenChanges{})
		want := []shared.Opcode{
			&opcode.AbortMerge{},
			&opcode.AbortRebase{},
			&opcode.RestoreOpenChanges{},
		}
		must.Eq(t, want, have)
	})

	t.Run("list does not contain the element", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "three"}
		have := slice.LowerAll(give, "last")
		want := []string{"one", "two", "three"}
		must.Eq(t, want, have)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		give := []int{1, 2, 1, 3, 1}
		have := slice.LowerAll(give, 1)
		want := []int{2, 3, 1}
		must.Eq(t, want, have)
	})
}
