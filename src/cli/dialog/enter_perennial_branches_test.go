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
}
