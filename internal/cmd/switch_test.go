package cmd_test

import (
	"regexp"
	"testing"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cmd"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/regexes"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()
	alpha := gitdomain.NewLocalBranchName("alpha")
	beta := gitdomain.NewLocalBranchName("beta")
	main := gitdomain.NewLocalBranchName("main")
	prototype := gitdomain.NewLocalBranchName("prototype")
	perennial := gitdomain.NewLocalBranchName("perennial")

	t.Run("SwitchBranchCursorPos", func(t *testing.T) {
		t.Parallel()
		t.Run("initialBranch is in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := dialog.SwitchBranchEntries{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "", OtherWorktree: false},
				{Branch: "alpha1", Indentation: "", OtherWorktree: false},
				{Branch: "beta", Indentation: "", OtherWorktree: false},
			}
			initialBranch := gitdomain.NewLocalBranchName("alpha1")
			have := cmd.SwitchBranchCursorPos(entries, initialBranch)
			want := 2
			must.EqOp(t, want, have)
		})
		t.Run("initialBranch is not in the entry list", func(t *testing.T) {
			t.Parallel()
			entries := dialog.SwitchBranchEntries{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "", OtherWorktree: false},
				{Branch: "beta", Indentation: "", OtherWorktree: false},
			}
			initialBranch := gitdomain.NewLocalBranchName("other")
			have := cmd.SwitchBranchCursorPos(entries, initialBranch)
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
				lineage := configdomain.NewLineageWith(configdomain.LineageData{
					alpha: main,
					beta:  main,
				})
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(alpha), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(beta), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
					{Branch: "beta", Indentation: "  ", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("one of the feature branches is in other worktree", func(t *testing.T) {
				t.Parallel()
				lineage := configdomain.NewLineageWith(configdomain.LineageData{
					alpha: main,
					beta:  main,
				})
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(alpha), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(beta), SyncStatus: gitdomain.SyncStatusOtherWorktree},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "alpha", Indentation: "  ", OtherWorktree: false},
					{Branch: "beta", Indentation: "  ", OtherWorktree: true},
				}
				must.Eq(t, want, have)
			})
		})

		t.Run("perennial branches", func(t *testing.T) {
			t.Parallel()
			perennial1 := gitdomain.NewLocalBranchName("perennial-1")
			lineage := configdomain.NewLineageWith(configdomain.LineageData{
				alpha: main,
				beta:  main,
			})
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{LocalName: Some(alpha), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{LocalName: Some(beta), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{LocalName: Some(perennial1), SyncStatus: gitdomain.SyncStatusLocalOnly},
			}
			branchTypes := []configdomain.BranchType{}
			branchesAndTypes := configdomain.BranchesAndTypes{}
			unknownBranchType := configdomain.BranchTypeFeatureBranch
			regexes := []*regexp.Regexp{}
			have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
			want := dialog.SwitchBranchEntries{
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
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(local), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{RemoteName: Some(remote), SyncStatus: gitdomain.SyncStatusRemoteOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "local", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("when disabled, does not display parent branches of local branches if they are remote only", func(t *testing.T) {
				t.Parallel()
				child := gitdomain.NewLocalBranchName("child")
				grandchild := gitdomain.NewLocalBranchName("grandchild")
				lineage := configdomain.NewLineageWith(configdomain.LineageData{
					child:      main,
					grandchild: child,
				})
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{RemoteName: Some(child.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusRemoteOnly},
					gitdomain.BranchInfo{LocalName: Some(grandchild), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: None[gitdomain.LocalBranchName](), RemoteName: Some(child.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(grandchild), RemoteName: Some(grandchild.AtRemote(gitdomain.RemoteOrigin)), SyncStatus: gitdomain.SyncStatusUpToDate},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "grandchild", Indentation: "    ", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("enabled", func(t *testing.T) {
				t.Parallel()
				local := gitdomain.NewLocalBranchName("local")
				remote := gitdomain.NewRemoteBranchName("origin/remote")
				lineage := configdomain.NewLineageWith(configdomain.LineageData{
					local: main,
				})
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{RemoteName: Some(remote), SyncStatus: gitdomain.SyncStatusRemoteOnly},
					gitdomain.BranchInfo{LocalName: Some(local), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, true, regexes)
				want := dialog.SwitchBranchEntries{
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
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "observed-1", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeObservedBranch},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeObservedBranch},
				}
				must.Eq(t, want, have)
			})
			t.Run("multiple branch types", func(t *testing.T) {
				t.Parallel()
				observed1 := gitdomain.NewLocalBranchName("observed-1")
				observed2 := gitdomain.NewLocalBranchName("observed-2")
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
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes := []*regexp.Regexp{}
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "observed-1", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeObservedBranch},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeObservedBranch},
					{Branch: "perennial", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypePerennialBranch},
				}
				must.Eq(t, want, have)
			})
		})

		t.Run("filter by regexes", func(t *testing.T) {
			t.Parallel()
			t.Run("no regex", func(t *testing.T) {
				t.Parallel()
				observed1 := gitdomain.NewLocalBranchName("observed-1")
				observed2 := gitdomain.NewLocalBranchName("observed-2")
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed1), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed2), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(perennial), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(prototype), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes, err := regexes.NewRegexes([]string{})
				must.NoError(t, err)
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "observed-1", Indentation: "", OtherWorktree: false},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false},
					{Branch: "perennial", Indentation: "", OtherWorktree: false},
					{Branch: "prototype", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("single regex", func(t *testing.T) {
				t.Parallel()
				observed1 := gitdomain.NewLocalBranchName("observed-1")
				observed2 := gitdomain.NewLocalBranchName("observed-2")
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(observed1), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed2), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(prototype), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(perennial), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes, err := regexes.NewRegexes([]string{"observed-"})
				must.NoError(t, err)
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "observed-1", Indentation: "", OtherWorktree: false},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
			t.Run("multiple regexes", func(t *testing.T) {
				t.Parallel()
				observed1 := gitdomain.NewLocalBranchName("observed-1")
				observed2 := gitdomain.NewLocalBranchName("observed-2")
				lineage := configdomain.NewLineage()
				branchInfos := gitdomain.BranchInfos{
					gitdomain.BranchInfo{LocalName: Some(main), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed1), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(observed2), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(perennial), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{LocalName: Some(prototype), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.BranchTypeFeatureBranch
				regexes, err := regexes.NewRegexes([]string{"observed-", "main"})
				must.NoError(t, err)
				have := cmd.SwitchBranchEntries(branchInfos, branchTypes, branchesAndTypes, lineage, unknownBranchType, false, regexes)
				want := dialog.SwitchBranchEntries{
					{Branch: "main", Indentation: "", OtherWorktree: false},
					{Branch: "observed-1", Indentation: "", OtherWorktree: false},
					{Branch: "observed-2", Indentation: "", OtherWorktree: false},
				}
				must.Eq(t, want, have)
			})
		})
	})
}
