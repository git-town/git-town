package git_test

import (
	"testing"

	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/gohacks"
	"github.com/git-town/git-town/v20/internal/gohacks/cache"
	"github.com/git-town/git-town/v20/internal/subshell"
	"github.com/git-town/git-town/v20/internal/test/testgit"
	"github.com/git-town/git-town/v20/internal/test/testruntime"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBackendCommands(t *testing.T) {
	t.Parallel()
	initial := gitdomain.NewLocalBranchName("initial")

	t.Run("BranchAuthors", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch")
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch,
			FileContent: "file1",
			FileName:    "file1",
			Message:     "first commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch,
			FileContent: "file2",
			FileName:    "file2",
			Message:     "second commit",
		})
		authors, err := runtime.Git.BranchAuthors(runtime.TestRunner, branch, initial)
		must.NoError(t, err)
		must.Eq(t, []gitdomain.Author{"user <email@example.com>"}, authors)
	})

	t.Run("BranchContainsMerges", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has a merge commit", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			err := runtime.Git.CreateAndCheckoutBranch(runtime, branch1)
			must.NoError(t, err)
			runtime.CreateBranch(branch2, branch1.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			runtime.CheckoutBranch(branch2)
			err = runtime.Git.MergeNoFastForward(runtime, configdomain.UseDefaultMessage(), branch1)
			must.NoError(t, err)
			have, err := runtime.Git.BranchContainsMerges(runtime, branch2, branch1)
			must.NoError(t, err)
			must.True(t, have)
		})
		t.Run("branch has no merge commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			err := runtime.Git.CreateAndCheckoutBranch(runtime, branch1)
			must.NoError(t, err)
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			have, err := runtime.Git.BranchContainsMerges(runtime, branch1, initial)
			must.NoError(t, err)
			must.False(t, have)
		})
	})

	t.Run("BranchExists", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch("b1", initial.BranchName())
		runner.CreateBranch("b2", initial.BranchName())
		must.True(t, runner.Git.BranchExists(runner, "b1"))
		must.True(t, runner.Git.BranchExists(runner, "b2"))
		must.False(t, runner.Git.BranchExists(runner, "b3"))
	})

	t.Run("BranchHasUnmergedChanges", func(t *testing.T) {
		t.Parallel()
		t.Run("branch without commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch")
			runtime.CreateBranch(branch, initial.BranchName())
			have, err := runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
			must.NoError(t, err)
			must.False(t, have)
		})
		t.Run("branch with commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(testgit.Commit{
				Branch:      initial,
				FileContent: "original content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			branch := gitdomain.NewLocalBranchName("branch")
			runtime.CreateBranch(branch, initial.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "modified content",
				FileName:    "file1",
				Message:     "commit 2",
			})
			have, err := runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
			must.NoError(t, err)
			must.True(t, have, must.Sprint("branch with commits that make changes"))
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "original content",
				FileName:    "file1",
				Message:     "commit 3",
			})
			have, err = runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
			must.NoError(t, err)
			must.False(t, have, must.Sprint("branch with commits that make no changes"))
		})
	})

	t.Run("BranchInSyncWithParent", func(t *testing.T) {
		t.Parallel()
		t.Run("child has the same commits as parent", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "parent_file",
				Message:     "add parent file",
			})
			err = local.Git.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			inSync, err := local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("parent has extra commit", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "file",
				Message:     "commit on both parent and child",
			})
			err = local.Git.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on parent",
			})
			inSync, err := local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.False(t, inSync)
		})
		t.Run("child has extra commit", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "file",
				Message:     "commit on both parent and child",
			})
			err = local.Git.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "child",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on child",
			})
			inSync, err := local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("empty parent", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			err = local.Git.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "child",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on child",
			})
			inSync, err := local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("both empty", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			err = local.Git.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			inSync, err := local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
	})

	t.Run("BranchInSyncWithTracking", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has no commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			inSync, err := local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("branch has local commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "local_file",
				Message:     "add local file",
			})
			inSync, err := local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.False(t, inSync)
		})
		t.Run("branch has remote commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			origin.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "remote_file",
				Message:     "add remote file",
			})
			local.Fetch()
			inSync, err := local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.False(t, inSync)
		})
		t.Run("branch has different local and remote commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "local_file",
				Message:     "add local file",
			})
			origin.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "remote_file",
				Message:     "add remote file",
			})
			local.Fetch()
			inSync, err := local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.False(t, inSync)
		})
	})

	t.Run("BranchesSnapshot", func(t *testing.T) {
		t.Parallel()
		t.Run("recognizes the active branch names", func(t *testing.T) {
			t.Parallel()
			t.Run("first branch is checked out", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("first-branch", initial.BranchName())
				runtime.CreateBranch("second-branch", initial.BranchName())
				runtime.CheckoutBranch("first-branch")
				snapshot, err := runtime.Git.BranchesSnapshot(runtime)
				must.NoError(t, err)
				must.Eq(t, Some[gitdomain.LocalBranchName]("first-branch"), snapshot.Active)
			})

			t.Run("second branch is checked out", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("first-branch", initial.BranchName())
				runtime.CreateBranch("second-branch", initial.BranchName())
				runtime.CheckoutBranch("second-branch")
				snapshot, err := runtime.Git.BranchesSnapshot(runtime)
				must.NoError(t, err)
				must.Eq(t, Some[gitdomain.LocalBranchName]("second-branch"), snapshot.Active)
			})

			t.Run("in the middle of a rebase", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("branch", initial.BranchName())
				runtime.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "branch content",
					FileName:    "file",
					Message:     "branch commit",
				})
				runtime.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "initial content",
					FileName:    "file",
					Message:     "initial commit",
				})
				runtime.CheckoutBranch("branch")
				err := runtime.RebaseAgainstBranch(initial)
				must.Error(t, err)
				rebaseInProgress, err := runtime.Git.HasRebaseInProgress(runtime)
				must.NoError(t, err)
				must.True(t, rebaseInProgress)
				snapshot, err := runtime.Git.BranchesSnapshot(runtime)
				must.NoError(t, err)
				must.Eq(t, None[gitdomain.LocalBranchName](), snapshot.Active)
			})
		})

		t.Run("recognizes the branch sync status", func(t *testing.T) {
			t.Parallel()
			t.Run("branch is ahead of its remote branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				origin.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content",
					FileName:    "file",
					Message:     "local and origin commit",
				})
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				local.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content 2",
					FileName:    "file",
					Message:     "local commit",
				})
				commits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some(initial),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(commits[0].SHA),
							SyncStatus: gitdomain.SyncStatusAhead,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(commits[1].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("branch is behind its remote branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				origin.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content",
					FileName:    "file",
					Message:     "local and origin commit",
				})
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				origin.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content 2",
					FileName:    "file",
					Message:     "origin commit",
				})
				local.Fetch()
				commits, err := origin.Git.CommitsInBranch(origin, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some(initial),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(commits[1].SHA),
							SyncStatus: gitdomain.SyncStatusBehind,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(commits[0].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("branch is ahead and behind its remote branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				origin.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content",
					FileName:    "file",
					Message:     "local and origin commit",
				})
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				local.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content 2",
					FileName:    "file",
					Message:     "local commit",
				})
				origin.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content 3",
					FileName:    "file",
					Message:     "origin commit",
				})
				local.Fetch()
				originCommits, err := origin.Git.CommitsInBranch(origin, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				localCommits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some(initial),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(localCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusNotInSync,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(originCommits[0].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("branch is in sync", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				origin.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content",
					FileName:    "file",
					Message:     "local and origin commit",
				})
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				commits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some(initial),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(commits[0].SHA),
							SyncStatus: gitdomain.SyncStatusUpToDate,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(commits[0].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("remote-only branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				origin.CreateAndCheckoutFeatureBranch("branch", initial.Location())
				origin.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "content",
					FileName:    "file",
					Message:     "origin commit",
				})
				localCommits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				local.Fetch()
				originBranchCommits, err := origin.Git.CommitsInBranch(origin, "branch", Some(initial))
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some(initial),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(localCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusUpToDate,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(localCommits[0].SHA),
						},
						gitdomain.BranchInfo{
							LocalName:  None[gitdomain.LocalBranchName](),
							LocalSHA:   None[gitdomain.SHA](),
							SyncStatus: gitdomain.SyncStatusRemoteOnly,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
							RemoteSHA:  Some(originBranchCommits[0].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("local-only branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				initialCommits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				local.CreateAndCheckoutFeatureBranch("branch", initial.Location())
				local.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "content",
					FileName:    "file",
					Message:     "local commit",
				})
				localBranchCommits, err := local.Git.CommitsInBranch(local, "branch", Some(initial))
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some[gitdomain.LocalBranchName]("branch"),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some[gitdomain.LocalBranchName]("branch"),
							LocalSHA:   Some(localBranchCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusLocalOnly,
							RemoteName: None[gitdomain.RemoteBranchName](),
							RemoteSHA:  None[gitdomain.SHA](),
						},
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(initialCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusUpToDate,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(initialCommits[0].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("branch is deleted at the remote", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				origin.CreateAndCheckoutFeatureBranch("branch", initial.Location())
				origin.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "content",
					FileName:    "file",
					Message:     "origin commit",
				})
				local.Fetch()
				local.CheckoutBranch("branch")
				origin.CheckoutBranch(initial)
				err := origin.Git.DeleteLocalBranch(origin, "branch")
				must.NoError(t, err)
				local.Fetch()
				initialCommits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				branchCommits, err := local.Git.CommitsInBranch(local, "branch", None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some[gitdomain.LocalBranchName]("branch"),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some[gitdomain.LocalBranchName]("branch"),
							LocalSHA:   Some(branchCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
							RemoteSHA:  None[gitdomain.SHA](),
						},
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(initialCommits[1].SHA),
							SyncStatus: gitdomain.SyncStatusUpToDate,
							RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
							RemoteSHA:  Some(initialCommits[1].SHA),
						},
					},
				}
				have, err := local.Git.BranchesSnapshot(local)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("branch is active in another worktree", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("branch", initial.BranchName())
				worktreeDir := t.TempDir()
				runtime.CreateWorktree(worktreeDir, "branch")
				commits, err := runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: Some[gitdomain.LocalBranchName]("initial"),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
							LocalSHA:   Some(commits[0].SHA),
							SyncStatus: gitdomain.SyncStatusOtherWorktree,
							RemoteName: None[gitdomain.RemoteBranchName](),
							RemoteSHA:  None[gitdomain.SHA](),
						},
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(commits[0].SHA),
							SyncStatus: gitdomain.SyncStatusLocalOnly,
							RemoteName: None[gitdomain.RemoteBranchName](),
							RemoteSHA:  None[gitdomain.SHA](),
						},
					},
				}
				have, err := runtime.Git.BranchesSnapshot(runtime)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})

			t.Run("in the middle of a rebase", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("branch", "initial")
				runtime.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "branch content",
					FileName:    "file",
					Message:     "branch commit",
				})
				runtime.CreateCommit(testgit.Commit{
					Branch:      "initial",
					FileContent: "initial content",
					FileName:    "file",
					Message:     "initial commit",
				})
				runtime.CheckoutBranch("branch")
				err := runtime.RebaseAgainstBranch("initial")
				must.Error(t, err)
				rebaseInProgress, err := runtime.Git.HasRebaseInProgress(runtime)
				must.NoError(t, err)
				must.True(t, rebaseInProgress)
				branchCommits, err := runtime.Git.CommitsInBranch(runtime, "branch", Some(initial))
				must.NoError(t, err)
				initialCommits, err := runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]())
				must.NoError(t, err)
				want := gitdomain.BranchesSnapshot{
					Active: None[gitdomain.LocalBranchName](),
					Branches: gitdomain.BranchInfos{
						gitdomain.BranchInfo{
							LocalName:  Some[gitdomain.LocalBranchName]("branch"),
							LocalSHA:   Some(branchCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusLocalOnly,
							RemoteName: None[gitdomain.RemoteBranchName](),
							RemoteSHA:  None[gitdomain.SHA](),
						},
						gitdomain.BranchInfo{
							LocalName:  Some(initial),
							LocalSHA:   Some(initialCommits[0].SHA),
							SyncStatus: gitdomain.SyncStatusLocalOnly,
							RemoteName: None[gitdomain.RemoteBranchName](),
							RemoteSHA:  None[gitdomain.SHA](),
						},
					},
				}
				have, err := runtime.Git.BranchesSnapshot(runtime)
				must.NoError(t, err)
				must.Eq(t, want, have)
			})
		})

		t.Run("square brackets in the commit message", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			local.CreateBranch("branch-1", initial.BranchName())
			local.CreateCommit(testgit.Commit{
				Branch:      "branch-1",
				FileContent: "content",
				FileName:    "file",
				Message:     "[ci skip] local commit",
			})
			local.CreateBranch("branch-2", initial.BranchName()) // Both local and remote
			local.CreateCommit(testgit.Commit{
				Branch:      "branch-2",
				FileContent: "content",
				FileName:    "file",
				Message:     "[ci skip] local and origin commit",
			})
			local.PushBranchToRemote(gitdomain.NewLocalBranchName("branch-2"), gitdomain.RemoteOrigin)
			origin.CreateBranch("branch-3", initial.BranchName()) // Remote only
			origin.CreateCommit(testgit.Commit{
				Branch:      "branch-3",
				FileContent: "content",
				FileName:    "file",
				Message:     "[ci skip] origin commit",
			})
			local.Fetch()
			branch1Commits, err := local.Git.CommitsInBranch(local, "branch-1", Some(initial))
			must.NoError(t, err)
			branch2Commits, err := local.Git.CommitsInBranch(local, "branch-2", Some(initial))
			must.NoError(t, err)
			initialCommits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			branch3Commits, err := local.Git.CommitsInBranch(local, "origin/branch-3", gitdomain.NewLocalBranchNameOption("origin/initial"))
			must.NoError(t, err)
			want := gitdomain.BranchesSnapshot{
				Active: gitdomain.NewLocalBranchNameOption("branch-2"),
				Branches: gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
						LocalSHA:   Some(branch1Commits[0].SHA),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
					},
					gitdomain.BranchInfo{
						LocalName:  gitdomain.NewLocalBranchNameOption("branch-2"),
						LocalSHA:   Some(branch2Commits[0].SHA),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
						RemoteSHA:  Some(branch2Commits[0].SHA),
					},
					gitdomain.BranchInfo{
						LocalName:  Some(initial),
						LocalSHA:   Some(initialCommits[1].SHA),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
						RemoteSHA:  Some(initialCommits[1].SHA),
					},
					gitdomain.BranchInfo{
						LocalName:  None[gitdomain.LocalBranchName](),
						LocalSHA:   None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-3")),
						RemoteSHA:  Some(branch3Commits[0].SHA),
					},
				},
			}
			have, err := local.Git.BranchesSnapshot(local)
			must.NoError(t, err)
			must.Eq(t, want, have)
		})

		t.Run("ignores symbolic refs", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.Run("git", "symbolic-ref", "refs/remotes/origin/master", "refs/remotes/origin/initial")
			must.NoError(t, err)
			commits, err := local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			want := gitdomain.BranchesSnapshot{
				Active: Some(initial),
				Branches: gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(initial),
						LocalSHA:   Some(commits[0].SHA),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/initial")),
						RemoteSHA:  Some(commits[0].SHA),
					},
				},
			}
			have, err := local.Git.BranchesSnapshot(local)
			must.NoError(t, err)
			must.Eq(t, want, have)
		})
	})

	t.Run("CheckoutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CheckoutBranch(branch)
		currentBranch, err := runtime.Git.CurrentBranch(runtime.TestRunner)
		must.NoError(t, err)
		must.EqOp(t, branch, currentBranch)
		runtime.CheckoutBranch(initial)
		currentBranch, err = runtime.Git.CurrentBranch(runtime.TestRunner)
		must.NoError(t, err)
		must.EqOp(t, initial, currentBranch)
	})

	t.Run("CommitsInBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("feature branch contains commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch1")
			runtime.CreateBranch(branch, initial.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file1",
				Message:  "commit 1",
			})
			runtime.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file2",
				Message:  "commit 2",
			})
			commits, err := runtime.Git.CommitsInBranch(runtime.TestRunner, branch, gitdomain.NewLocalBranchNameOption("initial"))
			must.NoError(t, err)
			haveMessages := commits.Messages()
			wantMessages := gitdomain.NewCommitMessages("commit 1", "commit 2")
			must.Eq(t, wantMessages, haveMessages)
		})
		t.Run("feature branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch1")
			runtime.CreateBranch(branch, initial.BranchName())
			commits, err := runtime.Git.CommitsInBranch(runtime, branch, gitdomain.NewLocalBranchNameOption("initial"))
			must.NoError(t, err)
			must.EqOp(t, 0, len(commits))
		})
		t.Run("main branch contains commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(testgit.Commit{
				Branch:   initial,
				FileName: "file1",
				Message:  "commit 1",
			})
			runtime.CreateCommit(testgit.Commit{
				Branch:   initial,
				FileName: "file2",
				Message:  "commit 2",
			})
			commits, err := runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			must.EqOp(t, 3, len(commits)) // 1 initial commit + 2 test commits
		})
		t.Run("main branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			commits, err := runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			must.EqOp(t, 1, len(commits)) // the initial commit
		})
	})

	t.Run("CommitsInFeatureBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("feature branch contains commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch1")
			runtime.CreateBranch(branch, initial.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file1",
				Message:  "commit 1\n\nbody line 1a\nbody line 1b",
			})
			runtime.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file2",
				Message:  "commit 2\n\nbody line 2a\nbody line 2b",
			})
			commits, err := runtime.Git.CommitsInFeatureBranch(runtime.TestRunner, branch, gitdomain.NewLocalBranchName("initial"))
			must.NoError(t, err)
			haveMessages := commits.Messages()
			// this method returns commit titles only, use CommitMessage() to get the full commit message for a particular commit
			wantMessages := gitdomain.NewCommitMessages("commit 1", "commit 2")
			must.Eq(t, wantMessages, haveMessages)
		})
		t.Run("feature branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch1")
			runtime.CreateBranch(branch, initial.BranchName())
			commits, err := runtime.Git.CommitsInFeatureBranch(runtime, branch, gitdomain.NewLocalBranchName("initial"))
			must.NoError(t, err)
			must.EqOp(t, 0, len(commits))
		})
	})

	t.Run("CurrentBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CheckoutBranch(initial)
		branch := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CheckoutBranch(branch)
		branch, err := runtime.Git.CurrentBranch(runtime)
		must.NoError(t, err)
		must.EqOp(t, branch, branch)
		runtime.CheckoutBranch(initial)
		branch, err = runtime.Git.CurrentBranch(runtime)
		must.NoError(t, err)
		must.EqOp(t, initial, branch)
	})

	t.Run("CurrentBranchDuringRebase", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch", "initial")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch",
			FileContent: "branch content",
			FileName:    "file",
			Message:     "branch commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "initial content",
			FileName:    "file",
			Message:     "initial commit",
		})
		runtime.CheckoutBranch("branch")
		err := runtime.RebaseAgainstBranch("initial")
		must.Error(t, err)
		rebaseInProgress, err := runtime.Git.HasRebaseInProgress(runtime)
		must.NoError(t, err)
		must.True(t, rebaseInProgress)
		have, err := runtime.Git.CurrentBranchDuringRebase(runtime)
		must.NoError(t, err)
		must.EqOp(t, "branch", have)
	})

	t.Run("CurrentBranchHasTrackingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has tracking branch", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.CreateGitTown(t)
			repoDir := t.TempDir()
			repo := testruntime.Clone(origin.TestRunner, repoDir)
			branch := gitdomain.NewLocalBranchName("branch")
			main := gitdomain.NewLocalBranchName("main")
			repo.CheckoutBranch(main)
			repo.CreateAndCheckoutFeatureBranch(branch, main.Location())
			repo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
			have := repo.Git.CurrentBranchHasTrackingBranch(repo)
			must.True(t, have)
		})
		t.Run("has no tracking branch", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.CreateGitTown(t)
			repoDir := t.TempDir()
			repo := testruntime.Clone(origin.TestRunner, repoDir)
			branch := gitdomain.NewLocalBranchName("branch")
			main := gitdomain.NewLocalBranchName("main")
			repo.CheckoutBranch(main)
			repo.CreateAndCheckoutFeatureBranch(branch, main.Location())
			have := repo.Git.CurrentBranchHasTrackingBranch(repo)
			must.False(t, have)
		})
	})

	t.Run("DefaultBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.SetDefaultGitBranch("main")
		have := runtime.Git.DefaultBranch(runtime)
		want := gitdomain.NewLocalBranchNameOption("main")
		must.Eq(t, want, have)
	})

	t.Run("DetectPhantomMergeConflicts", func(t *testing.T) {
		t.Parallel()
		t.Run("legit phantom merge conflict", func(t *testing.T) {
			t.Parallel()
			fullInfos := []git.FileConflictFullInfo{
				{
					Main: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Parent: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Current: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(fullInfos, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomMergeConflict{
				{FilePath: "file"},
			}
			must.Eq(t, want, have)
		})
		t.Run("permissions differ", func(t *testing.T) {
			t.Parallel()
			fullInfos := []git.FileConflictFullInfo{
				{
					Main: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Parent: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100644",
						SHA:        "111111",
					}),
					Current: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(fullInfos, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomMergeConflict{}
			must.Eq(t, want, have)
		})
		t.Run("file checksums between parent and main differ", func(t *testing.T) {
			t.Parallel()
			fullInfos := []git.FileConflictFullInfo{
				{
					Main: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Parent: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100644",
						SHA:        "222222",
					}),
					Current: Some(git.BlobInfo{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "222222",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(fullInfos, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomMergeConflict{}
			must.Eq(t, want, have)
		})
		t.Run("file names between parent and main differ", func(t *testing.T) {
			t.Parallel()
			fullInfos := []git.FileConflictFullInfo{
				{
					Main: Some(git.BlobInfo{
						FilePath:   "file-1",
						Permission: "100755",
						SHA:        "222222",
					}),
					Parent: Some(git.BlobInfo{
						FilePath:   "file-2",
						Permission: "100755",
						SHA:        "111111",
					}),
					Current: Some(git.BlobInfo{
						FilePath:   "file-2",
						Permission: "100755",
						SHA:        "111111",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(fullInfos, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomMergeConflict{}
			must.Eq(t, want, have)
		})
	})

	t.Run("FirstCommitMessageInBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is empty", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			must.NoError(t, repo.Git.CreateAndCheckoutBranch(repo.TestRunner, "branch"))
			have, err := repo.Git.FirstCommitMessageInBranch(repo.TestRunner, "branch", "main")
			must.NoError(t, err)
			must.Eq(t, None[gitdomain.CommitMessage](), have)
		})
		t.Run("branch has one commit", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			branch := gitdomain.NewLocalBranchName("branch")
			main := gitdomain.NewLocalBranchName("main")
			repo.CreateFeatureBranch(branch, main.BranchName())
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file",
				Message:  "my commit message",
			})
			have, err := repo.Git.FirstCommitMessageInBranch(repo.TestRunner, branch.BranchName(), main.BranchName())
			must.NoError(t, err)
			want := Some(gitdomain.CommitMessage("my commit message"))
			must.Eq(t, want, have)
		})
		t.Run("branch has multiple commits", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			branch := gitdomain.NewLocalBranchName("branch")
			main := gitdomain.NewLocalBranchName("main")
			repo.CreateFeatureBranch(branch, main.BranchName())
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file_1",
				Message:  "commit message 1",
			})
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file_2",
				Message:  "commit message 2",
			})
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file_3",
				Message:  "commit message 3",
			})
			have, err := repo.Git.FirstCommitMessageInBranch(repo.TestRunner, branch.BranchName(), main.BranchName())
			must.NoError(t, err)
			want := Some(gitdomain.CommitMessage("commit message 1"))
			must.Eq(t, want, have)
		})
		t.Run("branch doesn't exist", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			_, err := repo.Git.FirstCommitMessageInBranch(repo.TestRunner, "zonk", "main")
			must.Error(t, err)
		})
		t.Run("branch exists only at the remote", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.CreateGitTown(t)
			repoDir := t.TempDir()
			repo := testruntime.Clone(origin.TestRunner, repoDir)
			branch := gitdomain.NewLocalBranchName("branch")
			main := gitdomain.NewRemoteBranchName("origin/main")
			repo.CreateFeatureBranch(branch, main.BranchName())
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file_1",
				Message:  "commit message 1",
			})
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file_2",
				Message:  "commit message 2",
			})
			repo.CreateCommit(testgit.Commit{
				Branch:   branch,
				FileName: "file_3",
				Message:  "commit message 3",
			})
			repo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
			repo.CheckoutBranch(main.LocalBranchName())
			err := repo.Git.DeleteLocalBranch(repo.TestRunner, branch)
			must.NoError(t, err)
			have, err := repo.Git.FirstCommitMessageInBranch(repo.TestRunner, branch.TrackingBranch(gitdomain.RemoteOrigin).BranchName(), main.BranchName())
			must.NoError(t, err)
			want := Some(gitdomain.CommitMessage("commit message 1"))
			must.Eq(t, want, have)
		})
	})

	t.Run("FirstExistingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("first branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			runtime.CreateBranch(branch1, initial.BranchName())
			runtime.CreateBranch(branch2, initial.BranchName())
			have := runtime.Git.FirstExistingBranch(runtime, branch1, branch2)
			want := Some(branch1)
			must.Eq(t, want, have)
		})
		t.Run("second branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			runtime.CreateBranch(branch2, initial.BranchName())
			have := runtime.Git.FirstExistingBranch(runtime, branch1, branch2)
			want := Some(branch2)
			must.Eq(t, want, have)
		})
		t.Run("no branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			have := runtime.Git.FirstExistingBranch(runtime, branch1, branch2)
			want := None[gitdomain.LocalBranchName]()
			must.EqOp(t, want, have)
		})
	})

	t.Run("MergeFastForward", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch")
		runtime := testruntime.Create(t)
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch,
			FileContent: "file1",
			FileName:    "file1",
			Message:     "first commit",
		})
		runtime.CheckoutBranch(initial) // CreateCommit checks out `branch`, go back to `initial`.
		err := runtime.Git.MergeFastForward(runtime.TestRunner, branch.BranchName())
		must.NoError(t, err)
		commits, err := runtime.Git.CommitsInPerennialBranch(runtime) // Current branch.
		must.NoError(t, err)
		haveMessages := commits.Messages()
		wantMessages := gitdomain.NewCommitMessages("first commit", "initial commit")
		must.Eq(t, wantMessages, haveMessages)
	})

	t.Run("MergeNoFastForward", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch")
		runtime := testruntime.Create(t)
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch,
			FileContent: "file1",
			FileName:    "file1",
			Message:     "first commit",
		})
		runtime.CheckoutBranch(initial) // CreateCommit checks out `branch`, go back to `initial`.
		err := runtime.Git.MergeNoFastForward(runtime.TestRunner, configdomain.UseDefaultMessage(), branch)
		must.NoError(t, err)
		commits, err := runtime.Git.CommitsInPerennialBranch(runtime) // Current branch.
		must.NoError(t, err)
		haveMessages := commits.Messages()
		wantMessages := gitdomain.NewCommitMessages("Merge branch 'branch' into initial", "initial commit", "first commit")
		must.SliceContainsAll(t, wantMessages, haveMessages)
	})

	t.Run("MergeNoFastForwardWithCommitMessage", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch")
		runtime := testruntime.Create(t)
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch,
			FileContent: "file1",
			FileName:    "file1",
			Message:     "first commit",
		})
		mergeMessage := gitdomain.CommitMessage("merge message")
		runtime.CheckoutBranch(initial) // CreateCommit checks out `branch`, go back to `initial`.
		err := runtime.Git.MergeNoFastForward(runtime.TestRunner, configdomain.UseCustomMessage(mergeMessage), branch)
		must.NoError(t, err)
		commits, err := runtime.Git.CommitsInPerennialBranch(runtime) // Current branch.
		must.NoError(t, err)
		haveMessages := commits.Messages()
		wantMessages := gitdomain.NewCommitMessages("merge message", "initial commit", "first commit")
		must.SliceContainsAll(t, wantMessages, haveMessages)
	})

	t.Run("NewUnmergedStage", func(t *testing.T) {
		t.Parallel()
		tests := map[int]git.UnmergedStage{
			1: git.UnmergedStageBase,
			2: git.UnmergedStageCurrentBranch,
			3: git.UnmergedStageIncoming,
		}
		for give, want := range tests {
			have, err := git.NewUnmergedStage(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("PreviouslyCheckedOutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("feature1", initial.BranchName())
		runtime.CreateBranch("feature2", initial.BranchName())
		runtime.CheckoutBranch("feature1")
		runtime.CheckoutBranch("feature2")
		have := runtime.Git.PreviouslyCheckedOutBranch(runtime.TestRunner)
		must.Eq(t, gitdomain.NewLocalBranchNameOption("feature1"), have)
	})

	t.Run("Remotes", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		runtime.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		remotes, err := runtime.Git.Remotes(runtime.TestRunner)
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{gitdomain.RemoteOrigin}, remotes)
	})

	t.Run("RepoStatus", func(t *testing.T) {
		t.Run("OpenChanges", func(t *testing.T) {
			t.Parallel()
			t.Run("no open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				have, err := runtime.Git.RepoStatus(runtime)
				must.NoError(t, err)
				must.False(t, have.OpenChanges)
			})
			t.Run("has open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateFile("foo", "bar")
				have, err := runtime.Git.RepoStatus(runtime)
				must.NoError(t, err)
				must.True(t, have.OpenChanges)
			})
			t.Run("during rebase", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				branch1 := gitdomain.NewLocalBranchName("branch1")
				runtime.CreateBranch(branch1, initial.BranchName())
				runtime.CheckoutBranch(branch1)
				runtime.CreateCommit(testgit.Commit{
					Branch:      branch1,
					FileContent: "content on branch1",
					FileName:    "file",
					Message:     "Create file",
				})
				runtime.CheckoutBranch(initial)
				runtime.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content on initial",
					FileName:    "file",
					Message:     "Create file1",
				})
				_ = runtime.RebaseAgainstBranch(branch1) // this is expected to fail
				have, err := runtime.Git.RepoStatus(runtime)
				must.NoError(t, err)
				must.False(t, have.OpenChanges)
			})
			t.Run("during merge conflict", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				branch1 := gitdomain.NewLocalBranchName("branch1")
				runtime.CreateBranch(branch1, initial.BranchName())
				runtime.CheckoutBranch(branch1)
				runtime.CreateCommit(testgit.Commit{
					Branch:      branch1,
					FileContent: "content on branch1",
					FileName:    "file",
					Message:     "Create file",
				})
				runtime.CheckoutBranch(initial)
				runtime.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "content on initial",
					FileName:    "file",
					Message:     "Create file1",
				})
				_ = runtime.MergeBranch(branch1) // this is expected to fail
				have, err := runtime.Git.RepoStatus(runtime)
				must.NoError(t, err)
				must.False(t, have.OpenChanges)
			})
			t.Run("unstashed conflicting changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateFile("file", "stashed content")
				runtime.StashOpenFiles()
				runtime.CreateCommit(testgit.Commit{
					Branch:      initial,
					FileContent: "committed content",
					FileName:    "file",
					Message:     "Create file",
				})
				_ = runtime.UnstashOpenFiles() // this is expected to fail
				have, err := runtime.Git.RepoStatus(runtime)
				must.NoError(t, err)
				must.True(t, have.OpenChanges)
			})

			t.Run("status.short enabled", func(t *testing.T) {
				t.Parallel()
				t.Run("no open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					err := runtime.Run("git", "config", "status.short", "true")
					must.NoError(t, err)
					have, err := runtime.Git.RepoStatus(runtime)
					must.NoError(t, err)
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					err := runtime.Run("git", "config", "status.short", "true")
					must.NoError(t, err)
					have, err := runtime.Git.RepoStatus(runtime)
					must.NoError(t, err)
					must.True(t, have.OpenChanges)
				})
			})

			t.Run("status.branch enabled", func(t *testing.T) {
				t.Parallel()
				t.Run("no open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					err := runtime.Run("git", "config", "status.branch", "true")
					must.NoError(t, err)
					have, err := runtime.Git.RepoStatus(runtime)
					must.NoError(t, err)
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					err := runtime.Run("git", "config", "status.branch", "true")
					must.NoError(t, err)
					have, err := runtime.Git.RepoStatus(runtime)
					must.NoError(t, err)
					must.True(t, have.OpenChanges)
				})
			})
		})

		t.Run("RebaseInProgress", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have, err := runtime.Git.RepoStatus(runtime)
			must.NoError(t, err)
			must.False(t, have.RebaseInProgress)
		})
	})

	t.Run("RootDirectory", func(t *testing.T) {
		t.Parallel()
		t.Run("inside a Git repo", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have := runtime.Git.RootDirectory(runtime.TestRunner)
			must.True(t, have.IsSome())
		})
		t.Run("outside a Git repo", func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			runner := subshell.BackendRunner{
				Dir:             Some(dir),
				Verbose:         false,
				CommandsCounter: NewMutable(new(gohacks.Counter)),
			}
			cmds := git.Commands{
				CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
				RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
			}
			have := cmds.RootDirectory(runner)
			must.True(t, have.IsNone())
		})
	})

	t.Run("ShortenSHA", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch")
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch,
			FileContent: "file1",
			FileName:    "file1",
			Message:     "first commit",
		})
		commits, err := runtime.Git.CommitsInBranch(runtime.TestRunner, "branch", gitdomain.NewLocalBranchNameOption("initial"))
		must.NoError(t, err)
		have, err := runtime.Git.ShortenSHA(runtime, commits[0].SHA)
		must.NoError(t, err)
		must.True(t, len(commits[0].SHA.String()) == 40)
		must.True(t, len(have.String()) == 7)
		must.EqOp(t, have.String(), commits[0].SHA.String()[:7])
	})

	t.Run("StashEntries", func(t *testing.T) {
		t.Parallel()
		t.Run("some stash entries", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("file1", "content")
			runtime.StashOpenFiles()
			runtime.CreateFile("file2", "content")
			runtime.StashOpenFiles()
			have, err := runtime.Git.StashSize(runtime.TestRunner)
			want := gitdomain.StashSize(2)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		})
		t.Run("no stash entries", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have, err := runtime.Git.StashSize(runtime.TestRunner)
			want := gitdomain.StashSize(0)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		})
	})

	t.Run("lastBranchInRef", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"refs/remotes/origin/main": "main",
			"":                         "",
		}
		for give, want := range tests {
			have := git.LastBranchInRef(give)
			must.EqOp(t, want, have)
		}
	})
}
