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

	t.Run("RotateCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("currently selected alias doesn't exist", func(t *testing.T) {
			t.Parallel()
			model := dialog.AliasesModel{
				CurrentSelections: []dialog.AliasSelection{
					dialog.AliasSelectionNone,
				},
				OriginalSelections: []dialog.AliasSelection{
					dialog.AliasSelectionNone,
				},
				BubbleList: dialog.BubbleList{
					Cursor: 0,
				},
			}
			// rotate the first time to set to "checked"
			model.RotateCurrentEntry()
			want := []dialog.AliasSelection{
				dialog.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate a second time to set to "unchecked"
			model.RotateCurrentEntry()
			want = []dialog.AliasSelection{
				dialog.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
		t.Run("currently selected alias is set to the Git Town command", func(t *testing.T) {
			t.Parallel()
			model := dialog.AliasesModel{
				CurrentSelections: []dialog.AliasSelection{
					dialog.AliasSelectionGT,
				},
				OriginalSelections: []dialog.AliasSelection{
					dialog.AliasSelectionGT,
				},
				BubbleList: dialog.BubbleList{
					Cursor: 0,
				},
			}
			// rotate the first time to uncheck
			model.RotateCurrentEntry()
			want := []dialog.AliasSelection{
				dialog.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate the second time to check again
			model.RotateCurrentEntry()
			want = []dialog.AliasSelection{
				dialog.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
		t.Run("currently selected alias is set to an external command", func(t *testing.T) {
			t.Parallel()
			model := dialog.AliasesModel{
				CurrentSelections: []dialog.AliasSelection{
					dialog.AliasSelectionOther,
				},
				OriginalSelections: []dialog.AliasSelection{
					dialog.AliasSelectionOther,
				},
				BubbleList: dialog.BubbleList{
					Cursor: 0,
				},
			}
			// rotate the first time to check
			model.RotateCurrentEntry()
			want := []dialog.AliasSelection{
				dialog.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate the second time to uncheck
			model.RotateCurrentEntry()
			want = []dialog.AliasSelection{
				dialog.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate a third time to set to "other" again
			model.RotateCurrentEntry()
			want = []dialog.AliasSelection{
				dialog.AliasSelectionOther,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
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
