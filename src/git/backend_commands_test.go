package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/statistics"
	"github.com/git-town/git-town/v9/src/subshell"
	testgit "github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestBackendCommands(t *testing.T) {
	t.Parallel()

	t.Run("BranchAuthors", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch", "initial")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch",
			FileName:    "file1",
			FileContent: "file1",
			Message:     "first commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch",
			FileName:    "file2",
			FileContent: "file2",
			Message:     "second commit",
		})
		authors, err := runtime.Backend.BranchAuthors("branch", "initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"user <email@example.com>"}, authors)
	})

	t.Run(".CheckoutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch1", "initial")
		assert.NoError(t, runtime.Backend.CheckoutBranch("branch1"))
		currentBranch, err := runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "branch1", currentBranch)
		runtime.CheckoutBranch("initial")
		currentBranch, err = runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", currentBranch)
	})

	t.Run(".CreateFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		err := runtime.Backend.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		runtime.Config.Reload()
		assert.True(t, runtime.Config.BranchDurations().IsFeatureBranch("f1"))
		lineageHave := runtime.Config.Lineage()
		lineageWant := config.Lineage{}
		lineageWant["f1"] = "main"
		assert.Equal(t, lineageWant, lineageHave)
	})

	t.Run(".CurrentBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CheckoutBranch("initial")
		runtime.CreateBranch("b1", "initial")
		runtime.CheckoutBranch("b1")
		branch, err := runtime.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "b1", branch)
		runtime.CheckoutBranch("initial")
		branch, err = runtime.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch)
	})

	t.Run(".HasLocalBranch()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch("b1", "initial")
		runner.CreateBranch("b2", "initial")
		has, err := runner.Backend.HasLocalBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = runner.Backend.HasLocalBranch("b2")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = runner.Backend.HasLocalBranch("b3")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasOpenChanges()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		has, err := runtime.Backend.HasOpenChanges()
		assert.NoError(t, err)
		assert.False(t, has)
		runtime.CreateFile("foo", "bar")
		has, err = runtime.Backend.HasOpenChanges()
		assert.NoError(t, err)
		assert.True(t, has)
	})

	t.Run(".HasRebaseInProgress()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		has, err := runtime.Backend.HasRebaseInProgress()
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".LocalBranchesMainFirst()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch("b1", "initial")
		runner.CreateBranch("b2", "initial")
		origin.CreateBranch("b3", "initial")
		runner.Fetch()
		branches, err := runner.Backend.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2"}, branches)
	})

	t.Run("ParseVerboseBranchesOutput", func(t *testing.T) {
		t.Parallel()
		t.Run("recognizes the current branch", func(t *testing.T) {
			t.Run("marker is at the first entry", func(t *testing.T) {
				give := `
* branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
  branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, "branch-1", currentBranch)
			})
			t.Run("marker is at the middle entry", func(t *testing.T) {
				give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
* branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, "branch-2", currentBranch)
			})
			t.Run("marker is at the last entry", func(t *testing.T) {
				give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
  branch-2                     da796a69 [origin/branch-2] Commit message 2
* branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, "branch-3", currentBranch)
			})
		})

		t.Run("uses the remote branch name given by Git", func(t *testing.T) {
			t.Run("the remote branch name is the same as the local branch name", func(t *testing.T) {
				give := `
  branch-1                     11111111 [origin/branch-1: ahead 1] Commit message 1
	remotes/origin/branch-1      22222222 Commit message 2
	remotes/origin/branch-A      33333333 Commit message 3`[1:]
				have, _ := git.ParseVerboseBranchesOutput(give)
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						LocalSHA:   "11111111",
						SyncStatus: git.SyncStatusAhead,
						RemoteName: "origin/branch-1",
						RemoteSHA:  "22222222",
					},
					git.BranchSyncStatus{
						LocalName:  "",
						LocalSHA:   "",
						SyncStatus: git.SyncStatusRemoteOnly,
						RemoteName: "origin/branch-A",
						RemoteSHA:  "33333333",
					},
				}
				assert.Equal(t, want, have)
			})
			t.Run("the remote branch name is different from the local branch name", func(t *testing.T) {
				give := `
  branch-1                     11111111 [origin/branch-A: ahead 1] Commit message 1
	remotes/origin/branch-1      22222222 Commit message 2
	remotes/origin/branch-A      33333333 Commit message 3`[1:]
				have, _ := git.ParseVerboseBranchesOutput(give)
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						LocalSHA:   "11111111",
						SyncStatus: git.SyncStatusAhead,
						RemoteName: "origin/branch-A",
						RemoteSHA:  "33333333",
					},
					git.BranchSyncStatus{
						LocalName:  "",
						LocalSHA:   "",
						SyncStatus: git.SyncStatusRemoteOnly,
						RemoteName: "origin/branch-1",
						RemoteSHA:  "22222222",
					},
				}
				assert.Equal(t, want, have)
			})
		})

		t.Run("recognizes the branch sync status", func(t *testing.T) {
			t.Run("branch is ahead of its remote branch", func(t *testing.T) {
				give := `  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						SyncStatus: git.SyncStatusAhead,
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
			t.Run("recognizes branches that are behind their remote branch", func(t *testing.T) {
				give := `  branch-1                     01a7eded [origin/branch-1: behind 2] Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						SyncStatus: git.SyncStatusBehind,
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
			t.Run("recognizes branches that are in sync with their remote branch", func(t *testing.T) {
				give := `  branch-1                     01a7eded [origin/branch-1] Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						SyncStatus: git.SyncStatusUpToDate,
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
			t.Run("recognizes remote-only branches", func(t *testing.T) {
				give := `  remotes/origin/branch-1                     01a7eded Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						SyncStatus: git.SyncStatusRemoteOnly,
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
			t.Run("recognizes local-only branches", func(t *testing.T) {
				give := `  branch-1                     01a7eded Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						SyncStatus: git.SyncStatusLocalOnly,
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
			t.Run("recognizes branches that got deleted at the remote", func(t *testing.T) {
				give := `  branch-1                     01a7eded [origin/branch-1: gone] Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						LocalName:  "branch-1",
						SyncStatus: git.SyncStatusDeletedAtRemote,
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
		})
		t.Run("complex example", func(t *testing.T) {
			give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1a
* branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a
	branch-4                     2417dac3 [origin/branch-4: ahead 1, behind 31] Commit message 4a
  main                         024df944 [origin/main] Commit message on main (#1234)
  branch-5                     e4d6bc09 [origin/branch-5: gone] Commit message 5
  remotes/origin/branch-1      307a7bf4 Commit message 1b
  remotes/origin/branch-2      da796a69 Commit message 2
  remotes/origin/branch-3      bc39378a Commit message 3b
	remotes/origin/branch-4      ec643bcd Commit message 4b
`[1:]
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					LocalName:  "branch-1",
					RemoteName: "",
					SyncStatus: git.SyncStatusAhead,
					LocalSHA:   "",
					RemoteSHA:  "",
				},
				git.BranchSyncStatus{
					LocalName:  "branch-2",
					RemoteName: "",
					SyncStatus: git.SyncStatusUpToDate,
					LocalSHA:   "",
					RemoteSHA:  "",
				},
				git.BranchSyncStatus{
					LocalName:  "branch-3",
					RemoteName: "",
					SyncStatus: git.SyncStatusBehind,
					LocalSHA:   "",
					RemoteSHA:  "",
				},
				git.BranchSyncStatus{
					LocalName:  "branch-4",
					RemoteName: "",
					SyncStatus: git.SyncStatusAheadAndBehind,
					LocalSHA:   "",
					RemoteSHA:  "",
				},
				git.BranchSyncStatus{
					LocalName:  "main",
					RemoteName: "",
					SyncStatus: git.SyncStatusUpToDate,
					LocalSHA:   "",
					RemoteSHA:  "",
				},
				git.BranchSyncStatus{
					LocalName:  "branch-4",
					RemoteName: "",
					SyncStatus: git.SyncStatusDeletedAtRemote,
					LocalSHA:   "",
					RemoteSHA:  "",
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-2", currentBranch)
		})
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("feature1", "initial")
		runtime.CreateBranch("feature2", "initial")
		runtime.CheckoutBranch("feature1")
		runtime.CheckoutBranch("feature2")
		have := runtime.Backend.PreviouslyCheckedOutBranch()
		assert.Equal(t, "feature1", have)
	})

	t.Run(".Remotes()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		runtime.AddRemote(config.OriginRemote, origin.WorkingDir)
		remotes, err := runtime.Backend.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, config.Remotes{config.OriginRemote}, remotes)
	})

	t.Run(".RootDirectory", func(t *testing.T) {
		t.Parallel()
		t.Run("inside a Git repo", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			have := runtime.BackendCommands.RootDirectory()
			assert.Positive(t, len(have))
		})
		t.Run("outside a Git repo", func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			runner := subshell.BackendRunner{
				Dir:     &dir,
				Verbose: false,
				Stats:   &statistics.None{},
			}
			cmds := git.BackendCommands{
				BackendRunner:      runner,
				Config:             nil,
				CurrentBranchCache: &cache.String{},
				RemoteBranchCache:  &cache.Strings{},
				RemotesCache:       &cache.Strings{},
			}
			have := cmds.RootDirectory()
			assert.Empty(t, have)
		})
	})
}
