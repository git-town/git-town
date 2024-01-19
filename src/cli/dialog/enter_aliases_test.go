package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestEnterAliases(t *testing.T) {
	t.Parallel()
	t.Run("SelectAll", func(t *testing.T) {
		t.Parallel()
		model := dialog.AliasesModel{
			BubbleList: dialog.BubbleList{
				Entries: []string{"append", "hack", "diff-parent"},
			},
			Selections: []int{},
		}
		model.SelectAll()
		want := []int{0, 1, 2}
		must.Eq(t, model.Selections, want)
	})
}
