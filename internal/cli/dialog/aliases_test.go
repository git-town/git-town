package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/pkg/keys"
	"github.com/shoenig/test/must"
)

func TestAliases(t *testing.T) {
	t.Parallel()

	t.Run("AliasResult", func(t *testing.T) {
		t.Parallel()
		allAliasableCommands := keys.AliasableCommands{
			keys.AliasableCommandAppend,
			keys.AliasableCommandHack,
			keys.AliasableCommandPropose,
			keys.AliasableCommandSync,
		}
		existingAliases := configdomain.Aliases{
			keys.AliasableCommandAppend:  "town append",
			keys.AliasableCommandHack:    "other hack",
			keys.AliasableCommandPropose: "town propose",
			keys.AliasableCommandSync:    "other sync",
		}
		selections := []dialog.AliasSelection{
			dialog.AliasSelectionGT,
			dialog.AliasSelectionGT,
			dialog.AliasSelectionNone,
			dialog.AliasSelectionOther,
		}
		have := dialog.DetermineAliasResult(selections, allAliasableCommands, existingAliases)
		want := configdomain.Aliases{
			keys.AliasableCommandAppend: "town append",
			keys.AliasableCommandHack:   "town hack",
			keys.AliasableCommandSync:   "other sync",
		}
		must.Eq(t, want, have)
	})

	t.Run("AliasSelectionText", func(t *testing.T) {
		t.Parallel()
		t.Run("all commands selected", func(t *testing.T) {
			t.Parallel()
			give := keys.AllAliasableCommands()
			have := dialog.DetermineAliasSelectionText(give)
			want := "(all)"
			must.EqOp(t, want, have)
		})
		t.Run("no commands selected", func(t *testing.T) {
			t.Parallel()
			give := keys.AliasableCommands{}
			have := dialog.DetermineAliasSelectionText(give)
			want := "(none)"
			must.EqOp(t, want, have)
		})
		t.Run("some commands selected", func(t *testing.T) {
			t.Parallel()
			give := keys.AliasableCommands{
				keys.AliasableCommandAppend,
				keys.AliasableCommandHack,
				keys.AliasableCommandSync,
			}
			have := dialog.DetermineAliasSelectionText(give)
			want := "append, hack, sync"
			must.EqOp(t, want, have)
		})
	})

	t.Run("Checked", func(t *testing.T) {
		t.Parallel()
		model := dialog.AliasesModel{
			AllAliasableCommands: keys.AliasableCommands{
				keys.AliasableCommandAppend,
				keys.AliasableCommandHack,
				keys.AliasableCommandShip,
				keys.AliasableCommandSync,
			},
			CurrentSelections: []dialog.AliasSelection{
				dialog.AliasSelectionGT,
				dialog.AliasSelectionGT,
				dialog.AliasSelectionNone,
				dialog.AliasSelectionOther,
			},
		}
		have := model.Checked()
		want := keys.AliasableCommands{
			keys.AliasableCommandAppend,
			keys.AliasableCommandHack,
		}
		must.Eq(t, want, have)
	})

	t.Run("NewAliasSelections", func(t *testing.T) {
		t.Parallel()
		allAliasableCommands := keys.AliasableCommands{
			keys.AliasableCommandAppend,
			keys.AliasableCommandDiffParent,
			keys.AliasableCommandHack,
			keys.AliasableCommandRepo,
		}
		existingAliases := configdomain.Aliases{
			keys.AliasableCommandAppend:     "town append",
			keys.AliasableCommandDiffParent: "town diff-parent",
			keys.AliasableCommandRepo:       "other command",
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
		t.Run("currently selecting an alias that isn't currently set on disk", func(t *testing.T) {
			t.Parallel()
			model := dialog.AliasesModel{
				AllAliasableCommands: keys.AliasableCommands{
					keys.AliasableCommandAppend,
				},
				CurrentSelections: []dialog.AliasSelection{
					dialog.AliasSelectionNone,
				},
				OriginalAliases: configdomain.Aliases{},
				List: list.List[keys.AliasableCommand]{
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
		t.Run("currently selecting an alias that is set on disk to a Git Town command", func(t *testing.T) {
			t.Parallel()
			model := dialog.AliasesModel{
				AllAliasableCommands: keys.AliasableCommands{
					keys.AliasableCommandAppend,
				},
				CurrentSelections: []dialog.AliasSelection{
					dialog.AliasSelectionGT,
				},
				OriginalAliases: configdomain.Aliases{
					keys.AliasableCommandAppend: "town append",
				},
				List: list.List[keys.AliasableCommand]{
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
		t.Run("currently selecting an alias that is currently set on disk to an external command", func(t *testing.T) {
			t.Parallel()
			model := dialog.AliasesModel{
				AllAliasableCommands: keys.AliasableCommands{
					keys.AliasableCommandAppend,
				},
				CurrentSelections: []dialog.AliasSelection{
					dialog.AliasSelectionOther,
				},
				OriginalAliases: configdomain.Aliases{
					keys.AliasableCommandAppend: "other command",
				},
				List: list.List[keys.AliasableCommand]{
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
			List: list.List[keys.AliasableCommand]{
				Entries: list.NewEntries(
					keys.AliasableCommandAppend,
					keys.AliasableCommandHack,
					keys.AliasableCommandSync,
				),
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
			List: list.List[keys.AliasableCommand]{
				Entries: list.NewEntries(
					keys.AliasableCommandAppend,
					keys.AliasableCommandHack,
					keys.AliasableCommandDiffParent,
				),
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
