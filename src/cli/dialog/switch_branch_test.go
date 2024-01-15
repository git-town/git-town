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
		t.Run("initialBranch is not in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := []string{
				"main",
				"  alpha",
				"  beta",
			}
			initialBranch := gitdomain.NewLocalBranchName("other")
			have := dialog.SwitchBranchCursorPos(entries, initialBranch)
			want := 0
			must.EqOp(t, want, have)
		})
	})

	t.Run("SwitchBranchEntries", func(t *testing.T) {
		t.Parallel()
		t.Run("feature branches only", func(t *testing.T) {
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
		t.Run("feature and perennial branches", func(t *testing.T) {
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
		t.Run("parent is not checked out locally", func(t *testing.T) {
			t.Parallel()
			child := gitdomain.NewLocalBranchName("child")
			grandchild := gitdomain.NewLocalBranchName("grandchild")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				child:      main,
				grandchild: child,
			}
			localBranches := gitdomain.LocalBranchNames{grandchild, main}
			have := dialog.SwitchBranchEntries(localBranches, lineage)
			want := []string{
				"main",
				"  child",
				"    grandchild",
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


  ↑/k up   ↓/j down   enter/o accept   q/esc/ctrl-c abort`[1:]
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


  ↑/k up   ↓/j down   enter/o accept   q/esc/ctrl-c abort`[1:]
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


  ↑/k up   ↓/j down   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})
	})
}
