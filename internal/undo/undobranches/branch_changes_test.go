package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/undo/undobranches"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestChanges(t *testing.T) {
	t.Parallel()

	t.Run("local-only branch added", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{},
			Active:   gitdomain.NewLocalBranchNameOption("main"),
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
			Active: gitdomain.NewLocalBranchNameOption("branch-1"),
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
			LocalAdded:            gitdomain.NewLocalBranchNames("branch-1"),
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"branch-1": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PushHook:          false,
				PerennialBranches: gitdomain.LocalBranchNames{},
			},
		}
		haveProgram, _ := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
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
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{
				"perennial-branch": {
					Before: "111111",
					After:  "333333",
				},
				"feature-branch": {
					Before: "222222",
					After:  "444444",
				},
			},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
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
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				"origin/perennial-branch",
				"origin/feature-branch",
			},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
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
			Active: gitdomain.NewLocalBranchNameOption("branch-1"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{},
			Active:   gitdomain.NewLocalBranchNameOption("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{
				"branch-1": "111111",
			},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				PerennialBranches: gitdomain.LocalBranchNames{},
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"branch-1"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch added", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{},
			Active:   gitdomain.NewLocalBranchNameOption("main"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{
				"perennial-branch",
				"feature-branch",
			},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				"origin/perennial-branch",
				"origin/feature-branch",
			},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch changed locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			LocalChanged: undobranches.LocalBranchChange{
				"perennial-branch": {
					Before: "111111",
					After:  "333333",
				},
				"feature-branch": {
					Before: "222222",
					After:  "444444",
				},
			},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          true,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "555555"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			OmniChanged:   undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{
				undodomain.InconsistentChange{
					Before: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					},
					After: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "333333"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					},
				},
				undodomain.InconsistentChange{
					Before: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					},
					After: gitdomain.BranchInfo{
						Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "555555"}),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("666666")),
					},
				},
			},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "555555"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("555555")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "666666"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			OmniChanged: undobranches.LocalBranchChange{
				"main": {
					Before: "111111",
					After:  "444444",
				},
				"perennial-branch": {
					Before: "222222",
					After:  "555555",
				},
				"feature-branch": {
					Before: "333333",
					After:  "666666",
				},
			},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "main", "perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch deleted locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{
				"perennial-branch": "111111",
				"feature-branch":   "222222",
			},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch remote updated", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			},
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{
				"origin/perennial-branch": {
					Before: "111111",
					After:  "222222",
				},
				"origin/feature-branch": {
					Before: "333333",
					After:  "444444",
				},
			},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          true,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
				SetToSHA:    "333333",
				MustHaveSHA: "444444",
			},
			&opcodes.CheckoutIfExists{Branch: "feature-branch"},
		}
		must.Eq(t, wantProgram, haveProgram)
		must.Eq(t, nil, haveChangedBranches)
	})

	t.Run("omnibranch renamed locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "old", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("old"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "new", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("new"),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "old", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "new", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/old")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{
				{
					Before: gitdomain.NewLocalBranchName("old"),
					After:  gitdomain.NewLocalBranchName("new"),
				},
			},
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"old": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"new", "old"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch tracking branch deleted", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
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
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded:  gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{
				"origin/perennial-branch": "111111",
				"origin/feature-branch":   "222222",
			},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		must.Eq(t, nil, haveChangedBranches)
	})

	t.Run("omnibranch updates pulled down", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "333333"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			LocalChanged: undobranches.LocalBranchChange{
				"perennial-branch": {
					Before: "111111",
					After:  "222222",
				},
				"feature-branch": {
					Before: "333333",
					After:  "444444",
				},
			},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("omnibranch updates pushed up", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{
				"origin/perennial-branch": {
					Before: "111111",
					After:  "222222",
				},
				"origin/feature-branch": {
					Before: "333333",
					After:  "444444",
				},
			},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, _ := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("main"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-branch", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:            gitdomain.LocalBranchNames{"perennial-branch", "feature-branch"},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalRenamed:          []undobranches.LocalBranchRename{},
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "perennial-branch"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("sync with a new upstream remote", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("main"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalRenamed: []undobranches.LocalBranchRename{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				"upstream/main",
			},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			OmniChanged: undobranches.LocalBranchChange{
				"main": {
					Before: "111111",
					After:  "222222",
				},
			},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.LocalBranchNames{},
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"main"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})

	t.Run("upstream commit downloaded and branch shipped at the same time", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-branch", SHA: "222222"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("feature-branch"),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "main", SHA: "444444"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: gitdomain.NewLocalBranchNameOption("main"),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalRenamed:  []undobranches.LocalBranchRename{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved: undobranches.LocalBranchesSHAs{
				"feature-branch": "222222",
			},
			OmniChanged: undobranches.LocalBranchChange{
				"main": {
					Before: "111111",
					After:  "444444",
				},
			},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineageWith(configdomain.LineageData{
			"feature-branch": "main",
		})
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram, haveChangedBranches := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
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
		wantChangedBranches := gitdomain.LocalBranchNames{"feature-branch", "main"}
		must.Eq(t, wantChangedBranches, haveChangedBranches)
	})
}
