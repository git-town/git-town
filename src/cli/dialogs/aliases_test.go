package dialogs_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialogs"
	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestAliases(t *testing.T) {
	t.Parallel()

	t.Run("AliasResult", func(t *testing.T) {
		t.Parallel()
		allAliasableCommands := configdomain.AliasableCommands{
			configdomain.AliasableCommandAppend,
			configdomain.AliasableCommandHack,
			configdomain.AliasableCommandPropose,
			configdomain.AliasableCommandSync,
		}
		existingAliases := configdomain.Aliases{
			configdomain.AliasableCommandAppend:  "town append",
			configdomain.AliasableCommandHack:    "other hack",
			configdomain.AliasableCommandPropose: "town propose",
			configdomain.AliasableCommandSync:    "other sync",
		}
		selections := []dialogs.AliasSelection{
			dialogs.AliasSelectionGT,
			dialogs.AliasSelectionGT,
			dialogs.AliasSelectionNone,
			dialogs.AliasSelectionOther,
		}
		have := dialogs.DetermineAliasResult(selections, allAliasableCommands, existingAliases)
		want := configdomain.Aliases{
			configdomain.AliasableCommandAppend: "town append",
			configdomain.AliasableCommandHack:   "town hack",
			configdomain.AliasableCommandSync:   "other sync",
		}
		must.Eq(t, want, have)
	})

	t.Run("AliasSelectionText", func(t *testing.T) {
		t.Parallel()
		t.Run("all commands selected", func(t *testing.T) {
			t.Parallel()
			give := configdomain.AllAliasableCommands()
			have := dialogs.DetermineAliasSelectionText(give)
			want := "(all)"
			must.EqOp(t, want, have)
		})
		t.Run("no commands selected", func(t *testing.T) {
			t.Parallel()
			give := configdomain.AliasableCommands{}
			have := dialogs.DetermineAliasSelectionText(give)
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
			have := dialogs.DetermineAliasSelectionText(give)
			want := "append, hack, sync"
			must.EqOp(t, want, have)
		})
	})

	t.Run("Checked", func(t *testing.T) {
		t.Parallel()
		model := dialogs.AliasesModel{ //nolint:exhaustruct
			AllAliasableCommands: configdomain.AliasableCommands{
				configdomain.AliasableCommandAppend,
				configdomain.AliasableCommandHack,
				configdomain.AliasableCommandShip,
				configdomain.AliasableCommandSync,
			},
			CurrentSelections: []dialogs.AliasSelection{
				dialogs.AliasSelectionGT,
				dialogs.AliasSelectionGT,
				dialogs.AliasSelectionNone,
				dialogs.AliasSelectionOther,
			},
		}
		have := model.Checked()
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
		have := dialogs.NewAliasSelections(allAliasableCommands, existingAliases)
		want := []dialogs.AliasSelection{
			dialogs.AliasSelectionGT,
			dialogs.AliasSelectionGT,
			dialogs.AliasSelectionNone,
			dialogs.AliasSelectionOther,
		}
		must.Eq(t, want, have)
	})

	t.Run("RotateCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("currently selecting an alias that isn't currently set on disk", func(t *testing.T) {
			t.Parallel()
			model := dialogs.AliasesModel{
				AllAliasableCommands: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
				},
				CurrentSelections: []dialogs.AliasSelection{
					dialogs.AliasSelectionNone,
				},
				OriginalAliases: configdomain.Aliases{},
				BubbleList: components.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
					Cursor: 0,
				},
			}
			// rotate the first time to set to "checked"
			model.RotateCurrentEntry()
			want := []dialogs.AliasSelection{
				dialogs.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate a second time to set to "unchecked"
			model.RotateCurrentEntry()
			want = []dialogs.AliasSelection{
				dialogs.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
		t.Run("currently selecting an alias that is set on disk to a Git Town command", func(t *testing.T) {
			t.Parallel()
			model := dialogs.AliasesModel{
				AllAliasableCommands: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
				},
				CurrentSelections: []dialogs.AliasSelection{
					dialogs.AliasSelectionGT,
				},
				OriginalAliases: configdomain.Aliases{
					configdomain.AliasableCommandAppend: "town append",
				},
				BubbleList: components.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
					Cursor: 0,
				},
			}
			// rotate the first time to uncheck
			model.RotateCurrentEntry()
			want := []dialogs.AliasSelection{
				dialogs.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate the second time to check again
			model.RotateCurrentEntry()
			want = []dialogs.AliasSelection{
				dialogs.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
		t.Run("currently selecting an alias that is currently set on disk to an external command", func(t *testing.T) {
			t.Parallel()
			model := dialogs.AliasesModel{
				AllAliasableCommands: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
				},
				CurrentSelections: []dialogs.AliasSelection{
					dialogs.AliasSelectionOther,
				},
				OriginalAliases: configdomain.Aliases{
					configdomain.AliasableCommandAppend: "other command",
				},
				BubbleList: components.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
					Cursor: 0,
				},
			}
			// rotate the first time to check
			model.RotateCurrentEntry()
			want := []dialogs.AliasSelection{
				dialogs.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate the second time to uncheck
			model.RotateCurrentEntry()
			want = []dialogs.AliasSelection{
				dialogs.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate a third time to set to "other" again
			model.RotateCurrentEntry()
			want = []dialogs.AliasSelection{
				dialogs.AliasSelectionOther,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
	})

	t.Run("SelectAll", func(t *testing.T) {
		t.Parallel()
		model := dialogs.AliasesModel{ //nolint:exhaustruct
			BubbleList: components.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
				Entries: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
					configdomain.AliasableCommandHack,
					configdomain.AliasableCommandSync,
				},
			},
			CurrentSelections: []dialogs.AliasSelection{
				dialogs.AliasSelectionNone,
				dialogs.AliasSelectionNone,
				dialogs.AliasSelectionNone,
			},
		}
		model.SelectAll()
		want := []dialogs.AliasSelection{
			dialogs.AliasSelectionGT,
			dialogs.AliasSelectionGT,
			dialogs.AliasSelectionGT,
		}
		must.Eq(t, model.CurrentSelections, want)
	})

	t.Run("SelectNone", func(t *testing.T) {
		t.Parallel()
		model := dialogs.AliasesModel{ //nolint:exhaustruct
			BubbleList: components.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
				Entries: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
					configdomain.AliasableCommandHack,
					configdomain.AliasableCommandDiffParent,
				},
			},
			CurrentSelections: []dialogs.AliasSelection{
				dialogs.AliasSelectionGT,
				dialogs.AliasSelectionOther,
				dialogs.AliasSelectionGT,
			},
		}
		model.SelectNone()
		want := []dialogs.AliasSelection{
			dialogs.AliasSelectionNone,
			dialogs.AliasSelectionNone,
			dialogs.AliasSelectionNone,
		}
		must.Eq(t, model.CurrentSelections, want)
	})
}
