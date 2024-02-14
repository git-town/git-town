package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()

	t.Run("SwitchBranchCursorPos", func(t *testing.T) {
		t.Parallel()
		t.Run("initialBranch is in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: ""},
				{Branch: "alpha", Indentation: ""},
				{Branch: "alpha1", Indentation: ""},
				{Branch: "beta", Indentation: ""},
			}
			initialBranch := gitdomain.NewLocalBranchName("alpha1")
			have := dialog.SwitchBranchCursorPos(entries, initialBranch)
			want := 2
			must.EqOp(t, want, have)
		})
		t.Run("initialBranch is not in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: ""},
				{Branch: "alpha", Indentation: ""},
				{Branch: "beta", Indentation: ""},
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
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: ""},
				{Branch: "alpha", Indentation: "  "},
				{Branch: "beta", Indentation: "  "},
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
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: ""},
				{Branch: "alpha", Indentation: "  "},
				{Branch: "beta", Indentation: "  "},
				{Branch: "perennial-1", Indentation: ""},
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
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: ""},
				{Branch: "child", Indentation: "  "},
				{Branch: "grandchild", Indentation: "    "},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("View", func(t *testing.T) {
		t.Run("only the main branch exists", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				BubbleList: components.BubbleList[dialog.SwitchBranchEntry]{ //nolint:exhaustruct
					Cursor:       0,
					Entries:      []dialog.SwitchBranchEntry{{Branch: "main", Indentation: ""}},
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos: 0,
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
				BubbleList: components.BubbleList[dialog.SwitchBranchEntry]{ //nolint:exhaustruct
					Cursor: 0,
					Entries: []dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: ""},
						{Branch: "one", Indentation: ""},
						{Branch: "two", Indentation: ""},
					},
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


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("stacked changes", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				BubbleList: components.BubbleList[dialog.SwitchBranchEntry]{ //nolint:exhaustruct
					Cursor: 0,
					Entries: []dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: ""},
						{Branch: "alpha", Indentation: "  "},
						{Branch: "alpha1", Indentation: "    "},
						{Branch: "alpha2", Indentation: "    "},
						{Branch: "beta", Indentation: "  "},
						{Branch: "beta1", Indentation: "    "},
						{Branch: "other", Indentation: ""},
					},
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


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})
	})
}
