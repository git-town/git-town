package git_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/git"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/gohacks"
	"github.com/git-town/git-town/v18/internal/gohacks/cache"
	"github.com/git-town/git-town/v18/internal/subshell"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/git-town/git-town/v18/test/testgit"
	"github.com/git-town/git-town/v18/test/testruntime"
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
		authors, err := runtime.BranchAuthors(runtime.TestRunner, branch, initial)
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
			err := runtime.CreateAndCheckoutBranch(runtime, branch1)
			must.NoError(t, err)
			runtime.CreateBranch(branch2, branch1.BranchName())
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			runtime.CheckoutBranch(branch2)
			err = runtime.MergeNoFastForward(runtime, configdomain.UseDefaultMessage(), branch1)
			must.NoError(t, err)
			have, err := runtime.BranchContainsMerges(runtime, branch2, branch1)
			must.NoError(t, err)
			must.True(t, have)
		})
		t.Run("branch has no merge commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			err := runtime.CreateAndCheckoutBranch(runtime, branch1)
			must.NoError(t, err)
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch1,
				FileContent: "content",
				FileName:    "file1",
				Message:     "commit 1",
			})
			have, err := runtime.BranchContainsMerges(runtime, branch1, initial)
			must.NoError(t, err)
			must.False(t, have)
		})
	})

	t.Run("BranchHasUnmergedChanges", func(t *testing.T) {
		t.Parallel()
		t.Run("branch without commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch")
			runtime.CreateBranch(branch, initial.BranchName())
			have, err := runtime.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
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
			have, err := runtime.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
			must.NoError(t, err)
			must.True(t, have, must.Sprint("branch with commits that make changes"))
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "original content",
				FileName:    "file1",
				Message:     "commit 3",
			})
			have, err = runtime.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
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
			err := local.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "parent_file",
				Message:     "add parent file",
			})
			err = local.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			inSync, err := local.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("parent has extra commit", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "file",
				Message:     "commit on both parent and child",
			})
			err = local.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on parent",
			})
			inSync, err := local.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.False(t, inSync)
		})
		t.Run("child has extra commit", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "parent",
				FileContent: "content",
				FileName:    "file",
				Message:     "commit on both parent and child",
			})
			err = local.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "child",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on child",
			})
			inSync, err := local.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("empty parent", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			err = local.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "child",
				FileContent: "content 2",
				FileName:    "file",
				Message:     "commit only on child",
			})
			inSync, err := local.BranchInSyncWithParent(local.TestRunner, "child", "parent")
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("both empty", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "parent")
			must.NoError(t, err)
			err = local.CreateAndCheckoutBranch(local.TestRunner, "child")
			must.NoError(t, err)
			inSync, err := local.BranchInSyncWithParent(local.TestRunner, "child", "parent")
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
			err := local.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			inSync, err := local.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.True(t, inSync)
		})
		t.Run("branch has local commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			local.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "local_file",
				Message:     "add local file",
			})
			inSync, err := local.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.False(t, inSync)
		})
		t.Run("branch has remote commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			origin.CreateCommit(testgit.Commit{
				Branch:      "branch",
				FileContent: "content",
				FileName:    "remote_file",
				Message:     "add remote file",
			})
			local.Fetch()
			inSync, err := local.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.False(t, inSync)
		})
		t.Run("branch has different local and remote commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
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
			inSync, err := local.BranchInSyncWithTracking(local.TestRunner, "branch", gitdomain.RemoteOrigin)
			must.NoError(t, err)
			must.False(t, inSync)
		})
	})

	t.Run("CheckoutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch, initial.BranchName())
		runtime.CheckoutBranch(branch)
		currentBranch, err := runtime.CurrentBranch(runtime.TestRunner)
		must.NoError(t, err)
		must.EqOp(t, branch, currentBranch)
		runtime.CheckoutBranch(initial)
		currentBranch, err = runtime.CurrentBranch(runtime.TestRunner)
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
			commits, err := runtime.Commands.CommitsInBranch(runtime.TestRunner, branch, Some(gitdomain.NewLocalBranchName("initial")))
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
			commits, err := runtime.Commands.CommitsInBranch(runtime, branch, Some(gitdomain.NewLocalBranchName("initial")))
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
			commits, err := runtime.Commands.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			must.EqOp(t, 3, len(commits)) // 1 initial commit + 2 test commits
		})
		t.Run("main branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			commits, err := runtime.Commands.CommitsInBranch(runtime, initial, None[gitdomain.LocalBranchName]())
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
			commits, err := runtime.Commands.CommitsInFeatureBranch(runtime.TestRunner, branch, gitdomain.NewLocalBranchName("initial"))
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
			commits, err := runtime.Commands.CommitsInFeatureBranch(runtime, branch, gitdomain.NewLocalBranchName("initial"))
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
		branch, err := runtime.Commands.CurrentBranch(runtime)
		must.NoError(t, err)
		must.EqOp(t, branch, branch)
		runtime.CheckoutBranch(initial)
		branch, err = runtime.Commands.CurrentBranch(runtime)
		must.NoError(t, err)
		must.EqOp(t, initial, branch)
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
			have := repo.CurrentBranchHasTrackingBranch(repo)
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
			have := repo.CurrentBranchHasTrackingBranch(repo)
			must.False(t, have)
		})
	})

	t.Run("DefaultBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.SetDefaultGitBranch("main")
		have := runtime.Commands.DefaultBranch(runtime)
		want := Some(gitdomain.NewLocalBranchName("main"))
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
			have := git.DetectPhantomMergeConflicts(fullInfos, Some(gitdomain.NewLocalBranchName("alpha")), "main")
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
			have := git.DetectPhantomMergeConflicts(fullInfos, Some(gitdomain.NewLocalBranchName("alpha")), "main")
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
			have := git.DetectPhantomMergeConflicts(fullInfos, Some(gitdomain.NewLocalBranchName("alpha")), "main")
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
			have := git.DetectPhantomMergeConflicts(fullInfos, Some(gitdomain.NewLocalBranchName("alpha")), "main")
			want := []git.PhantomMergeConflict{}
			must.Eq(t, want, have)
		})
	})

	t.Run("FirstCommitMessageInBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is empty", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			must.NoError(t, repo.CreateAndCheckoutBranch(repo.TestRunner, "branch"))
			have, err := repo.FirstCommitMessageInBranch(repo.TestRunner, "branch", "main")
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
			have, err := repo.FirstCommitMessageInBranch(repo.TestRunner, branch.BranchName(), main.BranchName())
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
			have, err := repo.FirstCommitMessageInBranch(repo.TestRunner, branch.BranchName(), main.BranchName())
			must.NoError(t, err)
			want := Some(gitdomain.CommitMessage("commit message 1"))
			must.Eq(t, want, have)
		})
		t.Run("branch doesn't exist", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			_, err := repo.FirstCommitMessageInBranch(repo.TestRunner, "zonk", "main")
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
			err := repo.DeleteLocalBranch(repo.TestRunner, branch)
			must.NoError(t, err)
			have, err := repo.FirstCommitMessageInBranch(repo.TestRunner, branch.TrackingBranch(gitdomain.RemoteOrigin).BranchName(), main.BranchName())
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
			have := runtime.Commands.FirstExistingBranch(runtime, branch1, branch2)
			want := Some(branch1)
			must.Eq(t, want, have)
		})
		t.Run("second branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			runtime.CreateBranch(branch2, initial.BranchName())
			have := runtime.Commands.FirstExistingBranch(runtime, branch1, branch2)
			want := Some(branch2)
			must.Eq(t, want, have)
		})
		t.Run("no branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			have := runtime.Commands.FirstExistingBranch(runtime, branch1, branch2)
			want := None[gitdomain.LocalBranchName]()
			must.EqOp(t, want, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch("b1", initial.BranchName())
		runner.CreateBranch("b2", initial.BranchName())
		must.True(t, runner.Commands.HasLocalBranch(runner, "b1"))
		must.True(t, runner.Commands.HasLocalBranch(runner, "b2"))
		must.False(t, runner.Commands.HasLocalBranch(runner, "b3"))
	})

	t.Run("IsAhead", func(t *testing.T) {
		t.Parallel()
		t.Run("is actually ahead", func(t *testing.T) {
			t.Parallel()
			isAhead, remoteBranchName := git.IsAhead("branch-1", "[origin/branch-1: ahead 10] commit message")
			must.True(t, isAhead)
			must.EqOp(t, "origin/branch-1", remoteBranchName.String())
		})
		t.Run("is not ahead", func(t *testing.T) {
			t.Parallel()
			isAhead, remoteBranchName := git.IsAhead("branch-1", "[origin/branch-1: behind 10] commit message")
			must.False(t, isAhead)
			must.EqOp(t, "", remoteBranchName.String())
		})
	})

	t.Run("IsBehind", func(t *testing.T) {
		t.Parallel()
		has, branchNameOpt := git.IsBehind("production", "[origin/production: behind 1] initial commit")
		must.True(t, has)
		branchName, hasBranchName := branchNameOpt.Get()
		must.True(t, hasBranchName)
		must.EqOp(t, "origin/production", branchName)
	})

	t.Run("IsInSync", func(t *testing.T) {
		t.Parallel()
		t.Run("is actually in sync", func(t *testing.T) {
			t.Parallel()
			isInSync, remoteBranchName := git.IsInSync("branch-1", "[origin/branch-1] commit message")
			must.True(t, isInSync)
			must.EqOp(t, "origin/branch-1", remoteBranchName.String())
		})
	})

	t.Run("IsRemoteGone", func(t *testing.T) {
		t.Parallel()
		t.Run("remote is gone", func(t *testing.T) {
			t.Parallel()
			isGone, remoteBranchName := git.IsRemoteGone("branch-1", "[origin/branch-1: gone] commit message")
			must.True(t, isGone)
			must.Eq(t, Some(gitdomain.NewRemoteBranchName("origin/branch-1")), remoteBranchName)
		})
		t.Run("other sync status", func(t *testing.T) {
			t.Parallel()
			isGone, remoteBranchName := git.IsRemoteGone("branch-1", "[origin/branch-1: ahead] commit message")
			must.False(t, isGone)
			must.Eq(t, None[gitdomain.RemoteBranchName](), remoteBranchName)
		})
		t.Run("other text", func(t *testing.T) {
			t.Parallel()
			isGone, remoteBranchName := git.IsRemoteGone("branch-1", "[skip ci]")
			must.False(t, isGone)
			must.Eq(t, None[gitdomain.RemoteBranchName](), remoteBranchName)
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
		err := runtime.MergeFastForward(runtime.TestRunner, branch.BranchName())
		must.NoError(t, err)
		commits, err := runtime.Commands.CommitsInPerennialBranch(runtime) // Current branch.
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
		err := runtime.MergeNoFastForward(runtime.TestRunner, configdomain.UseDefaultMessage(), branch)
		must.NoError(t, err)
		commits, err := runtime.Commands.CommitsInPerennialBranch(runtime) // Current branch.
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
		err := runtime.MergeNoFastForward(runtime.TestRunner, configdomain.UseCustomMessage(mergeMessage), branch)
		must.NoError(t, err)
		commits, err := runtime.Commands.CommitsInPerennialBranch(runtime) // Current branch.
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

	t.Run("ParseVerboseBranchesOutput", func(t *testing.T) {
		t.Parallel()
		t.Run("recognizes the branch names", func(t *testing.T) {
			t.Parallel()
			t.Run("marker is at the first entry", func(t *testing.T) {
				t.Parallel()
				give := `
* branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
  branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, Some(gitdomain.NewLocalBranchName("branch-1")), currentBranch)
			})
			t.Run("marker is at the middle entry", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
* branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, Some(gitdomain.NewLocalBranchName("branch-2")), currentBranch)
			})
			t.Run("marker is at the last entry", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
  branch-2                     da796a69 [origin/branch-2] Commit message 2
* branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, Some(gitdomain.NewLocalBranchName("branch-3")), currentBranch)
			})
			t.Run("in the middle of a rebase", func(t *testing.T) {
				t.Parallel()
				give := `
				* (no branch, rebasing main) 214ba79 origin main commit
  feature                    62bf22e [origin/feature: ahead 1] feature commit
  main                       11716d4 [origin/main: ahead 1, behind 1] local main commit`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, None[gitdomain.LocalBranchName](), currentBranch)
			})
		})

		t.Run("recognize the branch sync status", func(t *testing.T) {
			t.Parallel()
			t.Run("branch is ahead of its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1: ahead 1] Commit message 1a
  remotes/origin/branch-1      222222 Commit message 1b`[1:]
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusAhead,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("branch is behind its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     1111111111111111111111111111111111111111 [origin/branch-1: behind 2] Commit message 1
  remotes/origin/branch-1      2222222222222222222222222222222222222222 Commit message 1b`[1:]
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("1111111111111111111111111111111111111111")),
						SyncStatus: gitdomain.SyncStatusBehind,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("2222222222222222222222222222222222222222")),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("branch is ahead and behind its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1: ahead 31, behind 2] Commit message 1a
  remotes/origin/branch-1      222222 Commit message 1b`[1:]
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusNotInSync,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("branch is in sync", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1] Commit message 1
  remotes/origin/branch-1      111111 Commit message 1`[1:]
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("111111")),
						SyncStatus: gitdomain.SyncStatusUpToDate,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("remote-only branch", func(t *testing.T) {
				t.Parallel()
				give := `
  remotes/origin/branch-1    222222 Commit message 2`[1:]
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  None[gitdomain.LocalBranchName](),
						LocalSHA:   None[gitdomain.SHA](),
						SyncStatus: gitdomain.SyncStatusRemoteOnly,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("local-only branch", func(t *testing.T) {
				t.Parallel()
				give := `  branch-1                     01a7eded Commit message 1`
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("01a7eded")),
						SyncStatus: gitdomain.SyncStatusLocalOnly,
						RemoteName: None[gitdomain.RemoteBranchName](),
						RemoteSHA:  None[gitdomain.SHA](),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("branch is deleted at the remote", func(t *testing.T) {
				t.Parallel()
				isGone, remoteBranchName := git.IsRemoteGone("branch-1", "[origin/branch-1: gone] commit message")
				must.True(t, isGone)
				must.Eq(t, Some(gitdomain.NewRemoteBranchName("origin/branch-1")), remoteBranchName)
			})

			t.Run("branch is active in another worktree", func(t *testing.T) {
				t.Parallel()
				give := `+ branch-1    3d0c4c13 (/path/to/other/worktree) [origin/branch-1] commit message`
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("3d0c4c13")),
						SyncStatus: gitdomain.SyncStatusOtherWorktree,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
						RemoteSHA:  None[gitdomain.SHA](),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
			})

			t.Run("in the middle of a rebase", func(t *testing.T) {
				t.Parallel()
				give := `
* (no branch, rebasing main) 214ba79 origin main commit
  feature                    62bf22e [origin/feature: ahead 1] feature commit
  main                       11716d4 [origin/main: ahead 1, behind 1] local main commit
  remotes/origin/feature     4989007 initial commit
  remotes/origin/main        214ba79 origin main commit`[1:]

				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("feature")),
						LocalSHA:   Some(gitdomain.NewSHA("62bf22e")),
						SyncStatus: gitdomain.SyncStatusAhead,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature")),
						RemoteSHA:  Some(gitdomain.NewSHA("4989007")),
					},
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("main")),
						LocalSHA:   Some(gitdomain.NewSHA("11716d4")),
						SyncStatus: gitdomain.SyncStatusNotInSync,
						RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
						RemoteSHA:  Some(gitdomain.SHA("214ba79")),
					},
				}
				have, active := git.ParseVerboseBranchesOutput(give)
				must.Eq(t, want, have)
				must.Eq(t, None[gitdomain.LocalBranchName](), active)
			})
		})

		t.Run("square brackets in the commit message", func(t *testing.T) {
			t.Parallel()
			give := `
  branch-1                 111111 [ci skip]
  branch-2                 222222 ️[origin/branch-2] [ci skip]
  remotes/origin/branch-2  222222 [ci skip]
  remotes/origin/branch-3  333333 [ci skip]`[1:]
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-2")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-3")),
					RemoteSHA:  Some(gitdomain.NewSHA("333333")),
				},
			}
			have, _ := git.ParseVerboseBranchesOutput(give)
			must.Eq(t, want, have)
		})

		t.Run("symbolic reference", func(t *testing.T) {
			t.Parallel()
			give := `
  main                                   4dc97db26 [origin/main] Commit 1
  remotes/origin/HEAD                    -> origin/main
  remotes/origin/main                    4dc97db26 Commit 1
  remotes/origin/master                  -> origin/main
`[1:]
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("4dc97db26")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("4dc97db26")),
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			must.Eq(t, want, have)
			must.Eq(t, None[gitdomain.LocalBranchName](), currentBranch)
		})

		t.Run("complex example", func(t *testing.T) {
			give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1a
* branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a
  main                         41c3f128 [origin/main: behind 2] Commit message on main (#1234)
  branch-4                     e4d6bc09 [origin/branch-4: gone] Commit message 4
+ branch-5                     55555555 (/path/to/other/worktree) [origin/branch-5] Commit message 5
  remotes/origin/branch-1      307a7bf4 Commit message 1b
  remotes/origin/branch-2      da796a69 Commit message 2
  remotes/origin/branch-3      bc39378a Commit message 3b
  remotes/origin/branch-5      55555555 Commit message 5
  remotes/origin/HEAD          -> origin/initial
  remotes/origin/main          02c192178 Commit message on main (#1234)
  remotes/upstream/HEAD        -> upstream/main

`[1:]
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("01a7eded")),
					SyncStatus: gitdomain.SyncStatusAhead,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("307a7bf4")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-2")),
					LocalSHA:   Some(gitdomain.NewSHA("da796a69")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
					RemoteSHA:  Some(gitdomain.NewSHA("da796a69")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-3")),
					LocalSHA:   Some(gitdomain.NewSHA("f4ebec0a")),
					SyncStatus: gitdomain.SyncStatusBehind,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-3")),
					RemoteSHA:  Some(gitdomain.NewSHA("bc39378a")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("41c3f128")),
					SyncStatus: gitdomain.SyncStatusBehind,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("02c192178")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-4")),
					LocalSHA:   Some(gitdomain.NewSHA("e4d6bc09")),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-4")),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-5")),
					LocalSHA:   Some(gitdomain.NewSHA("55555555")),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-5")),
					RemoteSHA:  Some(gitdomain.NewSHA("55555555")),
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			must.Eq(t, want, have)
			must.Eq(t, Some(gitdomain.NewLocalBranchName("branch-2")), currentBranch)
		})
	})

	t.Run("PreviouslyCheckedOutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("feature1", initial.BranchName())
		runtime.CreateBranch("feature2", initial.BranchName())
		runtime.CheckoutBranch("feature1")
		runtime.CheckoutBranch("feature2")
		have := runtime.Commands.PreviouslyCheckedOutBranch(runtime.TestRunner)
		must.Eq(t, Some(gitdomain.NewLocalBranchName("feature1")), have)
	})

	t.Run("Remotes", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		runtime.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		remotes, err := runtime.Commands.Remotes(runtime.TestRunner)
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{gitdomain.RemoteOrigin}, remotes)
	})

	t.Run("RepoStatus", func(t *testing.T) {
		t.Run("OpenChanges", func(t *testing.T) {
			t.Parallel()
			t.Run("no open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				have, err := runtime.Commands.RepoStatus(runtime)
				must.NoError(t, err)
				must.False(t, have.OpenChanges)
			})
			t.Run("has open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateFile("foo", "bar")
				have, err := runtime.Commands.RepoStatus(runtime)
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
				have, err := runtime.Commands.RepoStatus(runtime)
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
				have, err := runtime.Commands.RepoStatus(runtime)
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
				have, err := runtime.Commands.RepoStatus(runtime)
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
					have, err := runtime.Commands.RepoStatus(runtime)
					must.NoError(t, err)
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					err := runtime.Run("git", "config", "status.short", "true")
					must.NoError(t, err)
					have, err := runtime.Commands.RepoStatus(runtime)
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
					have, err := runtime.Commands.RepoStatus(runtime)
					must.NoError(t, err)
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					err := runtime.Run("git", "config", "status.branch", "true")
					must.NoError(t, err)
					have, err := runtime.Commands.RepoStatus(runtime)
					must.NoError(t, err)
					must.True(t, have.OpenChanges)
				})
			})
		})

		t.Run("RebaseInProgress", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have, err := runtime.Commands.RepoStatus(runtime)
			must.NoError(t, err)
			must.False(t, have.RebaseInProgress)
		})
	})

	t.Run("RootDirectory", func(t *testing.T) {
		t.Parallel()
		t.Run("inside a Git repo", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have := runtime.Commands.RootDirectory(runtime.TestRunner)
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

	t.Run("StashEntries", func(t *testing.T) {
		t.Parallel()
		t.Run("some stash entries", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("file1", "content")
			runtime.StashOpenFiles()
			runtime.CreateFile("file2", "content")
			runtime.StashOpenFiles()
			have, err := runtime.StashSize(runtime.TestRunner)
			want := gitdomain.StashSize(2)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		})
		t.Run("no stash entries", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have, err := runtime.StashSize(runtime.TestRunner)
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

	t.Run("parseActiveBranchDuringRebase", func(t *testing.T) {
		t.Parallel()
		t.Run("branch name is one word", func(t *testing.T) {
			t.Parallel()
			give := "* (no branch, rebasing feature)"
			have := git.ParseActiveBranchDuringRebase(give)
			want := gitdomain.NewLocalBranchName("feature")
			must.Eq(t, want, have)
		})
		t.Run("branch name is two words", func(t *testing.T) {
			t.Parallel()
			give := "* (no branch, rebasing feature branch)"
			have := git.ParseActiveBranchDuringRebase(give)
			want := gitdomain.NewLocalBranchName("feature branch")
			must.Eq(t, want, have)
		})
		t.Run("branch name is three words", func(t *testing.T) {
			t.Parallel()
			give := "* (no branch, rebasing the feature branch)"
			have := git.ParseActiveBranchDuringRebase(give)
			want := gitdomain.NewLocalBranchName("the feature branch")
			must.Eq(t, want, have)
		})
	})
}
