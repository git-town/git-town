package components_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestPerennialBranches(t *testing.T) {
	t.Parallel()

	t.Run("disableCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("entry is enabled", func(t *testing.T) {
			t.Parallel()
			model := components.CheckListModel[gitdomain.LocalBranchName]{
				List: list.List[gitdomain.LocalBranchName]{
					Cursor: 2,
				},
				Selections: []int{1, 2, 3},
			}
			model = model.DisableCurrentEntry()
			wantSelections := []int{1, 3}
			must.Eq(t, wantSelections, model.Selections)
		})
		t.Run("entry is disabled", func(t *testing.T) {
			t.Parallel()
			model := components.CheckListModel[gitdomain.LocalBranchName]{
				List: list.List[gitdomain.LocalBranchName]{
					Cursor: 2,
				},
				Selections: []int{1, 3},
			}
			model = model.DisableCurrentEntry()
			wantSelections := []int{1, 3}
			must.Eq(t, wantSelections, model.Selections)
		})
	})

	t.Run("enableCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("entry is disabled", func(t *testing.T) {
			t.Parallel()
			model := components.CheckListModel[gitdomain.LocalBranchName]{
				List: list.List[gitdomain.LocalBranchName]{
					Cursor: 2,
				},
				Selections: []int{1, 3},
			}
			model = model.EnableCurrentEntry()
			wantSelections := []int{1, 3, 2}
			must.Eq(t, wantSelections, model.Selections)
		})
		t.Run("entry is enabled", func(t *testing.T) {
			t.Parallel()
			model := components.CheckListModel[gitdomain.LocalBranchName]{
				List: list.List[gitdomain.LocalBranchName]{
					Cursor: 2,
				},
				Selections: []int{1, 2, 3},
			}
			model = model.EnableCurrentEntry()
			wantSelections := []int{1, 2, 3}
			must.Eq(t, wantSelections, model.Selections)
		})
	})

	t.Run("isRowChecked", func(t *testing.T) {
		t.Parallel()
		model := components.CheckListModel[gitdomain.LocalBranchName]{
			Selections: []int{2},
		}
		must.False(t, model.IsRowChecked(1))
		must.True(t, model.IsRowChecked(2))
		must.False(t, model.IsRowChecked(3))
	})

	t.Run("checkedEntries", func(t *testing.T) {
		t.Parallel()
		model := components.CheckListModel[gitdomain.LocalBranchName]{
			List: list.List[gitdomain.LocalBranchName]{
				Entries: list.NewEntries[gitdomain.LocalBranchName]("zero", "one", "two", "three"),
			},
			Selections: []int{1, 3},
		}
		have := model.CheckedEntries()
		want := gitdomain.NewLocalBranchNames("one", "three")
		must.Eq(t, want, have)
	})

	t.Run("toggleCurrentEntry", func(t *testing.T) {
		t.Parallel()
		model := components.CheckListModel[gitdomain.LocalBranchName]{
			List: list.List[gitdomain.LocalBranchName]{
				Cursor: 2,
			},
			Selections: []int{1, 3},
		}
		// enable the selected entry
		model = model.ToggleCurrentEntry()
		wantSelections := []int{1, 3, 2}
		must.Eq(t, wantSelections, model.Selections)
		// disable the selected entry
		model = model.ToggleCurrentEntry()
		wantSelections = []int{1, 3}
		must.Eq(t, wantSelections, model.Selections)
	})
}
