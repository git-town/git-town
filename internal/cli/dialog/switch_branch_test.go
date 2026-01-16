package dialog_test

import (
	"regexp"
	"testing"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/regexes"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestNewSwitchBranch(t *testing.T) {
	t.Parallel()
	alpha := gitdomain.NewLocalBranchName("alpha")
	beta := gitdomain.NewLocalBranchName("beta")
	main := gitdomain.NewLocalBranchName("main")
	prototype := gitdomain.NewLocalBranchName("prototype")
	perennial := gitdomain.NewLocalBranchName("perennial")

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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: alpha}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: beta}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: alpha}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: beta}), SyncStatus: gitdomain.SyncStatusOtherWorktree},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
				gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: alpha}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: beta}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: perennial1}), SyncStatus: gitdomain.SyncStatusLocalOnly},
			}
			branchTypes := []configdomain.BranchType{}
			branchesAndTypes := configdomain.BranchesAndTypes{}
			unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
			regexes := []*regexp.Regexp{}
			have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
				BranchInfos:       branchInfos,
				BranchTypes:       branchTypes,
				BranchesAndTypes:  branchesAndTypes,
				Lineage:           lineage,
				Regexes:           regexes,
				ShowAllBranches:   false,
				UnknownBranchType: unknownBranchType,
			})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: local}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{RemoteName: Some(remote), SyncStatus: gitdomain.SyncStatusRemoteOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{RemoteName: Some(gitdomain.NewRemoteBranchName("origin/child")), SyncStatus: gitdomain.SyncStatusRemoteOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: grandchild}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: None[gitdomain.BranchData](), RemoteName: Some(gitdomain.NewRemoteBranchName("origin/child")), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: grandchild}), RemoteName: Some(gitdomain.NewRemoteBranchName("origin/grandchild")), SyncStatus: gitdomain.SyncStatusUpToDate},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: local}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   true,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed1}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed2}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: prototype}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: perennial}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{configdomain.BranchTypeObservedBranch}
				branchesAndTypes := configdomain.BranchesAndTypes{
					observed1: configdomain.BranchTypeObservedBranch,
					observed2: configdomain.BranchTypeObservedBranch,
					prototype: configdomain.BranchTypePrototypeBranch,
					perennial: configdomain.BranchTypePerennialBranch,
					main:      configdomain.BranchTypeMainBranch,
				}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed1}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed2}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: prototype}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: perennial}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
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
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes := []*regexp.Regexp{}
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed1}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed2}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: perennial}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: prototype}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes, err := regexes.NewRegexes([]string{})
				must.NoError(t, err)
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed1}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed2}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: prototype}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: perennial}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes, err := regexes.NewRegexes([]string{"observed-"})
				must.NoError(t, err)
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: main}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed1}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: observed2}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: perennial}), SyncStatus: gitdomain.SyncStatusLocalOnly},
					gitdomain.BranchInfo{Local: Some(gitdomain.BranchData{Name: prototype}), SyncStatus: gitdomain.SyncStatusLocalOnly},
				}
				branchTypes := []configdomain.BranchType{}
				branchesAndTypes := configdomain.BranchesAndTypes{}
				unknownBranchType := configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)
				regexes, err := regexes.NewRegexes([]string{"observed-", "main"})
				must.NoError(t, err)
				have := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
					BranchInfos:       branchInfos,
					BranchTypes:       branchTypes,
					BranchesAndTypes:  branchesAndTypes,
					Lineage:           lineage,
					Regexes:           regexes,
					ShowAllBranches:   false,
					UnknownBranchType: unknownBranchType,
				})
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

