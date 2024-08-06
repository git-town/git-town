package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/undo/undobranches"
	"github.com/git-town/git-town/v15/internal/undo/undodomain"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/shoenig/test/must"
)

func TestChanges(t *testing.T) {
	t.Parallel()

	t.Run("local-only branch added", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{},
			Active:   Some(gitdomain.NewLocalBranchName("main")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("branch-1")),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
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
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("branch-1"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PushHook:          false,
				PerennialBranches: gitdomain.NewLocalBranchNames(),
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcodes.DeleteLocalBranch{Branch: gitdomain.NewLocalBranchName("branch-1")},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch removed", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("branch-1")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{},
			Active:   Some(gitdomain.NewLocalBranchName("main")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{
				gitdomain.NewLocalBranchName("branch-1"): gitdomain.NewSHA("111111"),
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
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           configdomain.Lineage{},
				PerennialBranches: gitdomain.NewLocalBranchNames(),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.CreateBranch{
				Branch:        gitdomain.NewLocalBranchName("branch-1"),
				StartingPoint: gitdomain.NewSHA("111111").Location(),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-1")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch changed", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				// a feature branch
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{
				gitdomain.NewLocalBranchName("perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("333333"),
				},
				gitdomain.NewLocalBranchName("feature-branch"): {
					Before: gitdomain.NewSHA("222222"),
					After:  gitdomain.NewSHA("444444"),
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
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				Lineage:           lineage,
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("222222"),
				Hard:        true,
			},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("333333"),
				SetToSHA:    gitdomain.NewSHA("111111"),
				Hard:        true,
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch pushed to origin", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"),
				gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
			},
			&opcodes.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("remote-only branch downloaded", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("main")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("main")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("perennial-branch"),
				gitdomain.NewLocalBranchName("feature-branch"),
			},
			LocalRemoved:          undobranches.LocalBranchesSHAs{},
			LocalChanged:          undobranches.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.DeleteLocalBranch{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcodes.DeleteLocalBranch{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch added", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{},
			Active:   Some(gitdomain.NewLocalBranchName("main")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("perennial-branch"),
				gitdomain.NewLocalBranchName("feature-branch"),
			},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"),
				gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			RemoteRemoved:         undobranches.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
			},
			&opcodes.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			&opcodes.DeleteLocalBranch{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcodes.DeleteLocalBranch{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{
				gitdomain.NewLocalBranchName("perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("333333"),
				},
				gitdomain.NewLocalBranchName("feature-branch"): {
					Before: gitdomain.NewSHA("222222"),
					After:  gitdomain.NewSHA("444444"),
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
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          true,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("222222"),
				Hard:        true,
			},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("333333"),
				SetToSHA:    gitdomain.NewSHA("111111"),
				Hard:        true,
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch remote updated", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		haveSpan := undobranches.NewBranchSpans(before, after)
		wantSpan := undobranches.BranchSpans{
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			},
			undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					RemoteName: Some(gitdomain.RemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
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
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
				gitdomain.NewRemoteBranchName("origin/feature-branch"): {
					Before: gitdomain.NewSHA("333333"),
					After:  gitdomain.NewSHA("444444"),
				},
			},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          true,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't reset the remote perennial branch since those are assumed to be protected against force-pushes
			// and we can't revert the commit on it since we cannot change the local perennial branch here.
			&opcodes.ResetRemoteBranchToSHA{
				Branch:      gitdomain.NewRemoteBranchName("origin/feature-branch"),
				SetToSHA:    gitdomain.NewSHA("333333"),
				MustHaveSHA: gitdomain.NewSHA("444444"),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("555555")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("555555")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("666666")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			OmniChanged: undobranches.LocalBranchChange{
				gitdomain.NewLocalBranchName("main"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("444444"),
				},
				gitdomain.NewLocalBranchName("perennial-branch"): {
					Before: gitdomain.NewSHA("222222"),
					After:  gitdomain.NewSHA("555555"),
				},
				gitdomain.NewLocalBranchName("feature-branch"): {
					Before: gitdomain.NewSHA("333333"),
					After:  gitdomain.NewSHA("666666"),
				},
			},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch: before.Active.GetOrPanic(),
			Config:      config,
			EndBranch:   after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{
				gitdomain.NewSHA("444444"),
			},
		})
		wantProgram := program.Program{
			// revert the commit on the perennial branch
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcodes.RevertCommit{SHA: gitdomain.NewSHA("444444")},
			&opcodes.PushCurrentBranch{CurrentBranch: gitdomain.NewLocalBranchName("main")},
			// reset the feature branch to the previous SHA
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.ResetCurrentBranchToSHA{MustHaveSHA: gitdomain.NewSHA("666666"), SetToSHA: gitdomain.NewSHA("333333"), Hard: true},
			&opcodes.ForcePushCurrentBranch{},
			// check out the initial branch
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("upstream commit downloaded and branch shipped at the same time", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("main")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved: undobranches.LocalBranchesSHAs{
				gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewSHA("222222"),
			},
			OmniChanged: undobranches.LocalBranchChange{
				gitdomain.NewLocalBranchName("main"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("444444"),
				},
			},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch: before.Active.GetOrPanic(),
			Config:      config,
			EndBranch:   after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{
				gitdomain.NewSHA("444444"),
			},
		})
		wantProgram := program.Program{
			// revert the undoable commit on the main branch
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcodes.RevertCommit{SHA: gitdomain.NewSHA("444444")},
			&opcodes.PushCurrentBranch{CurrentBranch: gitdomain.NewLocalBranchName("main")},
			// re-create the feature branch
			&opcodes.CreateBranch{Branch: gitdomain.NewLocalBranchName("feature-branch"), StartingPoint: gitdomain.NewSHA("222222").Location()},
			&opcodes.CreateTrackingBranch{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			// check out the initial branch
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("555555")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			OmniChanged:   undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{
				undodomain.InconsistentChange{
					Before: gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					},
					After: gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
						LocalSHA:   Some(gitdomain.NewSHA("333333")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					},
				},
				undodomain.InconsistentChange{
					Before: gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
						LocalSHA:   Some(gitdomain.NewSHA("222222")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					},
					After: gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
						LocalSHA:   Some(gitdomain.NewSHA("555555")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("666666")),
					},
				},
			},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the perennial branch because it cannot force-push the changes to the remote branch.
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("555555"),
				SetToSHA:    gitdomain.NewSHA("222222"),
				Hard:        true,
			},
			&opcodes.ResetRemoteBranchToSHA{
				Branch:      gitdomain.NewRemoteBranchName("origin/feature-branch"),
				MustHaveSHA: gitdomain.NewSHA("666666"),
				SetToSHA:    gitdomain.NewSHA("222222"),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pulled down", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{
				gitdomain.NewLocalBranchName("perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
				gitdomain.NewLocalBranchName("feature-branch"): {
					Before: gitdomain.NewSHA("333333"),
					After:  gitdomain.NewSHA("444444"),
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
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("333333"),
				Hard:        true,
			},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcodes.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("222222"),
				SetToSHA:    gitdomain.NewSHA("111111"),
				Hard:        true,
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pushed up", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  undobranches.LocalBranchesSHAs{},
			LocalChanged:  undobranches.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
				gitdomain.NewRemoteBranchName("origin/feature-branch"): {
					Before: gitdomain.NewSHA("333333"),
					After:  gitdomain.NewSHA("444444"),
				},
			},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the remote perennial branch because it cannot force-push the changes to it.
			&opcodes.ResetRemoteBranchToSHA{
				Branch:      gitdomain.NewRemoteBranchName("origin/feature-branch"),
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("333333"),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch deleted locally", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("main")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{
				gitdomain.NewLocalBranchName("perennial-branch"): gitdomain.NewSHA("111111"),
				gitdomain.NewLocalBranchName("feature-branch"):   gitdomain.NewSHA("222222"),
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
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcodes.CreateBranch{
				Branch:        gitdomain.NewLocalBranchName("feature-branch"),
				StartingPoint: gitdomain.NewSHA("222222").Location(),
			},
			&opcodes.CreateBranch{
				Branch:        gitdomain.NewLocalBranchName("perennial-branch"),
				StartingPoint: gitdomain.NewSHA("111111").Location(),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch tracking branch deleted", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-branch")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-branch")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded:  gitdomain.RemoteBranchNames{},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("111111"),
				gitdomain.NewRemoteBranchName("origin/feature-branch"):   gitdomain.NewSHA("222222"),
			},
			RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:           undobranches.LocalBranchesSHAs{},
			OmniChanged:           undobranches.LocalBranchChange{},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrDefault(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// don't re-create the tracking branch for the perennial branch
			// because those are protected
			&opcodes.CreateRemoteBranch{
				Branch: gitdomain.NewLocalBranchName("feature-branch"),
				SHA:    gitdomain.NewSHA("222222"),
			},
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("sync with a new upstream remote", func(t *testing.T) {
		t.Parallel()
		before := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("main")),
		}
		after := gitdomain.BranchesSnapshot{
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			},
			Active: Some(gitdomain.NewLocalBranchName("feature-branch")),
		}
		span := undobranches.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undobranches.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: undobranches.LocalBranchesSHAs{},
			LocalChanged: undobranches.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				gitdomain.NewRemoteBranchName("upstream/main"),
			},
			RemoteRemoved: undobranches.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
			OmniRemoved:   undobranches.LocalBranchesSHAs{},
			OmniChanged: undobranches.LocalBranchChange{
				gitdomain.NewLocalBranchName("main"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
			},
			InconsistentlyChanged: undodomain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		lineage := configdomain.NewLineage()
		lineage.Add(gitdomain.NewLocalBranchName("feature-branch"), gitdomain.NewLocalBranchName("main"))
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: &configdomain.UnvalidatedConfig{
				Lineage:           lineage,
				PerennialBranches: gitdomain.NewLocalBranchNames(),
				PushHook:          false,
			},
		}
		haveProgram := haveChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
			BeginBranch:              before.Active.GetOrPanic(),
			Config:                   config,
			EndBranch:                after.Active.GetOrPanic(),
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// No changes should happen here since all changes were syncs on perennial branches.
			// We don't want to undo these commits because that would undo commits
			// already committed to perennial branches by others for everybody on the team.
			&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}
