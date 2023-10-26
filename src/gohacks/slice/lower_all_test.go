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
		list := []string{"one", "two", "last"}
		slice.LowerAll(&list, "last")
		want := []string{"one", "two", "last"}
		must.Eq(t, want, list)
	})

	t.Run("list contains element in the middle", func(t *testing.T) {
		t.Parallel()
		list := []shared.Opcode{
			&opcode.AbortMerge{},
			&opcode.RestoreOpenChanges{},
			&opcode.AbortRebase{},
		}
		slice.LowerAll[shared.Opcode](&list, &opcode.RestoreOpenChanges{})
		want := []shared.Opcode{
			&opcode.AbortMerge{},
			&opcode.AbortRebase{},
			&opcode.RestoreOpenChanges{},
		}
		must.Eq(t, want, list)
	})

	t.Run("list does not contain the element", func(t *testing.T) {
		t.Parallel()
		list := []string{"one", "two", "three"}
		slice.LowerAll(&list, "last")
		want := []string{"one", "two", "three"}
		must.Eq(t, want, list)
	})

	t.Run("complex example", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 1, 3, 1}
		slice.LowerAll(&list, 1)
		want := []int{2, 3, 1}
		must.Eq(t, want, list)
	})
}
