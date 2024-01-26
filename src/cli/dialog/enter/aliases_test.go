package enter_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/dialog/enter"
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
		selections := []enter.AliasSelection{
			enter.AliasSelectionGT,
			enter.AliasSelectionGT,
			enter.AliasSelectionNone,
			enter.AliasSelectionOther,
		}
		have := enter.DetermineAliasResult(selections, allAliasableCommands, existingAliases)
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
			have := enter.DetermineAliasSelectionText(give)
			want := "(all)"
			must.EqOp(t, want, have)
		})
		t.Run("no commands selected", func(t *testing.T) {
			t.Parallel()
			give := configdomain.AliasableCommands{}
			have := enter.DetermineAliasSelectionText(give)
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
			have := enter.DetermineAliasSelectionText(give)
			want := "append, hack, sync"
			must.EqOp(t, want, have)
		})
	})

	t.Run("Checked", func(t *testing.T) {
		t.Parallel()
		model := enter.AliasesModel{ //nolint:exhaustruct
			AllAliasableCommands: configdomain.AliasableCommands{
				configdomain.AliasableCommandAppend,
				configdomain.AliasableCommandHack,
				configdomain.AliasableCommandShip,
				configdomain.AliasableCommandSync,
			},
			CurrentSelections: []enter.AliasSelection{
				enter.AliasSelectionGT,
				enter.AliasSelectionGT,
				enter.AliasSelectionNone,
				enter.AliasSelectionOther,
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
		have := enter.NewAliasSelections(allAliasableCommands, existingAliases)
		want := []enter.AliasSelection{
			enter.AliasSelectionGT,
			enter.AliasSelectionGT,
			enter.AliasSelectionNone,
			enter.AliasSelectionOther,
		}
		must.Eq(t, want, have)
	})

	t.Run("RotateCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("currently selecting an alias that isn't currently set on disk", func(t *testing.T) {
			t.Parallel()
			model := enter.AliasesModel{
				AllAliasableCommands: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
				},
				CurrentSelections: []enter.AliasSelection{
					enter.AliasSelectionNone,
				},
				OriginalAliases: configdomain.Aliases{},
				BubbleList: dialogcomponents.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
					Cursor: 0,
				},
			}
			// rotate the first time to set to "checked"
			model.RotateCurrentEntry()
			want := []enter.AliasSelection{
				enter.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate a second time to set to "unchecked"
			model.RotateCurrentEntry()
			want = []enter.AliasSelection{
				enter.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
		t.Run("currently selecting an alias that is set on disk to a Git Town command", func(t *testing.T) {
			t.Parallel()
			model := enter.AliasesModel{
				AllAliasableCommands: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
				},
				CurrentSelections: []enter.AliasSelection{
					enter.AliasSelectionGT,
				},
				OriginalAliases: configdomain.Aliases{
					configdomain.AliasableCommandAppend: "town append",
				},
				BubbleList: dialogcomponents.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
					Cursor: 0,
				},
			}
			// rotate the first time to uncheck
			model.RotateCurrentEntry()
			want := []enter.AliasSelection{
				enter.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate the second time to check again
			model.RotateCurrentEntry()
			want = []enter.AliasSelection{
				enter.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
		t.Run("currently selecting an alias that is currently set on disk to an external command", func(t *testing.T) {
			t.Parallel()
			model := enter.AliasesModel{
				AllAliasableCommands: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
				},
				CurrentSelections: []enter.AliasSelection{
					enter.AliasSelectionOther,
				},
				OriginalAliases: configdomain.Aliases{
					configdomain.AliasableCommandAppend: "other command",
				},
				BubbleList: dialogcomponents.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
					Cursor: 0,
				},
			}
			// rotate the first time to check
			model.RotateCurrentEntry()
			want := []enter.AliasSelection{
				enter.AliasSelectionGT,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate the second time to uncheck
			model.RotateCurrentEntry()
			want = []enter.AliasSelection{
				enter.AliasSelectionNone,
			}
			must.Eq(t, want, model.CurrentSelections)
			// rotate a third time to set to "other" again
			model.RotateCurrentEntry()
			want = []enter.AliasSelection{
				enter.AliasSelectionOther,
			}
			must.Eq(t, want, model.CurrentSelections)
		})
	})

	t.Run("SelectAll", func(t *testing.T) {
		t.Parallel()
		model := enter.AliasesModel{ //nolint:exhaustruct
			BubbleList: dialogcomponents.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
				Entries: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
					configdomain.AliasableCommandHack,
					configdomain.AliasableCommandSync,
				},
			},
			CurrentSelections: []enter.AliasSelection{
				enter.AliasSelectionNone,
				enter.AliasSelectionNone,
				enter.AliasSelectionNone,
			},
		}
		model.SelectAll()
		want := []enter.AliasSelection{
			enter.AliasSelectionGT,
			enter.AliasSelectionGT,
			enter.AliasSelectionGT,
		}
		must.Eq(t, model.CurrentSelections, want)
	})

	t.Run("SelectNone", func(t *testing.T) {
		t.Parallel()
		model := enter.AliasesModel{ //nolint:exhaustruct
			BubbleList: dialogcomponents.BubbleList[configdomain.AliasableCommand]{ //nolint:exhaustruct
				Entries: configdomain.AliasableCommands{
					configdomain.AliasableCommandAppend,
					configdomain.AliasableCommandHack,
					configdomain.AliasableCommandDiffParent,
				},
			},
			CurrentSelections: []enter.AliasSelection{
				enter.AliasSelectionGT,
				enter.AliasSelectionOther,
				enter.AliasSelectionGT,
			},
		}
		model.SelectNone()
		want := []enter.AliasSelection{
			enter.AliasSelectionNone,
			enter.AliasSelectionNone,
			enter.AliasSelectionNone,
		}
		must.Eq(t, model.CurrentSelections, want)
	})
}
