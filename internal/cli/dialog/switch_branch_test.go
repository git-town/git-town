package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
		t.Run("worktree", func(t *testing.T) {
			t.Parallel()
			t.Run("all branches are in the current worktree", func(t *testing.T) {
				t.Parallel()
				alpha := gitdomain.NewLocalBranchName("alpha")
				beta := gitdomain.NewLocalBranchName("beta")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				lineage.Add(alpha, main)
				lineage.Add(beta, main)
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(alpha), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(beta), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
				want := []dialog.SwitchBranchEntry{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
					{Branch: "beta", Indentation: "  ", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("one of the feature branches is in other worktree", func(t *testing.T) {
				t.Parallel()
				alpha := gitdomain.NewLocalBranchName("alpha")
				beta := gitdomain.NewLocalBranchName("beta")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				lineage.Add(alpha, main)
				lineage.Add(beta, main)
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(alpha), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(beta), SyncStatus: gitdomain.SyncStatusOtherWorktree},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
				want := []dialog.SwitchBranchEntry{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
					{Branch: "beta", Indentation: "  ", OtherWorktree: true},
				}
				must.Eq(t, want, have)
			})
		})
		t.Run("perennial branches", func(t *testing.T) {
			t.Parallel()
			alpha := gitdomain.NewLocalBranchName("alpha")
			beta := gitdomain.NewLocalBranchName("beta")
			perennial1 := gitdomain.NewLocalBranchName("perennial-1")
			main := gitdomain.NewLocalBranchName("main")
			lineage := configdomain.NewLineage()
			lineage.Add(alpha, main)
			lineage.Add(beta, main)
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{LocalName: Some(alpha), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{LocalName: Some(beta), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{LocalName: Some(perennial1), SyncStatus: gitdomain.SyncStatusLocalOnly},
			}
			branchTypes := []configdomain.BranchType{}
			branchesAndTypes := configdomain.BranchesAndTypes{}
			defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
			have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
			want := []dialog.SwitchBranchEntry{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
				{Branch: "beta", Indentation: "  ", OtherWorktree: false},
				{Branch: "perennial-1", Indentation: "", OtherWorktree: false},
			}
			must.Eq(t, want, have)
		})
		t.Run("--all flag", func(t *testing.T) {
			t.Parallel()
			t.Run("disabled", func(t *testing.T) {
				t.Parallel()
				local := gitdomain.NewLocalBranchName("local")
				remote := gitdomain.NewRemoteBranchName("origin/remote")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(local), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{RemoteName: Some(remote), SyncStatus: gitdomain.SyncStatusRemoteOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
				want := []dialog.SwitchBranchEntry{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "local", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("when disabled, does not display parent branches of local branches if they are remote only", func(t *testing.T) {
				t.Parallel()
				child := gitdomain.NewLocalBranchName("child")
				grandchild := gitdomain.NewLocalBranchName("grandchild")
				local := gitdomain.NewLocalBranchName("local")
				remote := gitdomain.NewLocalBranchName("remote")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				lineage.Add(local, main)
				lineage.Add(child, main)
				lineage.Add(grandchild, child)
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: None[gitdomain.LocalBranchName](), RemoteName: Some(remote.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusRemoteOnly},
					gitdomain.BranchInfo{LocalName: Some(local), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{RemoteName: Some(child.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusRemoteOnly},
					gitdomain.BranchInfo{LocalName: Some(grandchild), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: None[gitdomain.LocalBranchName](), RemoteName: Some(child.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(grandchild), RemoteName: Some(grandchild.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusUpToDate},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
				want := []dialog.SwitchBranchEntry{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "grandchild", Indentation: "    ", OtherWorktree: false},
					{Branch: "local", Indentation: "  ", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("enabled", func(t *testing.T) {
				t.Parallel()
				local := gitdomain.NewLocalBranchName("local")
				remote := gitdomain.NewLocalBranchName("remote")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				lineage.Add(local, main)
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: None[gitdomain.LocalBranchName](), RemoteName: Some(remote.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusRemoteOnly},
					gitdomain.BranchInfo{LocalName: Some(local), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, true)
				want := []dialog.SwitchBranchEntry{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "local", Indentation: "  ", OtherWorktree: false},
					{Branch: "remote", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
		})
		t.Run("filter by branch type", func(t *testing.T) {
			t.Parallel()
			t.Run("single branch type", func(t *testing.T) {
				t.Parallel()
				observed1 := gitdomain.NewLocalBranchName("observed-1")
				observed2 := gitdomain.NewLocalBranchName("observed-2")
				prototype := gitdomain.NewLocalBranchName("prototype")
				perennial := gitdomain.NewLocalBranchName("perennial")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(observed1), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed2), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(prototype), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(perennial), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{configdomain.BranchTypeObservedBranch}
				branchesAndTypes := configdomain.BranchesAndTypes{
					observed1: configdomain.BranchTypeObservedBranch,
					observed2: configdomain.BranchTypeObservedBranch,
					prototype: configdomain.BranchTypePrototypeBranch,
					perennial: configdomain.BranchTypePerennialBranch,
					main:      configdomain.BranchTypeMainBranch,
				}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
				want := []dialog.SwitchBranchEntry{
					{Branch: "observed-1", Indentation: "", OtherWorktree: false},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("multiple branch types", func(t *testing.T) {
				t.Parallel()
				observed1 := gitdomain.NewLocalBranchName("observed-1")
				observed2 := gitdomain.NewLocalBranchName("observed-2")
				prototype := gitdomain.NewLocalBranchName("prototype")
				perennial := gitdomain.NewLocalBranchName("perennial")
				main := gitdomain.NewLocalBranchName("main")
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(observed1), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed2), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(prototype), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(perennial), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{
					configdomain.BranchTypeObservedBranch,
					configdomain.BranchTypePerennialBranch,
				}
				branchesAndTypes := configdomain.BranchesAndTypes{
					observed1: configdomain.BranchTypeObservedBranch,
					observed2: configdomain.BranchTypeObservedBranch,
					prototype: configdomain.BranchTypePrototypeBranch,
					perennial: configdomain.BranchTypePerennialBranch,
					main:      configdomain.BranchTypeMainBranch,
				}
				defaultBranchType := configdomain.DefaultBranchType{BranchType: configdomain.BranchTypeFeatureBranch}
				have := dialog.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, defaultBranchType, false)
				want := []dialog.SwitchBranchEntry{
					{Branch: "observed-1", Indentation: "", OtherWorktree: false},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false},
					{Branch: "perennial", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
		})
	})

	t.Run("View", func(t *testing.T) {
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
