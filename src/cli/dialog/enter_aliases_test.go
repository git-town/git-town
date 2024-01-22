package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestEnterAliases(t *testing.T) {
	t.Parallel()

	t.Run("AliasSelectionText", func(t *testing.T) {
		t.Parallel()
		t.Run("all commands selected", func(t *testing.T) {
			t.Parallel()
			give := configdomain.AllAliasableCommands()
			have := dialog.AliasSelectionText(give)
			want := "(all)"
			must.EqOp(t, want, have)
		})
		t.Run("no commands selected", func(t *testing.T) {
			t.Parallel()
			give := configdomain.AliasableCommands{}
			have := dialog.AliasSelectionText(give)
			want := "(none)"
			must.EqOp(t, want, have)
		})
		t.Run("some commands selected", func(t *testing.T) {
			t.Parallel()
			give := configdomain.AliasableCommands{
				configdomain.AliasableCommandAppend,
				configdomain.AliasableCommandHack,
				configdomain.AliasableCommandSync,
			}
			have := dialog.AliasSelectionText(give)
			want := "append, hack, sync"
			must.EqOp(t, want, have)
		})
	})

	t.Run("Checked", func(t *testing.T) {
		t.Parallel()
		model := dialog.AliasesModel{
			CurrentSelections: []dialog.AliasSelection{
				dialog.AliasSelectionGT,
				dialog.AliasSelectionGT,
				dialog.AliasSelectionNone,
				dialog.AliasSelectionOther,
			},
		}
		aliasableCommands := configdomain.AliasableCommands{
			configdomain.AliasableCommandAppend,
			configdomain.AliasableCommandHack,
			configdomain.AliasableCommandShip,
			configdomain.AliasableCommandSync,
		}
		have := model.Checked(aliasableCommands)
		want := configdomain.AliasableCommands{
			configdomain.AliasableCommandAppend,
			configdomain.AliasableCommandHack,
		}
		must.Eq(t, want, have)
	})

	t.Run("NewAliasSelections", func(t *testing.T) {
		t.Parallel()
		allAliasableCommands := configdomain.AliasableCommands{
			configdomain.AliasableCommandAppend,
			configdomain.AliasableCommandDiffParent,
			configdomain.AliasableCommandHack,
			configdomain.AliasableCommandRepo,
		}
		existingAliases := configdomain.Aliases{
			configdomain.AliasableCommandAppend:     "town append",
			configdomain.AliasableCommandDiffParent: "town diff-parent",
			configdomain.AliasableCommandRepo:       "other command",
		}
		have := dialog.NewAliasSelections(allAliasableCommands, existingAliases)
		want := []dialog.AliasSelection{
			dialog.AliasSelectionGT,
			dialog.AliasSelectionGT,
			dialog.AliasSelectionNone,
			dialog.AliasSelectionOther,
		}
		must.Eq(t, want, have)
	})

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
