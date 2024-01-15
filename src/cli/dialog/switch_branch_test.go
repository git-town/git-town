package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()

	t.Run("SwitchBranchCursorPos", func(t *testing.T) {
		t.Parallel()
		t.Run("initialBranch is in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := []string{
				"main",
				"  alpha",
				"    alpha1",
				"  beta",
			}
			initialBranch := gitdomain.NewLocalBranchName("alpha1")
			have := dialog.SwitchBranchCursorPos(entries, initialBranch)
			want := 2
			must.EqOp(t, want, have)
		})
	})

	t.Run("SwitchBranchEntries", func(t *testing.T) {
		t.Parallel()
		t.Run("normal", func(t *testing.T) {
			t.Parallel()
			branchA := gitdomain.NewLocalBranchName("alpha")
			branchB := gitdomain.NewLocalBranchName("beta")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				branchA: main,
				branchB: main,
			}
			localBranches := gitdomain.LocalBranchNames{branchA, branchB, main}
			have := dialog.SwitchBranchEntries(localBranches, lineage)
			want := []string{
				"main",
				"  alpha",
				"  beta",
			}
			must.Eq(t, want, have)
		})
		t.Run("with perennial branches", func(t *testing.T) {
			t.Parallel()
			branchA := gitdomain.NewLocalBranchName("alpha")
			branchB := gitdomain.NewLocalBranchName("beta")
			perennial1 := gitdomain.NewLocalBranchName("perennial-1")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				branchA: main,
				branchB: main,
			}
			localBranches := gitdomain.LocalBranchNames{branchA, branchB, main, perennial1}
			have := dialog.SwitchBranchEntries(localBranches, lineage)
			want := []string{
				"main",
				"  alpha",
				"  beta",
				"perennial-1",
			}
			must.Eq(t, want, have)
		})
		t.Run("local grandparent branch with missing parent", func(t *testing.T) {
			t.Parallel()
			branchA := gitdomain.NewLocalBranchName("alpha")
			branchA1 := gitdomain.NewLocalBranchName("alpha-1")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				branchA:  main,
				branchA1: branchA,
			}
			localBranches := gitdomain.LocalBranchNames{branchA1, main}
			have := dialog.SwitchBranchEntries(localBranches, lineage)
			want := []string{
				"main",
				"  alpha",
				"    alpha-1",
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("View", func(t *testing.T) {
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
	})
}
