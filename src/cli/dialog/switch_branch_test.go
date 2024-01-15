package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()
	t.Run("only the main branch exists", func(t *testing.T) {
		t.Parallel()
		model := dialog.SwitchModel{
			BubbleList: dialog.BubbleList{ //nolint:exhaustruct
				Cursor:       0,
				Entries:      []string{"main"},
				MaxDigits:    1,
				NumberFormat: "%d",
			},
			InitialBranchPos: 0,
		}
		have := model.View()
		want := `
> main


  ↑/k up   ↓/j down   enter/o accept   esc/q abort`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("multiple top-level branches", func(t *testing.T) {
		t.Parallel()
		model := dialog.SwitchModel{
			BubbleList: dialog.BubbleList{ //nolint:exhaustruct
				Cursor:       0,
				Entries:      []string{"main", "one", "two"},
				MaxDigits:    1,
				NumberFormat: "%d",
			},
			InitialBranchPos: 0,
		}
		have := model.View()
		want := `
> main
  one
  two


  ↑/k up   ↓/j down   enter/o accept   esc/q abort`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("nested branches", func(t *testing.T) {
		t.Parallel()
		model := dialog.SwitchModel{
			BubbleList: dialog.BubbleList{ //nolint:exhaustruct
				Cursor:       0,
				Entries:      []string{"main", "  alpha", "    alpha1", "    alpha2", "  beta", "    beta1", "other"},
				MaxDigits:    1,
				NumberFormat: "%d",
			},
			InitialBranchPos: 0,
		}
		have := model.View()
		want := `
> main
    alpha
      alpha1
      alpha2
    beta
      beta1
  other


  ↑/k up   ↓/j down   enter/o accept   esc/q abort`[1:]
		must.EqOp(t, want, have)
	})
}
