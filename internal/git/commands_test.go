package git_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/cache"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/test/testgit"
	"github.com/git-town/git-town/v22/internal/test/testruntime"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
		authors := asserts.NoError1(runtime.Git.BranchAuthors(runtime.TestRunner, branch, initial))
		must.Eq(t, []gitdomain.Author{"user <email@example.com>"}, authors)
	})

	t.Run("BranchContainsMerges", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has a merge commit", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch2 := gitdomain.NewLocalBranchName("branch-2")
			asserts.NoError(runtime.Git.CreateAndCheckoutBranch(runtime, branch1))
			runtime.CreateBranch(branch2, branch1.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			runtime.CheckoutBranch(branch2)
			asserts.NoError(runtime.Git.MergeNoFastForward(runtime, configdomain.UseDefaultMessage(), branch1))
			have := asserts.NoError1(runtime.Git.BranchContainsMerges(runtime, branch2, branch1))
			must.True(t, have)
		})
		t.Run("branch has no merge commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			asserts.NoError(runtime.Git.CreateAndCheckoutBranch(runtime, branch1))
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			have := asserts.NoError1(runtime.Git.BranchContainsMerges(runtime, branch1, initial))
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
			have := asserts.NoError1(runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial))
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
			have := asserts.NoError1(runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial))
			must.True(t, have)
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "original content",
				FileName:    "file1",
				Message:     "commit 3",
			})
			have = asserts.NoError1(runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial))
			must.False(t, have)
		})
		t.Run("branch with same name as folder", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(testgit.Commit{
				Branch:      initial,
				FileContent: "content 1",
				FileName:    "test/file1",
				Message:     "commit 1",
			})
			branch := gitdomain.NewLocalBranchName("test")
			runtime.CreateBranch(branch, initial.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "content 2",
				FileName:    "file2",
				Message:     "commit 2",
			})
			have := asserts.NoError1(runtime.Git.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial))
			must.True(t, have, must.Sprint("branch with commits that make changes"))
		})
	})

	t.Run("BranchInSyncWithParent", func(t *testing.T) {
		t.Parallel()
		t.Run("child has the same commits as parent", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "parent_file",
				Message:     "add parent file",
			})
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.True(t, inSync)
		})
		t.Run("parent has extra commit", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "file",
				Message:     "commit on both parent and child",
			})
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on parent",
			})
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.False(t, inSync)
		})
		t.Run("child has extra commit", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "file",
				Message:     "commit on both parent and child",
			})
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			local.CreateCommit(testgit.Commit{
				Branch:      "child",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on child",
			})
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.True(t, inSync)
		})
		t.Run("empty parent", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			local.CreateCommit(testgit.Commit{
				Branch:      "child",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on child",
			})
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.True(t, inSync)
		})
		t.Run("both empty", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.True(t, inSync)
		})
		t.Run("child amends a commit from the parent", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			must.NoError(t, local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "parent content",
				FileName:    "file",
				Message:     "parent adds file",
			})
			must.NoError(t, local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			local.CreateFile("file", "child content")
			local.StageFiles("file")
			local.AmendCommit()
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.False(t, inSync)
		})
		t.Run("parent amends a commit", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			must.NoError(t, local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "parent content",
				FileName:    "file",
				Message:     "parent adds file",
			})
			must.NoError(t, local.Git.CreateBranch(local.TestRunner, "child", "parent"))
			local.CreateFile("file", "amended content")
			local.StageFiles("file")
			local.AmendCommit()
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.False(t, inSync)
		})
		t.Run("parent gets rebased", func(t *testing.T) {
			t.Parallel()
			local := testruntime.Create(t)
			must.NoError(t, local.Git.CreateAndCheckoutBranch(local.TestRunner, "parent"))
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "parent content",
				FileName:    "parent_file",
				Message:     "parent commit",
			})
			must.NoError(t, local.Git.CreateAndCheckoutBranch(local.TestRunner, "child"))
			local.CreateCommit(testgit.Commit{
				Branch:      "initial",
				FileContent: "initial content",
				FileName:    "initial_file",
				Message:     "initial commit",
			})
			local.CheckoutBranch("parent")
			asserts.NoError(local.RebaseAgainstBranch("initial"))
			inSync := asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "parent", "initial"))
			must.True(t, inSync)
			local.CheckoutBranch("child")
			inSync = asserts.NoError1(local.Git.BranchInSyncWithParent(local.TestRunner, "child", "parent"))
			must.False(t, inSync)
		})
	})

	t.Run("BranchInSyncWithTracking", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has no commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch"))
			asserts.NoError(local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false))
			inSync := asserts.NoError1(local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin))
			must.True(t, inSync)
		})
		t.Run("branch has local commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch"))
			asserts.NoError(local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false))
			local.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "local_file",
				Message:     "add local file",
			})
			inSync := asserts.NoError1(local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin))
			must.False(t, inSync)
		})
		t.Run("branch has remote commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch"))
			asserts.NoError(local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false))
			origin.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "remote_file",
				Message:     "add remote file",
			})
			local.Fetch()
			inSync := asserts.NoError1(local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin))
			must.False(t, inSync)
		})
		t.Run("branch has different local and remote commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			asserts.NoError(local.Git.CreateAndCheckoutBranch(local.TestRunner, "branch"))
			asserts.NoError(local.Git.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false))
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
			inSync := asserts.NoError1(local.Git.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin))
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
				snapshot := asserts.NoError1(runtime.Git.BranchesSnapshot(runtime))
				must.True(t, snapshot.Active.EqualSome("first-branch"))
			})

			t.Run("second branch is checked out", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("first-branch", initial.BranchName())
				runtime.CreateBranch("second-branch", initial.BranchName())
				runtime.CheckoutBranch("second-branch")
				snapshot := asserts.NoError1(runtime.Git.BranchesSnapshot(runtime))
				must.True(t, snapshot.Active.EqualSome("second-branch"))
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
				rebaseInProgress := asserts.NoError1(runtime.Git.HasRebaseInProgress(runtime))
				must.True(t, rebaseInProgress)
				snapshot := asserts.NoError1(runtime.Git.BranchesSnapshot(runtime))
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
				commits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
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
				commits := asserts.NoError1(origin.Git.CommitsInBranch(origin, initial, None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
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
				originCommits := asserts.NoError1(origin.Git.CommitsInBranch(origin, initial, None[gitdomain.LocalBranchName]()))
				localCommits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
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
				commits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
				must.Eq(t, want, have)
			})

			t.Run("remote-only branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				origin.CreateAndCheckoutFeatureBranch("branch", initial)
				origin.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "content",
					FileName:    "file",
					Message:     "origin commit",
				})
				localCommits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
				local.Fetch()
				originBranchCommits := asserts.NoError1(origin.Git.CommitsInBranch(origin, "branch", Some(initial)))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
				must.Eq(t, want, have)
			})

			t.Run("local-only branch", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				initialCommits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
				local.CreateAndCheckoutFeatureBranch("branch", initial)
				local.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "content",
					FileName:    "file",
					Message:     "local commit",
				})
				localBranchCommits := asserts.NoError1(local.Git.CommitsInBranch(local, "branch", Some(initial)))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
				must.Eq(t, want, have)
			})

			t.Run("branch is deleted at the remote", func(t *testing.T) {
				t.Parallel()
				origin := testruntime.Create(t)
				local := testruntime.Clone(origin.TestRunner, t.TempDir())
				origin.CreateAndCheckoutFeatureBranch("branch", initial)
				origin.CreateCommit(testgit.Commit{
					Branch:      "branch",
					FileContent: "content",
					FileName:    "file",
					Message:     "origin commit",
				})
				local.Fetch()
				local.CheckoutBranch("branch")
				origin.CheckoutBranch(initial)
				asserts.NoError(origin.Git.DeleteLocalBranch(origin, "branch"))
				local.Fetch()
				initialCommits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
				branchCommits := asserts.NoError1(local.Git.CommitsInBranch(local, "branch", None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(local.Git.BranchesSnapshot(local))
				must.Eq(t, want, have)
			})

			t.Run("branch is active in another worktree", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateBranch("branch", initial.BranchName())
				worktreeDir := t.TempDir()
				runtime.AddWorktree(worktreeDir, "branch")
				commits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(runtime.Git.BranchesSnapshot(runtime))
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
				rebaseInProgress := asserts.NoError1(runtime.Git.HasRebaseInProgress(runtime))
				must.True(t, rebaseInProgress)
				branchCommits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime, "branch", Some(initial)))
				initialCommits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]()))
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
				have := asserts.NoError1(runtime.Git.BranchesSnapshot(runtime))
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
			branch1Commits := asserts.NoError1(local.Git.CommitsInBranch(local, "branch-1", Some(initial)))
			branch2Commits := asserts.NoError1(local.Git.CommitsInBranch(local, "branch-2", Some(initial)))
			initialCommits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
			branch3Commits := asserts.NoError1(local.Git.CommitsInBranch(local, "origin/branch-3", gitdomain.NewLocalBranchNameOption("origin/initial")))
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
			have := asserts.NoError1(local.Git.BranchesSnapshot(local))
			must.Eq(t, want, have)
		})

		t.Run("ignores symbolic refs", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			asserts.NoError(local.Run("git", "symbolic-ref", "refs/remotes/origin/master", "refs/remotes/origin/initial"))
			commits := asserts.NoError1(local.Git.CommitsInBranch(local, initial, None[gitdomain.LocalBranchName]()))
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
			have := asserts.NoError1(local.Git.BranchesSnapshot(local))
			must.Eq(t, want, have)
		})
	})

	t.Run("CheckoutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CheckoutBranch(branch)
		currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime.TestRunner)).GetOrPanic()
		must.EqOp(t, branch, currentBranch)
		runtime.CheckoutBranch(initial)
		currentBranch = asserts.NoError1(runtime.Git.CurrentBranch(runtime.TestRunner)).GetOrPanic()
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
			commits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime.TestRunner, branch, gitdomain.NewLocalBranchNameOption("initial")))
			haveMessages := commits.Messages()
			wantMessages := gitdomain.NewCommitMessages("commit 1", "commit 2")
			must.Eq(t, wantMessages, haveMessages)
		})
		t.Run("feature branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch1")
			runtime.CreateBranch(branch, initial.BranchName())
			commits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime, branch, gitdomain.NewLocalBranchNameOption("initial")))
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
			commits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]()))
			must.EqOp(t, 3, len(commits)) // 1 initial commit + 2 test commits
		})
		t.Run("main branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			commits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]()))
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
			commits := asserts.NoError1(runtime.Git.CommitsInFeatureBranch(runtime.TestRunner, branch, gitdomain.NewBranchName("initial")))
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
			commits := asserts.NoError1(runtime.Git.CommitsInFeatureBranch(runtime, branch, gitdomain.NewBranchName("initial")))
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
		branch = asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
		must.EqOp(t, branch, branch)
		runtime.CheckoutBranch(initial)
		branch = asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
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
		rebaseInProgress := asserts.NoError1(runtime.Git.HasRebaseInProgress(runtime))
		must.True(t, rebaseInProgress)
		have := asserts.NoError1(runtime.Git.CurrentBranchDuringRebase(runtime)).GetOrPanic()
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
			repo.CreateAndCheckoutFeatureBranch(branch, main)
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
			repo.CreateAndCheckoutFeatureBranch(branch, main)
			have := repo.Git.CurrentBranchHasTrackingBranch(repo)
			must.False(t, have)
		})
	})

	t.Run("DefaultBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.SetDefaultGitBranch("main")
		have := gitconfig.DefaultBranch(runtime)
		want := gitdomain.NewLocalBranchNameOption("main")
		must.Eq(t, want, have)
	})

	t.Run("DetectPhantomMergeConflicts", func(t *testing.T) {
		t.Parallel()
		t.Run("legit phantom merge conflict", func(t *testing.T) {
			t.Parallel()
			mergeConflicts := []git.MergeConflict{
				{
					Root: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Parent: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Current: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(mergeConflicts, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomConflict{
				{
					FilePath:   "file",
					Resolution: gitdomain.ConflictResolutionOurs,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("permissions differ", func(t *testing.T) {
			t.Parallel()
			mergeConflicts := []git.MergeConflict{
				{
					Root: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Parent: Some(git.Blob{
						FilePath:   "file",
						Permission: "100644",
						SHA:        "111111",
					}),
					Current: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(mergeConflicts, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomConflict{}
			must.Eq(t, want, have)
		})
		t.Run("file checksums between parent and main differ", func(t *testing.T) {
			t.Parallel()
			mergeConflicts := []git.MergeConflict{
				{
					Root: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "111111",
					}),
					Parent: Some(git.Blob{
						FilePath:   "file",
						Permission: "100644",
						SHA:        "222222",
					}),
					Current: Some(git.Blob{
						FilePath:   "file",
						Permission: "100755",
						SHA:        "222222",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(mergeConflicts, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomConflict{}
			must.Eq(t, want, have)
		})
		t.Run("file names between parent and main differ", func(t *testing.T) {
			t.Parallel()
			mergeConflicts := []git.MergeConflict{
				{
					Root: Some(git.Blob{
						FilePath:   "file-1",
						Permission: "100755",
						SHA:        "222222",
					}),
					Parent: Some(git.Blob{
						FilePath:   "file-2",
						Permission: "100755",
						SHA:        "111111",
					}),
					Current: Some(git.Blob{
						FilePath:   "file-2",
						Permission: "100755",
						SHA:        "111111",
					}),
				},
			}
			have := git.DetectPhantomMergeConflicts(mergeConflicts, gitdomain.NewLocalBranchNameOption("alpha"), "main")
			want := []git.PhantomConflict{}
			must.Eq(t, want, have)
		})
	})

	t.Run("FirstCommitMessageInBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is empty", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			must.NoError(t, repo.Git.CreateAndCheckoutBranch(repo.TestRunner, "branch"))
			have := asserts.NoError1(repo.Git.FirstCommitMessageInBranch(repo.TestRunner, "branch", "main"))
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
			have := asserts.NoError1(repo.Git.FirstCommitMessageInBranch(repo.TestRunner, branch.BranchName(), main.BranchName()))
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
			have := asserts.NoError1(repo.Git.FirstCommitMessageInBranch(repo.TestRunner, branch.BranchName(), main.BranchName()))
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
			asserts.NoError(repo.Git.DeleteLocalBranch(repo.TestRunner, branch))
			have := asserts.NoError1(repo.Git.FirstCommitMessageInBranch(repo.TestRunner, branch.TrackingBranch(gitdomain.RemoteOrigin).BranchName(), main.BranchName()))
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
		asserts.NoError(runtime.Git.MergeFastForward(runtime.TestRunner, branch.BranchName()))
		commits := asserts.NoError1(runtime.Git.CommitsInPerennialBranch(runtime)) // Current branch.
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
		asserts.NoError(runtime.Git.MergeNoFastForward(runtime.TestRunner, configdomain.UseDefaultMessage(), branch))
		commits := asserts.NoError1(runtime.Git.CommitsInPerennialBranch(runtime)) // Current branch.
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
		asserts.NoError(runtime.Git.MergeNoFastForward(runtime.TestRunner, configdomain.UseCustomMessage(mergeMessage), branch))
		commits := asserts.NoError1(runtime.Git.CommitsInPerennialBranch(runtime)) // Current branch.
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
			have := asserts.NoError1(git.NewUnmergedStage(give))
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

	t.Run("RebaseInProgress", func(t *testing.T) {
		t.Parallel()
		t.Run("not in progress", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			repoStatus := asserts.NoError1(runtime.Git.RepoStatus(runtime))
			must.False(t, repoStatus.RebaseInProgress)
		})
		t.Run("in progress", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("branch1")
			runtime.CreateBranch(branch1, initial.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content 1",
				FileName:    "file",
				Message:     "commit 1",
			})
			branch2 := gitdomain.NewLocalBranchName("branch2")
			runtime.CreateBranch(branch2, initial.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch2,
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit 2",
			})
			runtime.CheckoutBranch(branch2)
			err := runtime.RebaseAgainstBranch(branch1)
			must.Error(t, err)
			repoStatus := asserts.NoError1(runtime.Git.RepoStatus(runtime))
			must.True(t, repoStatus.RebaseInProgress)
		})
	})

	t.Run("Remotes", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		runtime.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		remotes := asserts.NoError1(runtime.Git.Remotes(runtime.TestRunner))
		must.Eq(t, gitdomain.Remotes{gitdomain.RemoteOrigin}, remotes)
	})

	t.Run("RepoStatus", func(t *testing.T) {
		t.Run("OpenChanges", func(t *testing.T) {
			t.Parallel()
			t.Run("no open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
				must.False(t, have.OpenChanges)
			})
			t.Run("has open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateFile("foo", "bar")
				have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
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
				have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
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
				have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
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
				have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
				must.True(t, have.OpenChanges)
			})

			t.Run("status.short enabled", func(t *testing.T) {
				t.Parallel()
				t.Run("no open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					asserts.NoError(runtime.Run("git", "config", "status.short", "true"))
					have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					asserts.NoError(runtime.Run("git", "config", "status.short", "true"))
					have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
					must.True(t, have.OpenChanges)
				})
			})

			t.Run("status.branch enabled", func(t *testing.T) {
				t.Parallel()
				t.Run("no open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					asserts.NoError(runtime.Run("git", "config", "status.branch", "true"))
					have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					asserts.NoError(runtime.Run("git", "config", "status.branch", "true"))
					have := asserts.NoError1(runtime.Git.RepoStatus(runtime))
					must.True(t, have.OpenChanges)
				})
			})
		})
	})

	t.Run("ResetCurrentBranchToSHA", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch1 := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch1, initial.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch1,
			FileContent: "file1",
			FileName:    "file1",
			Message:     "commit 1",
		})
		branch2 := gitdomain.NewLocalBranchName("branch2")
		runtime.CreateBranch(branch2, branch1.BranchName())
		runtime.CreateCommit(testgit.Commit{
			Branch:      branch2,
			FileContent: "file2",
			FileName:    "file2",
			Message:     "commit 2",
		})
		branch1SHA := runtime.SHAforBranch(branch1)
		branch2SHA := runtime.SHAforBranch(branch2)
		must.NotEqOp(t, branch1SHA, branch2SHA)
		runtime.CheckoutBranch(branch1)
		asserts.NoError(runtime.Git.ResetCurrentBranchToSHA(runtime, branch2SHA))
		newBranch1SHA := runtime.SHAforBranch(branch1)
		must.EqOp(t, newBranch1SHA, branch2SHA)
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
		commits := asserts.NoError1(runtime.Git.CommitsInBranch(runtime.TestRunner, "branch", gitdomain.NewLocalBranchNameOption("initial")))
		have := asserts.NoError1(runtime.Git.ShortenSHA(runtime, commits[0].SHA))
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
			have := asserts.NoError1(runtime.Git.StashSize(runtime.TestRunner))
			want := gitdomain.StashSize(2)
			must.EqOp(t, want, have)
		})
		t.Run("no stash entries", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have := asserts.NoError1(runtime.Git.StashSize(runtime.TestRunner))
			want := gitdomain.StashSize(0)
			must.EqOp(t, want, have)
		})
	})

	t.Run("BranchesAvailableInCurrentWorktree", func(t *testing.T) {
		t.Parallel()
		t.Run("includes branches not checked out anywhere", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("branch1", initial.BranchName())
			runtime.CreateBranch("branch2", initial.BranchName())
			runtime.CheckoutBranch("branch1")
			available := asserts.NoError1(runtime.Git.BranchesAvailableInCurrentWorktree(runtime))
			must.SliceContainsAll(t, gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("initial"),
				gitdomain.NewLocalBranchName("branch1"),
				gitdomain.NewLocalBranchName("branch2"),
			}, available)
		})

		t.Run("excludes branch checked out in another worktree", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("branch1", initial.BranchName())
			runtime.CreateBranch("branch2", initial.BranchName())
			worktreeDir := t.TempDir()
			runtime.AddWorktree(worktreeDir, "branch1")
			available := asserts.NoError1(runtime.Git.BranchesAvailableInCurrentWorktree(runtime))
			must.SliceContainsAll(t, gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("initial"),
				gitdomain.NewLocalBranchName("branch2"),
			}, available)
		})

		t.Run("excludes remote branches", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			origin.CreateBranch("remote-branch", initial.BranchName())
			local.Fetch()
			available := asserts.NoError1(local.Git.BranchesAvailableInCurrentWorktree(local))
			must.SliceContainsAll(t, gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("initial"),
			}, available)
		})

		t.Run("empty repository", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			available := asserts.NoError1(runtime.Git.BranchesAvailableInCurrentWorktree(runtime))
			must.SliceContainsAll(t, gitdomain.LocalBranchNames{
				gitdomain.NewLocalBranchName("initial"),
			}, available)
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
