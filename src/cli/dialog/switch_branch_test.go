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
			InitialBranch: "main",
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
			InitialBranch: "main",
		}
		have := model.View()
		want := `
> main
  one
  two


  ↑/k up   ↓/j down   enter/o accept   esc/q abort`[1:]
		must.EqOp(t, want, have)
	})
}
