package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/shoenig/test/must"
)

func TestChanges(t *testing.T) {
	t.Parallel()

	t.Run("local-only branch added", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames(),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("branch-1"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{},
			Active:   gitdomain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("branch-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
			Active: gitdomain.NewLocalBranchName("branch-1"),
		}
		haveSpan := undo.NewBranchSpans(before, after)
		wantSpan := undo.BranchSpans{
			undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("branch-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:            gitdomain.NewLocalBranchNames("branch-1"),
			LocalRemoved:          domain.LocalBranchesSHAs{},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcode.DeleteLocalBranch{
				Branch: gitdomain.NewLocalBranchName("branch-1"),
				Force:  true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch removed", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames(),
		}
		lineage := configdomain.Lineage{}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("branch-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
			Active: gitdomain.NewLocalBranchName("branch-1"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{},
			Active:   gitdomain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{
				gitdomain.NewLocalBranchName("branch-1"): gitdomain.NewSHA("111111"),
			},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.CreateBranch{
				Branch:        gitdomain.NewLocalBranchName("branch-1"),
				StartingPoint: gitdomain.NewSHA("111111").Location(),
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-1")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch changed", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				// a feature branch
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{
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
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("222222"),
				Hard:        true,
			},
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("333333"),
				SetToSHA:    gitdomain.NewSHA("111111"),
				Hard:        true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch pushed to origin", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"),
				gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
			},
			&opcode.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("remote-only branch downloaded", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("perennial-branch"),
				gitdomain.NewLocalBranchName("feature-branch"),
			},
			LocalRemoved:          domain.LocalBranchesSHAs{},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.DeleteLocalBranch{
				Branch: gitdomain.NewLocalBranchName("perennial-branch"),
				Force:  true,
			},
			&opcode.DeleteLocalBranch{
				Branch: gitdomain.NewLocalBranchName("feature-branch"),
				Force:  true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch added", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{},
			Active:   gitdomain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("perennial-branch"),
				gitdomain.NewLocalBranchName("feature-branch"),
			},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"),
				gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
			},
			&opcode.DeleteTrackingBranch{
				Branch: gitdomain.NewRemoteBranchName("origin/feature-branch"),
			},
			&opcode.DeleteLocalBranch{
				Branch: gitdomain.NewLocalBranchName("perennial-branch"),
				Force:  true,
			},
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcode.DeleteLocalBranch{
				Branch: gitdomain.NewLocalBranchName("feature-branch"),
				Force:  true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{
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
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               false,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("222222"),
				Hard:        true,
			},
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("333333"),
				SetToSHA:    gitdomain.NewSHA("111111"),
				Hard:        true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch remote updated", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("333333"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
				gitdomain.NewRemoteBranchName("origin/feature-branch"): {
					Before: gitdomain.NewSHA("333333"),
					After:  gitdomain.NewSHA("444444"),
				},
			},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               false,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't reset the remote perennial branch since those are assumed to be protected against force-pushes
			// and we can't revert the commit on it since we cannot change the local perennial branch here.
			&opcode.ResetRemoteBranchToSHA{
				Branch:      gitdomain.NewRemoteBranchName("origin/feature-branch"),
				SetToSHA:    gitdomain.NewSHA("333333"),
				MustHaveSHA: gitdomain.NewSHA("444444"),
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("main"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("333333"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("main"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("555555"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("666666"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:   domain.LocalBranchesSHAs{},
			OmniChanged: domain.LocalBranchChange{
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
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:       lineage,
			BranchTypes:   branchTypes,
			InitialBranch: before.Active,
			FinalBranch:   after.Active,
			NoPushHook:    true,
			UndoablePerennialCommits: []gitdomain.SHA{
				gitdomain.NewSHA("444444"),
			},
		})
		wantProgram := program.Program{
			// revert the commit on the perennial branch
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcode.RevertCommit{SHA: gitdomain.NewSHA("444444")},
			&opcode.PushCurrentBranch{CurrentBranch: gitdomain.NewLocalBranchName("main"), NoPushHook: true},
			// reset the feature branch to the previous SHA
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{MustHaveSHA: gitdomain.NewSHA("666666"), SetToSHA: gitdomain.NewSHA("333333"), Hard: true},
			&opcode.ForcePushCurrentBranch{NoPushHook: true},
			// check out the initial branch
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("upstream commit downloaded and branch shipped at the same time", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("main"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  gitdomain.NewSHA("333333"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("main"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  gitdomain.NewSHA("333333"),
				},
			},
			Active: gitdomain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved: domain.LocalBranchesSHAs{
				gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewSHA("222222"),
			},
			OmniChanged: domain.LocalBranchChange{
				gitdomain.NewLocalBranchName("main"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("444444"),
				},
			},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:       lineage,
			BranchTypes:   branchTypes,
			InitialBranch: before.Active,
			FinalBranch:   after.Active,
			NoPushHook:    true,
			UndoablePerennialCommits: []gitdomain.SHA{
				gitdomain.NewSHA("444444"),
			},
		})
		wantProgram := program.Program{
			// revert the undoable commit on the main branch
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("main")},
			&opcode.RevertCommit{SHA: gitdomain.NewSHA("444444")},
			&opcode.PushCurrentBranch{CurrentBranch: gitdomain.NewLocalBranchName("main"), NoPushHook: true},
			// re-create the feature branch
			&opcode.CreateBranch{Branch: gitdomain.NewLocalBranchName("feature-branch"), StartingPoint: gitdomain.NewSHA("222222").Location()},
			&opcode.CreateTrackingBranch{Branch: gitdomain.NewLocalBranchName("feature-branch"), NoPushHook: true},
			// check out the initial branch
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:   domain.LocalBranchesSHAs{},
			OmniChanged:   domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{
				domain.InconsistentChange{
					Before: domain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   gitdomain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  gitdomain.NewSHA("111111"),
					},
					After: domain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   gitdomain.NewSHA("333333"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  gitdomain.NewSHA("444444"),
					},
				},
				domain.InconsistentChange{
					Before: domain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
						LocalSHA:   gitdomain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
						RemoteSHA:  gitdomain.NewSHA("222222"),
					},
					After: domain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
						LocalSHA:   gitdomain.NewSHA("555555"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
						RemoteSHA:  gitdomain.NewSHA("666666"),
					},
				},
			},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the perennial branch because it cannot force-push the changes to the remote branch.
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("555555"),
				SetToSHA:    gitdomain.NewSHA("222222"),
				Hard:        true,
			},
			&opcode.ResetRemoteBranchToSHA{
				Branch:      gitdomain.NewRemoteBranchName("origin/feature-branch"),
				MustHaveSHA: gitdomain.NewSHA("666666"),
				SetToSHA:    gitdomain.NewSHA("222222"),
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pulled down", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{
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
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("333333"),
				Hard:        true,
			},
			&opcode.Checkout{Branch: gitdomain.NewLocalBranchName("perennial-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: gitdomain.NewSHA("222222"),
				SetToSHA:    gitdomain.NewSHA("111111"),
				Hard:        true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pushed up", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("333333"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    gitdomain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   gitdomain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
				gitdomain.NewRemoteBranchName("origin/feature-branch"): {
					Before: gitdomain.NewSHA("333333"),
					After:  gitdomain.NewSHA("444444"),
				},
			},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the remote perennial branch because it cannot force-push the changes to it.
			&opcode.ResetRemoteBranchToSHA{
				Branch:      gitdomain.NewRemoteBranchName("origin/feature-branch"),
				MustHaveSHA: gitdomain.NewSHA("444444"),
				SetToSHA:    gitdomain.NewSHA("333333"),
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch deleted locally", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{
				gitdomain.NewLocalBranchName("perennial-branch"): gitdomain.NewSHA("111111"),
				gitdomain.NewLocalBranchName("feature-branch"):   gitdomain.NewSHA("222222"),
			},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           gitdomain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.CreateBranch{
				Branch:        gitdomain.NewLocalBranchName("feature-branch"),
				StartingPoint: gitdomain.NewSHA("222222").Location(),
			},
			&opcode.CreateBranch{
				Branch:        gitdomain.NewLocalBranchName("perennial-branch"),
				StartingPoint: gitdomain.NewSHA("111111").Location(),
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch tracking branch deleted", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-branch"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded:  gitdomain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{
				gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("111111"),
				gitdomain.NewRemoteBranchName("origin/feature-branch"):   gitdomain.NewSHA("222222"),
			},
			RemoteChanged:         map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:           domain.LocalBranchesSHAs{},
			OmniChanged:           domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// don't re-create the tracking branch for the perennial branch
			// because those are protected
			&opcode.CreateRemoteBranch{
				Branch:     gitdomain.NewLocalBranchName("feature-branch"),
				SHA:        gitdomain.NewSHA("222222"),
				NoPushHook: true,
			},
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("sync with a new upstream remote", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames(),
		}
		lineage := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature-branch"): gitdomain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("main"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
			},
			Active: gitdomain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("main"),
					LocalSHA:   gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
				},
			},
			Active: gitdomain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   gitdomain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded: gitdomain.RemoteBranchNames{
				gitdomain.NewRemoteBranchName("upstream/main"),
			},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[gitdomain.RemoteBranchName]domain.Change[gitdomain.SHA]{},
			OmniRemoved:   domain.LocalBranchesSHAs{},
			OmniChanged: domain.LocalBranchChange{
				gitdomain.NewLocalBranchName("main"): {
					Before: gitdomain.NewSHA("111111"),
					After:  gitdomain.NewSHA("222222"),
				},
			},
			InconsistentlyChanged: domain.InconsistentChanges{},
		}
		must.Eq(t, wantChanges, haveChanges)
		haveProgram := haveChanges.UndoProgram(undo.BranchChangesUndoProgramArgs{
			Lineage:                  lineage,
			BranchTypes:              branchTypes,
			InitialBranch:            before.Active,
			FinalBranch:              after.Active,
			NoPushHook:               true,
			UndoablePerennialCommits: []gitdomain.SHA{},
		})
		wantProgram := program.Program{
			// No changes should happen here since all changes were syncs on perennial branches.
			// We don't want to undo these commits because that would undo commits
			// already committed to perennial branches by others for everybody on the team.
			&opcode.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}
