package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	"github.com/shoenig/test/must"
)

func TestBranchSpans(t *testing.T) {
	t.Parallel()
	// BranchSpans are tested in branch_changes_test.go

	t.Run("RemoveRemoteOnlyBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("removes remote-only branches that got updated", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  None[gitdomain.LocalBranchName](),
						LocalSHA:   None[gitdomain.SHA](),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  None[gitdomain.LocalBranchName](),
						LocalSHA:   None[gitdomain.SHA](),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
				},
			}
			want := undobranches.BranchSpans{}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("removes remote-only branches that got added", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: None[gitdomain.BranchInfo](),
					After: Some(gitdomain.BranchInfo{
						LocalName:  None[gitdomain.LocalBranchName](),
						LocalSHA:   None[gitdomain.SHA](),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
				},
			}
			want := undobranches.BranchSpans{}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("removes remote-only branches that got removed", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  None[gitdomain.LocalBranchName](),
						LocalSHA:   None[gitdomain.SHA](),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
					After: None[gitdomain.BranchInfo](),
				},
			}
			want := undobranches.BranchSpans{}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("keeps local branches that got changed", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("222222")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
				},
			}
			want := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("222222")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
				},
			}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("keeps omni-branches that got changed", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("222222")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
				},
			}
			want := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("222222")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
					}),
				},
			}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("keeps local branches that got their tracking branches created", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
					}),
				},
			}
			want := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					}),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
					}),
				},
			}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("keeps local branches that got added", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: None[gitdomain.BranchInfo](),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					}),
				},
			}
			want := undobranches.BranchSpans{
				{
					Before: None[gitdomain.BranchInfo](),
					After: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					}),
				},
			}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})

		t.Run("keeps local branches that got removed", func(t *testing.T) {
			give := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					}),
					After: None[gitdomain.BranchInfo](),
				},
			}
			want := undobranches.BranchSpans{
				{
					Before: Some(gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
					}),
					After: None[gitdomain.BranchInfo](),
				},
			}
			have := give.RemoveRemoteOnlyBranches()
			must.Eq(t, want, have)
		})
	})
}
