package dialog_test

import (
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()

	t.Run("View", func(t *testing.T) {
		t.Parallel()
		t.Run("only the main branch exists", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor:       0,
					Entries:      newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{{Branch: "main", Indentation: "", OtherWorktree: false}}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
			}
			have := model.View()
			dim := "\x1b[2m"
			reset := "\x1b[0m"
			want := `
> main  ` + dim + `(main)` + reset + `



  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`[1:]
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("multiple top-level branches", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false},
						{Branch: "one", Indentation: "", OtherWorktree: false},
						{Branch: "two", Indentation: "", OtherWorktree: true},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
			}
			have := stripansi.Strip(model.View())
			dim := "\x1b[2m"
			reset := "\x1b[0m"
			want := `
> main
  one
` + dim + `+ two` + reset + `


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("stacked changes", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false},
						{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
						{Branch: "alpha1", Indentation: "    ", OtherWorktree: false},
						{Branch: "alpha2", Indentation: "    ", OtherWorktree: true},
						{Branch: "beta", Indentation: "  ", OtherWorktree: false},
						{Branch: "beta1", Indentation: "    ", OtherWorktree: false},
						{Branch: "other", Indentation: "", OtherWorktree: false},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
			}
			have := model.View()
			dim := "\x1b[2m"
			reset := "\x1b[0m"
			want := `
> main
    alpha
      alpha1
` + dim + `+     alpha2` + reset + `
    beta
      beta1
  other


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("uncommitted changes", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor:       0,
					Entries:      newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{{Branch: "main", Indentation: "", OtherWorktree: false}}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: true,
			}
			have := model.View()
			want := `
` +
				"\x1b[36;1m" +
				`uncommitted changes
` +
				"\x1b[0m" +
				`

> main


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})
	})
}

func newSwitchBranchBubbleListEntries(entries []dialog.SwitchBranchEntry) []list.Entry[dialog.SwitchBranchEntry] {
	result := make([]list.Entry[dialog.SwitchBranchEntry], len(entries))
	for e, entry := range entries {
		result[e] = list.Entry[dialog.SwitchBranchEntry]{
			Data:    entry,
			Enabled: !entry.OtherWorktree,
			Text:    entry.String(),
		}
	}
	return result
}
