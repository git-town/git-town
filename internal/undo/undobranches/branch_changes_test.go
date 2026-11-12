package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/config"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v24/internal/undo/undobranches"
	"github.com/git-town/git-town/v24/internal/undo/undodomain"
	"github.com/git-town/git-town/v24/internal/vm/opcodes"
	"github.com/git-town/git-town/v24/internal/vm/program"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestChanges(t *testing.T) {
	t.Parallel()

	t.Run("local-only branch added", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Active:   gitdomain.LocalBranchNameOpt("main"),
			Branches: gitdomain.BranchInfos{},
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
			},
			Active: gitdomain.LocalBranchNameOpt("branch-1"),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.NewLocalBranchNames("branch-1"),
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged:           undobranches.LocalBranchChange{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"branch-1": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.LocalBranchNames{},
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.CheckoutIfNeeded{Branch: "main"},
			&opcodes.BranchLocalDelete{Branch: "branch-1"},
			&opcodes.CheckoutIfExists{Branch: "main"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch changed", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				// a feature branch
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged: undobranches.LocalBranchChange{
				"feature-branch": {
					After:  "444444",
					Before: "222222",
				},
				"perennial-branch": {
					After:  "333333",
					Before: "111111",
				},
			},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.CheckoutIfNeeded{Branch: "feature-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "444444",
				SetToSHA:    "222222",
			},
			&opcodes.CheckoutIfNeeded{Branch: "perennial-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "333333",
				SetToSHA:    "111111",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch pushed to origin", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged:           undobranches.LocalBranchChange{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				"origin/perennial-branch",
				"origin/feature-branch",
			},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.BranchTrackingDelete{
				Branch: "origin/perennial-branch",
			},
			&opcodes.BranchTrackingDelete{
				Branch: "origin/feature-branch",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch removed", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("branch-1"),
		}
		after := gitdomain.BranchesSnapshot{
			Active:   gitdomain.LocalBranchNameOpt("main"),
			Branches: gitdomain.BranchInfos{},
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved: undobranches.LocalBranchesSHAs{
				"branch-1": "111111",
			},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				PerennialBranches: gitdomain.LocalBranchNames{},
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.BranchCreate{
				Branch:        "branch-1",
				StartingPoint: "111111",
			},
			&opcodes.CheckoutIfExists{Branch: "branch-1"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch added", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Active:   gitdomain.LocalBranchNameOpt("main"),
			Branches: gitdomain.BranchInfos{},
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded: gitdomain.LocalBranchNames{
				"perennial-branch",
				"feature-branch",
			},
			LocalChanged: undobranches.LocalBranchChange{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			OmniChanged:  undobranches.LocalBranchChange{},
			OmniRemoved:  undobranches.LocalBranchesSHAs{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				"origin/perennial-branch",
				"origin/feature-branch",
			},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.BranchLocalDelete{Branch: "perennial-branch"},
			&opcodes.CheckoutIfNeeded{Branch: "main"},
			&opcodes.BranchLocalDelete{Branch: "feature-branch"},
			&opcodes.BranchTrackingDelete{
				Branch: "origin/perennial-branch",
			},
			&opcodes.BranchTrackingDelete{
				Branch: "origin/feature-branch",
			},
			&opcodes.CheckoutIfExists{Branch: "main"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged: undobranches.LocalBranchChange{
				"feature-branch": {
					After:  "444444",
					Before: "222222",
				},
				"perennial-branch": {
					After:  "333333",
					Before: "111111",
				},
			},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          true,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.CheckoutIfNeeded{Branch: "feature-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "444444",
				SetToSHA:    "222222",
			},
			&opcodes.CheckoutIfNeeded{Branch: "perennial-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "333333",
				SetToSHA:    "111111",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "555555"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("666666")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{
				undodomain.InconsistentChange{
					Before: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
						RemoteSHA:  Some(gitdomain.SHA("111111")),
					},
					After: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
						RemoteSHA:  Some(gitdomain.SHA("444444")),
					},
				},
				undodomain.InconsistentChange{
					Before: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
						RemoteSHA:  Some(gitdomain.SHA("222222")),
					},
					After: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "555555"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
						RemoteSHA:  Some(gitdomain.SHA("666666")),
					},
				},
			},
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalChanged:  undobranches.LocalBranchChange{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the perennial branch because it cannot force-push the changes to the remote branch.
			&opcodes.CheckoutIfNeeded{Branch: "feature-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "555555",
				SetToSHA:    "222222",
			},
			&opcodes.BranchRemoteSetToSHAIfNeeded{
				Branch:      "origin/feature-branch",
				MustHaveSHA: "666666",
				SetToSHA:    "222222",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("333333")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "555555"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("555555")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "666666"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("666666")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged: undobranches.LocalBranchChange{
				"feature-branch": {
					After:  "666666",
					Before: "333333",
				},
				"main": {
					After:  "444444",
					Before: "111111",
				},
				"perennial-branch": {
					After:  "555555",
					Before: "222222",
				},
			},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{"444444"},
		})
		wantProgram := program.Program{
			// revert the commit on the perennial branch
			&opcodes.CheckoutIfNeeded{Branch: "main"},
			&opcodes.CommitRevertIfNeeded{SHA: "444444"},
			&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: "main", TrackingBranch: "origin/main"},
			// reset the feature branch to the previous SHA
			&opcodes.CheckoutIfNeeded{Branch: "feature-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{MustHaveSHA: "666666", SetToSHA: "333333"},
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: "feature-branch", ForceIfIncludes: true, TrackingBranch: "origin/feature-branch"},
			// check out the initial branch
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch deleted locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved: undobranches.LocalBranchesSHAs{
				"feature-branch":   "222222",
				"perennial-branch": "111111",
			},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.BranchCreate{
				Branch:        "feature-branch",
				StartingPoint: "222222",
			},
			&opcodes.BranchCreate{
				Branch:        "perennial-branch",
				StartingPoint: "111111",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch remote updated", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("333333")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			},
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged:           undobranches.LocalBranchChange{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{
				"origin/feature-branch": {
					After:  "444444",
					Before: "333333",
				},
				"origin/perennial-branch": {
					After:  "222222",
					Before: "111111",
				},
			},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          true,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't reset the remote perennial branch since those are assumed to be protected against force-pushes
			// and we can't revert the commit on it since we cannot change the local perennial branch here.
			&opcodes.BranchRemoteSetToSHAIfNeeded{
				Branch:      "origin/feature-branch",
				MustHaveSHA: "444444",
				SetToSHA:    "333333",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch renamed locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "old", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("old"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "new", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("new"),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "old", SHA: "111111"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "new", SHA: "111111"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{
				{
					Before: gitdomain.LocalBranchName("old"),
					After:  gitdomain.LocalBranchName("new"),
				},
			},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"old": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrZero(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.BranchLocalRename{
				NewName: "old",
				OldName: "new",
			},
			&opcodes.CheckoutIfExists{
				Branch: "old",
			},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch tracking branch deleted", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged:           undobranches.LocalBranchChange{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{
				"origin/feature-branch":   "222222",
				"origin/perennial-branch": "111111",
			},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrZero(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// don't re-create the tracking branch for the perennial branch
			// because those are protected
			&opcodes.BranchRemoteCreate{
				Branch: "feature-branch",
				SHA:    "222222",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pulled down", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged: undobranches.LocalBranchChange{
				"feature-branch": {
					After:  "444444",
					Before: "333333",
				},
				"perennial-branch": {
					After:  "222222",
					Before: "111111",
				},
			},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			OmniChanged:   undobranches.LocalBranchChange{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.CheckoutIfNeeded{Branch: "feature-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "444444",
				SetToSHA:    "333333",
			},
			&opcodes.CheckoutIfNeeded{Branch: "perennial-branch"},
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				MustHaveSHA: "222222",
				SetToSHA:    "111111",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pushed up", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("333333")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged:           undobranches.LocalBranchChange{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{
				"origin/feature-branch": {
					After:  "444444",
					Before: "333333",
				},
				"origin/perennial-branch": {
					After:  "222222",
					Before: "111111",
				},
			},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the remote perennial branch because it cannot force-push the changes to it.
			&opcodes.BranchRemoteSetToSHAIfNeeded{
				Branch:      "origin/feature-branch",
				MustHaveSHA: "444444",
				SetToSHA:    "333333",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("remote-only branch downloaded", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("main"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{"perennial-branch", "feature-branch"},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged:           undobranches.LocalBranchChange{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.BranchLocalDelete{Branch: "perennial-branch"},
			&opcodes.BranchLocalDelete{Branch: "feature-branch"},
			&opcodes.CheckoutIfExists{Branch: "main"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("sync with a new upstream remote", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("main"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged: undobranches.LocalBranchChange{
				"main": {
					After:  "222222",
					Before: "111111",
				},
			},
			OmniRemoved: undobranches.LocalBranchesSHAs{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				"upstream/main",
			},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.LocalBranchNames{},
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// No changes should happen here since all changes were syncs on perennial branches.
			// We don't want to undo these commits because that would undo commits
			// already committed to perennial branches by others for everybody on the team.
			&opcodes.CheckoutIfExists{Branch: "main"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("upstream commit downloaded and branch shipped at the same time", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.SHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.SHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.SHA("333333")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.RemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.SHA("444444")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.RemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.SHA("333333")),
				},
			},
			Active: gitdomain.LocalBranchNameOpt("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			InconsistentlyChanged: undodomain.InconsistentChanges{},
			LocalAdded:            gitdomain.LocalBranchNames{},
			LocalChanged:          undobranches.LocalBranchChange{},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			OmniChanged: undobranches.LocalBranchChange{
				"main": {
					After:  "444444",
					Before: "111111",
				},
			},
			OmniRemoved: undobranches.LocalBranchesSHAs{
				"feature-branch": "222222",
			},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			BranchInfos:              before.Branches,
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			FinalMessages:            stringslice.NewCollector(),
			UndoablePerennialCommits: []gitdomain.SHA{"444444"},
		})
		wantProgram := program.Program{
			// revert the undoable commit on the main branch
			&opcodes.CheckoutIfNeeded{Branch: "main"},
			&opcodes.CommitRevertIfNeeded{SHA: "444444"},
			&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: "main", TrackingBranch: "origin/main"},
			// re-create the feature branch
			&opcodes.BranchCreate{Branch: "feature-branch", StartingPoint: "222222"},
			&opcodes.BranchTrackingCreate{Branch: "feature-branch"},
			// check out the initial branch
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}
