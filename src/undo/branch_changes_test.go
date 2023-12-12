package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
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
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames(),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("branch-1"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{},
			Active:   domain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
			Active: domain.NewLocalBranchName("branch-1"),
		}
		haveSpan := undo.NewBranchSpans(before, after)
		wantSpan := undo.BranchSpans{
			undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
		}
		must.Eq(t, wantSpan, haveSpan)
		haveChanges := haveSpan.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:            domain.NewLocalBranchNames("branch-1"),
			LocalRemoved:          domain.LocalBranchesSHAs{},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: domain.NewLocalBranchName("main")},
			&opcode.DeleteLocalBranch{
				Branch: domain.NewLocalBranchName("branch-1"),
				Force:  true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch removed", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames(),
		}
		lineage := configdomain.Lineage{}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
			Active: domain.NewLocalBranchName("branch-1"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{},
			Active:   domain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{
				domain.NewLocalBranchName("branch-1"): domain.NewSHA("111111"),
			},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.CreateBranch{
				Branch:        domain.NewLocalBranchName("branch-1"),
				StartingPoint: domain.NewSHA("111111").Location(),
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-1")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch changed", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				// a feature branch
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{
				domain.NewLocalBranchName("perennial-branch"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("333333"),
				},
				domain.NewLocalBranchName("feature-branch"): {
					Before: domain.NewSHA("222222"),
					After:  domain.NewSHA("444444"),
				},
			},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: domain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("444444"),
				SetToSHA:    domain.NewSHA("222222"),
				Hard:        true,
			},
			&opcode.Checkout{Branch: domain.NewLocalBranchName("perennial-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("333333"),
				SetToSHA:    domain.NewSHA("111111"),
				Hard:        true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("local-only branch pushed to origin", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded: domain.RemoteBranchNames{
				domain.NewRemoteBranchName("origin/perennial-branch"),
				domain.NewRemoteBranchName("origin/feature-branch"),
			},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.DeleteTrackingBranch{
				Branch: domain.NewRemoteBranchName("origin/perennial-branch"),
			},
			&opcode.DeleteTrackingBranch{
				Branch: domain.NewRemoteBranchName("origin/feature-branch"),
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("remote-only branch downloaded", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: domain.LocalBranchNames{
				domain.NewLocalBranchName("perennial-branch"),
				domain.NewLocalBranchName("feature-branch"),
			},
			LocalRemoved:          domain.LocalBranchesSHAs{},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.DeleteLocalBranch{
				Branch: domain.NewLocalBranchName("perennial-branch"),
				Force:  true,
			},
			&opcode.DeleteLocalBranch{
				Branch: domain.NewLocalBranchName("feature-branch"),
				Force:  true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch added", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{},
			Active:   domain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: domain.LocalBranchNames{
				domain.NewLocalBranchName("perennial-branch"),
				domain.NewLocalBranchName("feature-branch"),
			},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded: domain.RemoteBranchNames{
				domain.NewRemoteBranchName("origin/perennial-branch"),
				domain.NewRemoteBranchName("origin/feature-branch"),
			},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.DeleteTrackingBranch{
				Branch: domain.NewRemoteBranchName("origin/perennial-branch"),
			},
			&opcode.DeleteTrackingBranch{
				Branch: domain.NewRemoteBranchName("origin/feature-branch"),
			},
			&opcode.DeleteLocalBranch{
				Branch: domain.NewLocalBranchName("perennial-branch"),
				Force:  true,
			},
			&opcode.Checkout{Branch: domain.NewLocalBranchName("main")},
			&opcode.DeleteLocalBranch{
				Branch: domain.NewLocalBranchName("feature-branch"),
				Force:  true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{
				domain.NewLocalBranchName("perennial-branch"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("333333"),
				},
				domain.NewLocalBranchName("feature-branch"): {
					Before: domain.NewSHA("222222"),
					After:  domain.NewSHA("444444"),
				},
			},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: domain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("444444"),
				SetToSHA:    domain.NewSHA("222222"),
				Hard:        true,
			},
			&opcode.Checkout{Branch: domain.NewLocalBranchName("perennial-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("333333"),
				SetToSHA:    domain.NewSHA("111111"),
				Hard:        true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch remote updated", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("333333"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    domain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   domain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{
				domain.NewRemoteBranchName("origin/perennial-branch"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("222222"),
				},
				domain.NewRemoteBranchName("origin/feature-branch"): {
					Before: domain.NewSHA("333333"),
					After:  domain.NewSHA("444444"),
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't reset the remote perennial branch since those are assumed to be protected against force-pushes
			// and we can't revert the commit on it since we cannot change the local perennial branch here.
			&opcode.ResetRemoteBranchToSHA{
				Branch:      domain.NewRemoteBranchName("origin/feature-branch"),
				SetToSHA:    domain.NewSHA("333333"),
				MustHaveSHA: domain.NewSHA("444444"),
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("main"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("333333"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("main"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("555555"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("555555"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("666666"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("666666"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    domain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   domain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
			OmniRemoved:   domain.LocalBranchesSHAs{},
			OmniChanged: domain.LocalBranchChange{
				domain.NewLocalBranchName("main"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("444444"),
				},
				domain.NewLocalBranchName("perennial-branch"): {
					Before: domain.NewSHA("222222"),
					After:  domain.NewSHA("555555"),
				},
				domain.NewLocalBranchName("feature-branch"): {
					Before: domain.NewSHA("333333"),
					After:  domain.NewSHA("666666"),
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
			UndoablePerennialCommits: []domain.SHA{
				domain.NewSHA("444444"),
			},
		})
		wantProgram := program.Program{
			// revert the commit on the perennial branch
			&opcode.Checkout{Branch: domain.NewLocalBranchName("main")},
			&opcode.RevertCommit{SHA: domain.NewSHA("444444")},
			&opcode.PushCurrentBranch{CurrentBranch: domain.NewLocalBranchName("main"), NoPushHook: true},
			// reset the feature branch to the previous SHA
			&opcode.Checkout{Branch: domain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{MustHaveSHA: domain.NewSHA("666666"), SetToSHA: domain.NewSHA("333333"), Hard: true},
			&opcode.ForcePushCurrentBranch{NoPushHook: true},
			// check out the initial branch
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("upstream commit downloaded and branch shipped at the same time", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("main"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  domain.NewSHA("333333"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("main"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  domain.NewSHA("333333"),
				},
			},
			Active: domain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    domain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   domain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
			OmniRemoved: domain.LocalBranchesSHAs{
				domain.NewLocalBranchName("feature-branch"): domain.NewSHA("222222"),
			},
			OmniChanged: domain.LocalBranchChange{
				domain.NewLocalBranchName("main"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("444444"),
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
			UndoablePerennialCommits: []domain.SHA{
				domain.NewSHA("444444"),
			},
		})
		wantProgram := program.Program{
			// revert the undoable commit on the main branch
			&opcode.Checkout{Branch: domain.NewLocalBranchName("main")},
			&opcode.RevertCommit{SHA: domain.NewSHA("444444")},
			&opcode.PushCurrentBranch{CurrentBranch: domain.NewLocalBranchName("main"), NoPushHook: true},
			// re-create the feature branch
			&opcode.CreateBranch{Branch: domain.NewLocalBranchName("feature-branch"), StartingPoint: domain.NewSHA("222222").Location()},
			&opcode.CreateTrackingBranch{Branch: domain.NewLocalBranchName("feature-branch"), NoPushHook: true},
			// check out the initial branch
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("555555"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("666666"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    domain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   domain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
			OmniRemoved:   domain.LocalBranchesSHAs{},
			OmniChanged:   domain.LocalBranchChange{},
			InconsistentlyChanged: domain.InconsistentChanges{
				domain.InconsistentChange{
					Before: domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
					After: domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   domain.NewSHA("333333"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  domain.NewSHA("444444"),
					},
				},
				domain.InconsistentChange{
					Before: domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
						RemoteSHA:  domain.NewSHA("222222"),
					},
					After: domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("555555"),
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
						RemoteSHA:  domain.NewSHA("666666"),
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the perennial branch because it cannot force-push the changes to the remote branch.
			&opcode.Checkout{Branch: domain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("555555"),
				SetToSHA:    domain.NewSHA("222222"),
				Hard:        true,
			},
			&opcode.ResetRemoteBranchToSHA{
				Branch:      domain.NewRemoteBranchName("origin/feature-branch"),
				MustHaveSHA: domain.NewSHA("666666"),
				SetToSHA:    domain.NewSHA("222222"),
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pulled down", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{
				domain.NewLocalBranchName("perennial-branch"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("222222"),
				},
				domain.NewLocalBranchName("feature-branch"): {
					Before: domain.NewSHA("333333"),
					After:  domain.NewSHA("444444"),
				},
			},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.Checkout{Branch: domain.NewLocalBranchName("feature-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("444444"),
				SetToSHA:    domain.NewSHA("333333"),
				Hard:        true,
			},
			&opcode.Checkout{Branch: domain.NewLocalBranchName("perennial-branch")},
			&opcode.ResetCurrentBranchToSHA{
				MustHaveSHA: domain.NewSHA("222222"),
				SetToSHA:    domain.NewSHA("111111"),
				Hard:        true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch updates pushed up", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("333333"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:    domain.LocalBranchNames{},
			LocalRemoved:  domain.LocalBranchesSHAs{},
			LocalChanged:  domain.LocalBranchChange{},
			RemoteAdded:   domain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{
				domain.NewRemoteBranchName("origin/perennial-branch"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("222222"),
				},
				domain.NewRemoteBranchName("origin/feature-branch"): {
					Before: domain.NewSHA("333333"),
					After:  domain.NewSHA("444444"),
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			// It doesn't revert the remote perennial branch because it cannot force-push the changes to it.
			&opcode.ResetRemoteBranchToSHA{
				Branch:      domain.NewRemoteBranchName("origin/feature-branch"),
				MustHaveSHA: domain.NewSHA("444444"),
				SetToSHA:    domain.NewSHA("333333"),
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch deleted locally", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("main"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded: domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{
				domain.NewLocalBranchName("perennial-branch"): domain.NewSHA("111111"),
				domain.NewLocalBranchName("feature-branch"):   domain.NewSHA("222222"),
			},
			LocalChanged:          domain.LocalBranchChange{},
			RemoteAdded:           domain.RemoteBranchNames{},
			RemoteRemoved:         domain.RemoteBranchesSHAs{},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			&opcode.CreateBranch{
				Branch:        domain.NewLocalBranchName("feature-branch"),
				StartingPoint: domain.NewSHA("222222").Location(),
			},
			&opcode.CreateBranch{
				Branch:        domain.NewLocalBranchName("perennial-branch"),
				StartingPoint: domain.NewSHA("111111").Location(),
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("omnibranch tracking branch deleted", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/feature-branch"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-branch"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-branch"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded:  domain.RemoteBranchNames{},
			RemoteRemoved: domain.RemoteBranchesSHAs{
				domain.NewRemoteBranchName("origin/perennial-branch"): domain.NewSHA("111111"),
				domain.NewRemoteBranchName("origin/feature-branch"):   domain.NewSHA("222222"),
			},
			RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			// don't re-create the tracking branch for the perennial branch
			// because those are protected
			&opcode.CreateRemoteBranch{
				Branch:     domain.NewLocalBranchName("feature-branch"),
				SHA:        domain.NewSHA("222222"),
				NoPushHook: true,
			},
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("feature-branch")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})

	t.Run("sync with a new upstream remote", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames(),
		}
		lineage := configdomain.Lineage{
			domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
		}
		before := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("main"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			},
			Active: domain.NewLocalBranchName("main"),
		}
		after := domain.BranchesSnapshot{
			Branches: domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("main"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			},
			Active: domain.NewLocalBranchName("feature-branch"),
		}
		span := undo.NewBranchSpans(before, after)
		haveChanges := span.Changes()
		wantChanges := undo.BranchChanges{
			LocalAdded:   domain.LocalBranchNames{},
			LocalRemoved: domain.LocalBranchesSHAs{},
			LocalChanged: domain.LocalBranchChange{},
			RemoteAdded: domain.RemoteBranchNames{
				domain.NewRemoteBranchName("upstream/main"),
			},
			RemoteRemoved: domain.RemoteBranchesSHAs{},
			RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
			OmniRemoved:   domain.LocalBranchesSHAs{},
			OmniChanged: domain.LocalBranchChange{
				domain.NewLocalBranchName("main"): {
					Before: domain.NewSHA("111111"),
					After:  domain.NewSHA("222222"),
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
			UndoablePerennialCommits: []domain.SHA{},
		})
		wantProgram := program.Program{
			// No changes should happen here since all changes were syncs on perennial branches.
			// We don't want to undo these commits because that would undo commits
			// already committed to perennial branches by others for everybody on the team.
			&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("main")},
		}
		must.Eq(t, wantProgram, haveProgram)
	})
}
