package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()

	t.Run("SwitchBranchCursorPos", func(t *testing.T) {
		t.Parallel()
		t.Run("initialBranch is in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "", OtherWorktree: false},
				{Branch: "alpha1", Indentation: "", OtherWorktree: false},
				{Branch: "beta", Indentation: "", OtherWorktree: false},
			}
			initialBranch := gitdomain.NewLocalBranchName("alpha1")
			have := dialog.SwitchBranchCursorPos(entries, initialBranch)
			want := 2
			must.EqOp(t, want, have)
		})
		t.Run("initialBranch is not in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "", OtherWorktree: false},
				{Branch: "beta", Indentation: "", OtherWorktree: false},
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
			alpha := gitdomain.NewLocalBranchName("alpha")
			beta := gitdomain.NewLocalBranchName("beta")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				alpha: main,
				beta:  main,
			}
			localBranches := gitdomain.LocalBranchNames{alpha, beta, main}
			allBranches := gitdomain.BranchInfos{
				gitdomain.BranchInfo{LocalName: alpha, SyncStatus: gitdomain.SyncStatusLocalOnly}, //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: beta, SyncStatus: gitdomain.SyncStatusLocalOnly},  //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: main, SyncStatus: gitdomain.SyncStatusLocalOnly},  //nolint:exhaustruct
			}
			have := dialog.SwitchBranchEntries(localBranches, lineage, allBranches)
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
				{Branch: "beta", Indentation: "  ", OtherWorktree: false},
			}
			must.Eq(t, want, have)
		})
		t.Run("feature branch in other worktree", func(t *testing.T) {
			t.Parallel()
			alpha := gitdomain.NewLocalBranchName("alpha")
			beta := gitdomain.NewLocalBranchName("beta")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				alpha: main,
				beta:  main,
			}
			localBranches := gitdomain.LocalBranchNames{alpha, beta, main}
			allBranches := gitdomain.BranchInfos{
				gitdomain.BranchInfo{LocalName: alpha, SyncStatus: gitdomain.SyncStatusLocalOnly},    //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: beta, SyncStatus: gitdomain.SyncStatusOtherWorktree}, //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: main, SyncStatus: gitdomain.SyncStatusLocalOnly},     //nolint:exhaustruct
			}
			have := dialog.SwitchBranchEntries(localBranches, lineage, allBranches)
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
				{Branch: "beta", Indentation: "  ", OtherWorktree: true},
			}
			must.Eq(t, want, have)
		})
		t.Run("feature and perennial branches", func(t *testing.T) {
			t.Parallel()
			alpha := gitdomain.NewLocalBranchName("alpha")
			beta := gitdomain.NewLocalBranchName("beta")
			perennial1 := gitdomain.NewLocalBranchName("perennial-1")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				alpha: main,
				beta:  main,
			}
			localBranches := gitdomain.LocalBranchNames{alpha, beta, main, perennial1}
			allBranches := gitdomain.BranchInfos{
				gitdomain.BranchInfo{LocalName: alpha, SyncStatus: gitdomain.SyncStatusLocalOnly},      //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: beta, SyncStatus: gitdomain.SyncStatusLocalOnly},       //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: main, SyncStatus: gitdomain.SyncStatusLocalOnly},       //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: perennial1, SyncStatus: gitdomain.SyncStatusLocalOnly}, //nolint:exhaustruct
			}
			have := dialog.SwitchBranchEntries(localBranches, lineage, allBranches)
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
				{Branch: "beta", Indentation: "  ", OtherWorktree: false},
				{Branch: "perennial-1", Indentation: "", OtherWorktree: false},
			}
			must.Eq(t, want, have)
		})
		t.Run("parent exists remotely but is not checked out locally", func(t *testing.T) {
			t.Parallel()
			child := gitdomain.NewLocalBranchName("child")
			grandchild := gitdomain.NewLocalBranchName("grandchild")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.Lineage{
				child:      main,
				grandchild: child,
			}
			localBranches := gitdomain.LocalBranchNames{grandchild, main}
			allBranches := gitdomain.BranchInfos{
				gitdomain.BranchInfo{LocalName: gitdomain.EmptyLocalBranchName(), RemoteName: child.BranchName().RemoteName(), SyncStatus: gitdomain.SyncStatusRemoteOnly}, //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: grandchild, SyncStatus: gitdomain.SyncStatusLocalOnly},                                                                     //nolint:exhaustruct
				gitdomain.BranchInfo{LocalName: main, SyncStatus: gitdomain.SyncStatusLocalOnly},                                                                           //nolint:exhaustruct
			}
			have := dialog.SwitchBranchEntries(localBranches, lineage, allBranches)
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "child", Indentation: "  ", OtherWorktree: false},
				{Branch: "grandchild", Indentation: "    ", OtherWorktree: false},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("View", func(t *testing.T) {
		t.Run("only the main branch exists", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{ //nolint:exhaustruct
					Cursor:       0,
					Entries:      newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{{Branch: "main", Indentation: "", OtherWorktree: false}}),
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
				List: list.List[dialog.SwitchBranchEntry]{ //nolint:exhaustruct
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries([]dialog.SwitchBranchEntry{
						{Branch: "main", Indentation: "", OtherWorktree: false},
						{Branch: "one", Indentation: "", OtherWorktree: false},
						{Branch: "two", Indentation: "", OtherWorktree: true},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos: 0,
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
				List: list.List[dialog.SwitchBranchEntry]{ //nolint:exhaustruct
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
				InitialBranchPos: 0,
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
	})
}

func newSwitchBranchBubbleListEntries(entries []dialog.SwitchBranchEntry) []list.Entry[dialog.SwitchBranchEntry] {
	result := make([]list.Entry[dialog.SwitchBranchEntry], len(entries))
	for e, entry := range entries {
		result[e] = list.Entry[dialog.SwitchBranchEntry]{
			Data: entry,
			Text: entry.String(),
		}
	}
	return result
}
