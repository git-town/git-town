package git_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks"
	"github.com/git-town/git-town/v15/internal/gohacks/cache"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/subshell"
	testgit "github.com/git-town/git-town/v15/test/git"
	"github.com/git-town/git-town/v15/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestBackendCommands(t *testing.T) {
	t.Parallel()
	initial := gitdomain.NewLocalBranchName("initial")

	t.Run("BranchAuthors", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch")
		runtime.CreateBranch(branch, initial)
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
		authors, err := runtime.TestCommands.BranchAuthors(runtime.TestRunner, branch, initial)
		must.NoError(t, err)
		must.Eq(t, []gitdomain.Author{"user <email@example.com>"}, authors)
	})

	t.Run("BranchHasUnmergedChanges", func(t *testing.T) {
		t.Parallel()
		t.Run("branch without commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch := gitdomain.NewLocalBranchName("branch")
			runtime.CreateBranch(branch, initial)
			have, err := runtime.TestCommands.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
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
			runtime.CreateBranch(branch, initial)
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "modified content",
				FileName:    "file1",
				Message:     "commit 2",
			})
			have, err := runtime.TestCommands.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
			must.NoError(t, err)
			must.True(t, have, must.Sprint("branch with commits that make changes"))
			runtime.CreateCommit(testgit.Commit{
				Branch:      branch,
				FileContent: "original content",
				FileName:    "file1",
				Message:     "commit 3",
			})
			have, err = runtime.TestCommands.BranchHasUnmergedChanges(runtime.TestRunner, branch, initial)
			must.NoError(t, err)
			must.False(t, have, must.Sprint("branch with commits that make no changes"))
		})
	})

	t.Run("CheckoutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		branch := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch, initial)
		runtime.TestCommands.CheckoutBranch(branch)
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
			runtime.CreateBranch(branch, initial)
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
			runtime.CreateBranch(branch, initial)
			commits, err := runtime.Commands.CommitsInBranch(runtime.TestCommands, branch, Some(gitdomain.NewLocalBranchName("initial")))
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
			commits, err := runtime.Commands.CommitsInBranch(runtime.TestCommands, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			must.EqOp(t, 3, len(commits)) // 1 initial commit + 2 test commits
		})
		t.Run("main branch contains no commits", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			commits, err := runtime.Commands.CommitsInBranch(runtime.TestCommands, initial, None[gitdomain.LocalBranchName]())
			must.NoError(t, err)
			must.EqOp(t, 1, len(commits)) // the initial commit
		})
	})

	t.Run("CurrentBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CheckoutBranch(initial)
		branch := gitdomain.NewLocalBranchName("branch1")
		runtime.CreateBranch(branch, initial)
		runtime.CheckoutBranch(branch)
		branch, err := runtime.Commands.CurrentBranch(runtime.TestCommands)
		must.NoError(t, err)
		must.EqOp(t, branch, branch)
		runtime.CheckoutBranch(initial)
		branch, err = runtime.Commands.CurrentBranch(runtime.TestCommands)
		must.NoError(t, err)
		must.EqOp(t, initial, branch)
	})

	t.Run("DefaultBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.SetDefaultGitBranch("main")
		have := runtime.Commands.DefaultBranch(runtime.TestCommands)
		want := Some(gitdomain.NewLocalBranchName("main"))
		must.Eq(t, want, have)
	})

	t.Run("FirstExistingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("first branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			runtime.CreateBranch(branch1, initial)
			runtime.CreateBranch(branch2, initial)
			have := runtime.Commands.FirstExistingBranch(runtime.TestCommands, branch1, branch2)
			want := Some(branch1)
			must.Eq(t, want, have)
		})
		t.Run("second branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			runtime.CreateBranch(branch2, initial)
			have := runtime.Commands.FirstExistingBranch(runtime.TestCommands, branch1, branch2)
			want := Some(branch2)
			must.Eq(t, want, have)
		})
		t.Run("no branch matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			branch1 := gitdomain.NewLocalBranchName("b1")
			branch2 := gitdomain.NewLocalBranchName("b2")
			have := runtime.Commands.FirstExistingBranch(runtime.TestCommands, branch1, branch2)
			want := None[gitdomain.LocalBranchName]()
			must.EqOp(t, want, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch(gitdomain.NewLocalBranchName("b1"), initial)
		runner.CreateBranch(gitdomain.NewLocalBranchName("b2"), initial)
		must.True(t, runner.Commands.HasLocalBranch(runner.TestCommands, gitdomain.NewLocalBranchName("b1")))
		must.True(t, runner.Commands.HasLocalBranch(runner.TestCommands, gitdomain.NewLocalBranchName("b2")))
		must.False(t, runner.Commands.HasLocalBranch(runner.TestCommands, gitdomain.NewLocalBranchName("b3")))
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

	t.Run("RepoStatus", func(t *testing.T) {
		t.Run("HasOpenChanges", func(t *testing.T) {
			t.Parallel()
			t.Run("no open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
				must.NoError(t, err)
				must.False(t, have.OpenChanges)
			})
			t.Run("has open changes", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				runtime.CreateFile("foo", "bar")
				have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
				must.NoError(t, err)
				must.True(t, have.OpenChanges)
			})
			t.Run("during rebase", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				branch1 := gitdomain.NewLocalBranchName("branch1")
				runtime.CreateBranch(branch1, initial)
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
				have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
				must.NoError(t, err)
				must.False(t, have.OpenChanges)
			})
			t.Run("during merge conflict", func(t *testing.T) {
				t.Parallel()
				runtime := testruntime.Create(t)
				branch1 := gitdomain.NewLocalBranchName("branch1")
				runtime.CreateBranch(branch1, initial)
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
				have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
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
				have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
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
					have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
					must.NoError(t, err)
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					err := runtime.Run("git", "config", "status.short", "true")
					must.NoError(t, err)
					have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
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
					have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
					must.NoError(t, err)
					must.False(t, have.OpenChanges)
				})
				t.Run("open changes", func(t *testing.T) {
					t.Parallel()
					runtime := testruntime.Create(t)
					runtime.CreateFile("file", "stashed content")
					err := runtime.Run("git", "config", "status.branch", "true")
					must.NoError(t, err)
					have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
					must.NoError(t, err)
					must.True(t, have.OpenChanges)
				})
			})
		})

		t.Run("RebaseInProgress", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have, err := runtime.Commands.RepoStatus(runtime.TestCommands)
			must.NoError(t, err)
			must.False(t, have.RebaseInProgress)
		})
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

		t.Run("recognizes the branch sync status", func(t *testing.T) {
			t.Parallel()
			t.Run("branch is ahead of its remote branch", func(t *testing.T) {
				t.Parallel()
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
				t.Run("determineSyncStatus", func(t *testing.T) {
					t.Parallel()
					give := `
  branch-1                     111111 [origin/branch-1: ahead 1] Commit message 1a
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
			})

			t.Run("branch is behind its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1: behind 2] Commit message 1
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

			t.Run("branch is in sync with its remote branch", func(t *testing.T) {
				t.Parallel()
				t.Run("IsInSync", func(t *testing.T) {
					t.Parallel()
					t.Run("is actually in sync", func(t *testing.T) {
						t.Parallel()
						isInSync, remoteBranchName := git.IsInSync("branch-1", "[origin/branch-1] commit message")
						must.True(t, isInSync)
						must.EqOp(t, "origin/branch-1", remoteBranchName.String())
					})
				})
				t.Run("ParseVerboseBranchesOutput", func(t *testing.T) {
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
			})

			t.Run("remote gone", func(t *testing.T) {
				t.Parallel()
				give := `  branch-1                     01a7eded [origin/branch-1: gone] Commit message 1`
				want := gitdomain.BranchInfos{
					gitdomain.BranchInfo{
						LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
						LocalSHA:   Some(gitdomain.NewSHA("01a7eded")),
						SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
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
						SyncStatus: gitdomain.SyncStatusNotInSync,
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
					SyncStatus: gitdomain.SyncStatusNotInSync,
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
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-3")),
					RemoteSHA:  Some(gitdomain.NewSHA("bc39378a")),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("main")),
					LocalSHA:   Some(gitdomain.NewSHA("41c3f128")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
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
		runtime.CreateBranch(gitdomain.NewLocalBranchName("feature1"), initial)
		runtime.CreateBranch(gitdomain.NewLocalBranchName("feature2"), initial)
		runtime.CheckoutBranch(gitdomain.NewLocalBranchName("feature1"))
		runtime.CheckoutBranch(gitdomain.NewLocalBranchName("feature2"))
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
				CurrentBranchCache: &cache.LocalBranchWithPrevious{},
				RemotesCache:       &cache.Remotes{},
			}
			have := cmds.RootDirectory(runner)
			must.True(t, have.IsNone())
		})
	})

	t.Run("ShouldPushBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("branch has no commits", func(t *testing.T) {
			t.Parallel()
			origin := testruntime.Create(t)
			local := testruntime.Clone(origin.TestRunner, t.TempDir())
			err := local.CreateAndCheckoutBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			err = local.CreateTrackingBranch(local.TestRunner, "branch", gitdomain.RemoteOrigin, false)
			must.NoError(t, err)
			shouldPush, err := local.ShouldPushBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			must.False(t, shouldPush)
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
			shouldPush, err := local.ShouldPushBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			must.True(t, shouldPush)
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
			shouldPush, err := local.ShouldPushBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			must.True(t, shouldPush)
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
			shouldPush, err := local.ShouldPushBranch(local.TestRunner, "branch")
			must.NoError(t, err)
			must.True(t, shouldPush)
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
}
