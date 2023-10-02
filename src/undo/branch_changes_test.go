package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/stretchr/testify/assert"
)

func TestChanges(t *testing.T) {
	t.Parallel()

	t.Run("Steps", func(t *testing.T) {
		t.Parallel()
		t.Run("local-only branch added", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames(),
			}
			lineage := config.Lineage{
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
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
				Active: domain.NewLocalBranchName("branch-1"),
			}
			haveSpan := undo.NewBranchSpans(before, after)
			wantSpan := undo.BranchSpans{
				undo.BranchSpan{
					Before: domain.BranchInfo{
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
						SyncStatus: domain.SyncStatusUpToDate,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
					After: domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
			}
			assert.Equal(t, wantSpan, haveSpan)
			haveChanges := haveSpan.Changes()
			wantChanges := undo.BranchChanges{
				LocalAdded:            domain.NewLocalBranchNames("branch-1"),
				LocalRemoved:          domain.LocalBranchesSHAs{},
				LocalChanged:          domain.LocalBranchChange{},
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("main")},
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("branch-1"),
						Parent: domain.NewLocalBranchName("main").Location(),
						Force:  true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("main")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("local-only branch removed", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames(),
			}
			lineage := config.Lineage{}
			before := domain.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
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
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.CreateBranchStep{
						Branch:        domain.NewLocalBranchName("branch-1"),
						StartingPoint: domain.NewSHA("111111").Location(),
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-1")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("local-only branch changed", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
				domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
			}
			before := domain.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
					// a feature branch
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
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
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("444444"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
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
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("feature-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("444444"),
						SetToSHA:    domain.NewSHA("222222"),
						Hard:        true,
					},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("perennial-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("333333"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("local-only branch pushed to origin", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
				domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
			}
			before := domain.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
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
				RemoteAdded: []domain.RemoteBranchName{
					domain.NewRemoteBranchName("origin/perennial-branch"),
					domain.NewRemoteBranchName("origin/feature-branch"),
				},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.DeleteTrackingBranchStep{
						Branch: domain.NewRemoteBranchName("origin/perennial-branch"),
					},
					&steps.DeleteTrackingBranchStep{
						Branch: domain.NewRemoteBranchName("origin/feature-branch"),
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("remote-only branch downloaded", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
				domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
			}
			before := domain.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
						SyncStatus: domain.SyncStatusRemoteOnly,
						RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
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
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("perennial-branch"),
						Parent: domain.LocalBranchName{}.Location(),
						Force:  true,
					},
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("feature-branch"),
						Parent: domain.NewLocalBranchName("main").Location(),
						Force:  true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("main")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch added", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
				RemoteAdded: []domain.RemoteBranchName{
					domain.NewRemoteBranchName("origin/perennial-branch"),
					domain.NewRemoteBranchName("origin/feature-branch"),
				},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.DeleteTrackingBranchStep{
						Branch: domain.NewRemoteBranchName("origin/perennial-branch"),
					},
					&steps.DeleteTrackingBranchStep{
						Branch: domain.NewRemoteBranchName("origin/feature-branch"),
					},
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("perennial-branch"),
						Parent: domain.LocalBranchName{}.Location(),
						Force:  true,
					},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("main")},
					&steps.DeleteLocalBranchStep{
						Branch: domain.NewLocalBranchName("feature-branch"),
						Parent: domain.NewLocalBranchName("main").Location(),
						Force:  true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("main")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch changed locally", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               false,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("feature-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("444444"),
						SetToSHA:    domain.NewSHA("222222"),
						Hard:        true,
					},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("perennial-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("333333"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch remote updated", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
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
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               false,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					// It doesn't reset the remote perennial branch since those are assumed to be protected against force-pushes
					// and we can't revert the commit on it since we cannot change the local perennial branch here.
					&steps.ResetRemoteBranchToSHAStep{
						Branch:      domain.NewRemoteBranchName("origin/feature-branch"),
						SetToSHA:    domain.NewSHA("333333"),
						MustHaveSHA: domain.NewSHA("444444"),
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
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
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:       lineage,
				BranchTypes:   branchTypes,
				InitialBranch: before.Active,
				FinalBranch:   after.Active,
				NoPushHook:    true,
				UndoablePerennialCommits: []domain.SHA{
					domain.NewSHA("444444"),
				},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					// revert the commit on the perennial branch
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("main")},
					&steps.RevertCommitStep{SHA: domain.NewSHA("444444")},
					&steps.PushCurrentBranchStep{CurrentBranch: domain.NewLocalBranchName("main"), NoPushHook: true},
					// reset the feature branch to the previous SHA
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("feature-branch")},
					&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: domain.NewSHA("666666"), SetToSHA: domain.NewSHA("333333"), Hard: true},
					&steps.ForcePushCurrentBranchStep{NoPushHook: true},
					// check out the initial branch
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("upstream commit downloaded and branch shipped at the same time", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
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
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
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
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
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
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:       lineage,
				BranchTypes:   branchTypes,
				InitialBranch: before.Active,
				FinalBranch:   after.Active,
				NoPushHook:    true,
				UndoablePerennialCommits: []domain.SHA{
					domain.NewSHA("444444"),
				},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					// revert the undoable commit on the main branch
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("main")},
					&steps.RevertCommitStep{SHA: domain.NewSHA("444444")},
					&steps.PushCurrentBranchStep{CurrentBranch: domain.NewLocalBranchName("main"), NoPushHook: true},
					// re-create the feature branch
					&steps.CreateBranchStep{Branch: domain.NewLocalBranchName("feature-branch"), StartingPoint: domain.NewSHA("222222").Location()},
					&steps.CreateTrackingBranchStep{Branch: domain.NewLocalBranchName("feature-branch"), NoPushHook: true},
					// check out the initial branch
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
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
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					// It doesn't revert the perennial branch because it cannot force-push the changes to the remote branch.
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("feature-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("555555"),
						SetToSHA:    domain.NewSHA("222222"),
						Hard:        true,
					},
					&steps.ResetRemoteBranchToSHAStep{
						Branch:      domain.NewRemoteBranchName("origin/feature-branch"),
						MustHaveSHA: domain.NewSHA("666666"),
						SetToSHA:    domain.NewSHA("222222"),
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch updates pulled down", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
				domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
			}
			before := domain.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusBehind,
						RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  domain.NewSHA("222222"),
					},
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("333333"),
						SyncStatus: domain.SyncStatusBehind,
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
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("feature-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("444444"),
						SetToSHA:    domain.NewSHA("333333"),
						Hard:        true,
					},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("perennial-branch")},
					&steps.ResetCurrentBranchToSHAStep{
						MustHaveSHA: domain.NewSHA("222222"),
						SetToSHA:    domain.NewSHA("111111"),
						Hard:        true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch updates pushed up", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
				domain.NewLocalBranchName("feature-branch"): domain.NewLocalBranchName("main"),
			}
			before := domain.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("perennial-branch"),
						LocalSHA:   domain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusAhead,
						RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("444444"),
						SyncStatus: domain.SyncStatusAhead,
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
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
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
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					// It doesn't revert the remote perennial branch because it cannot force-push the changes to it.
					&steps.ResetRemoteBranchToSHAStep{
						Branch:      domain.NewRemoteBranchName("origin/feature-branch"),
						MustHaveSHA: domain.NewSHA("444444"),
						SetToSHA:    domain.NewSHA("333333"),
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch deleted locally", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
						SyncStatus: domain.SyncStatusRemoteOnly,
						RemoteName: domain.NewRemoteBranchName("origin/perennial-branch"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
					domain.BranchInfo{
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
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
				RemoteAdded:           []domain.RemoteBranchName{},
				RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					&steps.CreateBranchStep{
						Branch:        domain.NewLocalBranchName("feature-branch"),
						StartingPoint: domain.NewSHA("222222").Location(),
					},
					&steps.CreateBranchStep{
						Branch:        domain.NewLocalBranchName("perennial-branch"),
						StartingPoint: domain.NewSHA("111111").Location(),
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("omnibranch tracking branch deleted", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
			}
			lineage := config.Lineage{
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
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("feature-branch"),
						LocalSHA:   domain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
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
				RemoteAdded:  []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{
					domain.NewRemoteBranchName("origin/perennial-branch"): domain.NewSHA("111111"),
					domain.NewRemoteBranchName("origin/feature-branch"):   domain.NewSHA("222222"),
				},
				RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
				OmniRemoved:           domain.LocalBranchesSHAs{},
				OmniChanged:           domain.LocalBranchChange{},
				InconsistentlyChanged: domain.InconsistentChanges{},
			}
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				List: []steps.Step{
					// don't re-create the tracking branch for the perennial branch
					// because those are protected
					&steps.CreateRemoteBranchStep{
						Branch:     domain.NewLocalBranchName("feature-branch"),
						SHA:        domain.NewSHA("222222"),
						NoPushHook: true,
					},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("feature-branch")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})

		t.Run("sync with a new upstream remote", func(t *testing.T) {
			t.Parallel()
			branchTypes := domain.BranchTypes{
				MainBranch:        domain.NewLocalBranchName("main"),
				PerennialBranches: domain.NewLocalBranchNames(),
			}
			lineage := config.Lineage{
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
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
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
				RemoteAdded: []domain.RemoteBranchName{ // TODO: replace with domain.RemoteBranchNames everywhere
					domain.NewRemoteBranchName("upstream/main"),
				},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{}, // TODO: replace with domain.RemoteBranchesSHAs everywhere
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
			assert.Equal(t, wantChanges, haveChanges)
			haveSteps := haveChanges.UndoSteps(undo.StepsArgs{
				Lineage:                  lineage,
				BranchTypes:              branchTypes,
				InitialBranch:            before.Active,
				FinalBranch:              after.Active,
				NoPushHook:               true,
				UndoablePerennialCommits: []domain.SHA{},
			})
			wantSteps := runstate.StepList{
				// No changes should happen here since all changes were syncs on perennial branches.
				// We don't want to undo these commits because that would undo commits
				// already committed to perennial branches by others for everybody on the team.
				List: []steps.Step{
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("main")},
				},
			}
			assert.Equal(t, wantSteps, haveSteps)
		})
	})
}