func TestSwitchBranch(t *testing.T) {
	t.Parallel()

	t.Run("SwitchBranchEntries", func(t *testing.T) {
		t.Parallel()
		t.Run("ContainsBranch", func(t *testing.T) {
			t.Parallel()
			t.Run("contains the branch", func(t *testing.T) {
				t.Parallel()
				entries := dialog.SwitchBranchEntries{
					{Branch: "branch-1"},
					{Branch: "branch-2"},
				}
				must.True(t, entries.ContainsBranch("branch-1"))
			})
			t.Run("does not contain the branch", func(t *testing.T) {
				t.Parallel()
				entries := dialog.SwitchBranchEntries{
					{Branch: "branch-1"},
				}
				must.False(t, entries.ContainsBranch("branch-2"))
			})
			t.Run("empty", func(t *testing.T) {
				t.Parallel()
				entries := dialog.SwitchBranchEntries{}
				must.False(t, entries.ContainsBranch("branch-2"))
			})
		})

		t.Run("IndexOf", func(t *testing.T) {
			t.Parallel()
			entries := dialog.SwitchBranchEntries{
				{Branch: "main", Indentation: "", OtherWorktree: false},
				{Branch: "alpha", Indentation: "", OtherWorktree: false},
				{Branch: "alpha1", Indentation: "", OtherWorktree: false},
				{Branch: "beta", Indentation: "", OtherWorktree: false},
			}
			tests := map[gitdomain.LocalBranchName]int{
				"alpha1": 2,
				"other":  0,
			}
			for give, want := range tests {
				must.EqOp(t, want, entries.IndexOf(give))
			}
		})
	})

	t.Run("View", func(t *testing.T) {
		t.Parallel()
		t.Run("only the main branch exists", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries(dialog.SwitchBranchEntries{
						{Branch: "main", Indentation: "", OtherWorktree: false},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   Some(0),
				UncommittedChanges: false,
				DisplayBranchTypes: configdomain.DisplayTypes{
					Quantifier:  configdomain.QuantifierNo,
					BranchTypes: []configdomain.BranchType{},
				},
			}
			have := model.View()
			want := `
> main


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   a all   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("multiple top-level branches", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries(dialog.SwitchBranchEntries{
						{Branch: "main", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeMainBranch},
						{Branch: "one", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeFeatureBranch},
						{Branch: "two", Indentation: "", OtherWorktree: true, Type: configdomain.BranchTypeFeatureBranch},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   Some(0),
				UncommittedChanges: false,
				DisplayBranchTypes: configdomain.DisplayTypes{
					Quantifier:  configdomain.QuantifierNo,
					BranchTypes: []configdomain.BranchType{},
				},
			}
			have := model.View()
			dim := "\x1b[2m"
			reset := "\x1b[0m"
			want := `
> main
  one
` + dim + `+ two` + reset + `


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   a all   enter/o accept   q/esc/ctrl-c abort`
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("stacked changes", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries(dialog.SwitchBranchEntries{
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
				InitialBranchPos:   Some(0),
				UncommittedChanges: false,
				DisplayBranchTypes: configdomain.DisplayTypes{
					Quantifier:  configdomain.QuantifierNo,
					BranchTypes: []configdomain.BranchType{},
				},
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


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   a all   enter/o accept   q/esc/ctrl-c abort`
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("stacked changes with types", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries(dialog.SwitchBranchEntries{
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
				InitialBranchPos:   Some(0),
				UncommittedChanges: false,
				DisplayBranchTypes: configdomain.DisplayTypes{
					Quantifier:  configdomain.QuantifierAll,
					BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeMainBranch},
				},
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


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   a all   enter/o accept   q/esc/ctrl-c abort`
			want = want[1:]
			must.EqOp(t, want, have)
		})

		t.Run("uncommitted changes", func(t *testing.T) {
			t.Parallel()
			model := dialog.SwitchModel{
				List: list.List[dialog.SwitchBranchEntry]{
					Cursor: 0,
					Entries: newSwitchBranchBubbleListEntries(dialog.SwitchBranchEntries{
						{Branch: "main", Indentation: "", OtherWorktree: false, Type: configdomain.BranchTypeMainBranch},
					}),
					MaxDigits:    1,
					NumberFormat: "%d",
				},
				InitialBranchPos:   Some(0),
				UncommittedChanges: true,
				DisplayBranchTypes: configdomain.DisplayTypes{
					Quantifier:  configdomain.QuantifierNo,
					BranchTypes: []configdomain.BranchType{},
				},
			}
			have := model.View()
			cyanBold := "\x1b[36;1m"
			reset := "\x1b[0m"
			want := `
` + cyanBold + `uncommitted changes` + reset + `


> main


  ↑/k up   ↓/j down   ←/u 10 up   →/d 10 down   a all   enter/o accept   q/esc/ctrl-c abort`[1:]
			must.EqOp(t, want, have)
		})
	})
}

func newSwitchBranchBubbleListEntries(entries dialog.SwitchBranchEntries) []list.Entry[dialog.SwitchBranchEntry] {
	result := make([]list.Entry[dialog.SwitchBranchEntry], len(entries))
	for e, entry := range entries {
		result[e] = list.Entry[dialog.SwitchBranchEntry]{
			Data:     entry,
			Disabled: entry.OtherWorktree,
			Text:     entry.String(),
		}
	}
	return result
}
