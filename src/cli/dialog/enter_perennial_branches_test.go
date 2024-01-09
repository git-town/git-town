package dialog //nolint:testpackage

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/shoenig/test/must"
)

func TestPerennialBranchesModel(t *testing.T) {
	t.Parallel()
	t.Run("disableCurrentEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("entry is enabled", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections:    []int{1, 2, 3},
				selectedColor: termenv.Style{},
			}
			model.disableCurrentEntry()
			wantSelections := []int{1, 3}
			must.Eq(t, wantSelections, model.selections)
		})
		t.Run("entry is disabled", func(t *testing.T) {
			t.Parallel()
			model := perennialBranchesModel{
				bubbleList: bubbleList{ //nolint:exhaustruct
					cursor: 2,
				},
				selections:    []int{1, 3},
				selectedColor: termenv.Style{},
			}
			model.disableCurrentEntry()
			wantSelections := []int{1, 3}
			must.Eq(t, wantSelections, model.selections)
		})
	})
}
