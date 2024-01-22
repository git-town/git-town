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
			CurrentSelections: []dialog.AliasSelection{
				dialog.AliasSelectionNone,
				dialog.AliasSelectionNone,
				dialog.AliasSelectionNone,
			},
		}
		model.SelectAll()
		want := []dialog.AliasSelection{
			dialog.AliasSelectionGT,
			dialog.AliasSelectionGT,
			dialog.AliasSelectionGT,
		}
		must.Eq(t, model.CurrentSelections, want)
	})

	t.Run("SelectNone", func(t *testing.T) {
		t.Parallel()
		model := dialog.AliasesModel{
			BubbleList: dialog.BubbleList{
				Entries: []string{"append", "hack", "diff-parent"},
			},
			CurrentSelections: []dialog.AliasSelection{
				dialog.AliasSelectionGT,
				dialog.AliasSelectionOther,
				dialog.AliasSelectionGT,
			},
		}
		model.SelectNone()
		want := []dialog.AliasSelection{
			dialog.AliasSelectionNone,
			dialog.AliasSelectionNone,
			dialog.AliasSelectionNone,
		}
		must.Eq(t, model.CurrentSelections, want)
	})
}
