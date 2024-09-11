package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
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
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
				DisplayBranchTypes: false,
			}
			have := model.View()
			want := `
> main


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("multiple top-level branches", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeMainBranch},
						{Branch: "one", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "two", Indentation: "", OtherWorktree: true, Type: configdomain.BranchTypeFeatureBranch},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
				DisplayBranchTypes: false,
			}
			have := model.View()
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
						{Branch: "main", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeMainBranch},
						{Branch: "alpha", Indentation: "  ", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "alpha1", Indentation: "    ", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "alpha2", Indentation: "    ", OtherWorktree: true, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "beta", Indentation: "  ", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "beta1", Indentation: "    ", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "other", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
				DisplayBranchTypes: false,
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

		t.Run("stacked changes with types", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeMainBranch},
						{Branch: "alpha", Indentation: "  ", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "alpha1", Indentation: "    ", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "alpha2", Indentation: "    ", OtherWorktree: true, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "beta", Indentation: "  ", OtherWorktree: false, Type: configdomain.BranchTypeObservedBranch},
						{Branch: "beta1", Indentation: "    ", OtherWorktree: false, Type: configdomain.BranchTypeObservedBranch},
						{Branch: "other", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeParkedBranch},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: false,
				DisplayBranchTypes: true,
			}
			have := model.View()
			dim := "\x1b[2m"
			reset := "\x1b[0m"
			want := `
> main  ` + dim + `(main)` + reset + `
    alpha  ` + dim + `(feature)` + reset + `
      alpha1  ` + dim + `(feature)` + reset + `
` + dim + `+     alpha2` + reset + `  ` + dim + `(feature)` + reset + `
    beta  ` + dim + `(observed)` + reset + `
      beta1  ` + dim + `(observed)` + reset + `
  other  ` + dim + `(parked)` + reset + `


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("uncommitted changes", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeMainBranch},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   0,
				UncommittedChanges: true,
				DisplayBranchTypes: false,
			}
			have := model.View()
			cyanBold := "\x1b[36;1m"
			reset := "\x1b[0m"
			want := `
` + cyanBold + `uncommitted changes` + reset + `

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
