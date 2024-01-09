package dialog //nolint:testpackage

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestPerennialBranchesModel(t *testing.T) {
	t.Parallel()

	t.Run("disableCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("entry is enabled", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{ //nolint:exhaustruct
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections: []int{1, 2, 3},
			}
			model.disableCurrentEntry()
			wantSelections := []int{1, 3}
			must.Eq(t, wantSelections, model.selections)
		})
		t.Run("entry is disabled", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{ //nolint:exhaustruct
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections: []int{1, 3},
			}
			model.disableCurrentEntry()
			wantSelections := []int{1, 3}
			must.Eq(t, wantSelections, model.selections)
		})
	})

	t.Run("enableCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("entry is disabled", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{ //nolint:exhaustruct
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections: []int{1, 3},
			}
			model.enableCurrentEntry()
			wantSelections := []int{1, 3, 2}
			must.Eq(t, wantSelections, model.selections)
		})
		t.Run("entry is enabled", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{ //nolint:exhaustruct
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections: []int{1, 2, 3},
			}
			model.enableCurrentEntry()
			wantSelections := []int{1, 2, 3}
			must.Eq(t, wantSelections, model.selections)
		})
	})

	t.Run("isSelectedRowChecked", func(t *testing.T) {
		t.Parallel()
		t.Run("selected row is checked", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{ //nolint:exhaustruct
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections: []int{2},
			}
			must.True(t, model.isSelectedRowChecked())
		})
		t.Run("selected row is not checked", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{ //nolint:exhaustruct
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 1,
				},
				selections: []int{2},
			}
			must.False(t, model.isSelectedRowChecked())
		})
	})

	t.Run("isRowChecked", func(t *testing.T) {
		t.Parallel()
		model := perennialBranchesModel{ //nolint:exhaustruct
			selections: []int{2},
		}
		must.False(t, model.isRowChecked(1))
		must.True(t, model.isRowChecked(2))
		must.False(t, model.isRowChecked(3))
	})

	t.Run("checkedEntries", func(t *testing.T) {
		t.Parallel()
		model := perennialBranchesModel{ //nolint:exhaustruct
			bubbleList: bubbleList{ //nolint:exhaustruct
				entries: []string{"zero", "one", "two", "three"},
			},
			selections: []int{1, 3},
		}
		have := model.checkedEntries()
		want := []string{"one", "three"}
		must.Eq(t, want, have)
	})

	t.Run("toggleCurrentEntry", func(t *testing.T) {
		t.Parallel()
		model := perennialBranchesModel{ //nolint:exhaustruct
			bubbleList: bubbleList{ //nolint:exhaustruct
				cursor: 2,
			},
			selections: []int{1, 3},
		}
		// enable the selected entry
		model.toggleCurrentEntry()
		wantSelections := []int{1, 3, 2}
		must.Eq(t, wantSelections, model.selections)
		// disable the selected entry
		model.toggleCurrentEntry()
		wantSelections = []int{1, 3}
		must.Eq(t, wantSelections, model.selections)
	})
}
