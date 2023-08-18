package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
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
		runtime.CreateBranch(domain.NewLocalBranchName("branch"), domain.NewLocalBranchName("initial"))
		runtime.CreateCommit(testgit.Commit{
			Branch:      domain.NewLocalBranchName("branch"),
			FileName:    "file1",
			FileContent: "file1",
			Message:     "first commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      domain.NewLocalBranchName("branch"),
			FileName:    "file2",
			FileContent: "file2",
			Message:     "second commit",
		})
		authors, err := runtime.Backend.BranchAuthors(domain.NewLocalBranchName("branch"), domain.NewLocalBranchName("initial"))
		assert.NoError(t, err)
		assert.Equal(t, []string{"user <email@example.com>"}, authors)
	})

	t.Run(".CheckoutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("initial"))
		assert.NoError(t, runtime.Backend.CheckoutBranch(domain.NewLocalBranchName("branch1")))
		currentBranch, err := runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, domain.NewLocalBranchName("branch1"), currentBranch)
		runtime.CheckoutBranch(domain.NewLocalBranchName("initial"))
		currentBranch, err = runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, domain.NewLocalBranchName("initial"), currentBranch)
	})

	t.Run(".CreateFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		err := runtime.Backend.CreateFeatureBranch(domain.NewLocalBranchName("f1"))
		assert.NoError(t, err)
		runtime.Config.Reload()
		assert.True(t, runtime.Config.BranchDurations().IsFeatureBranch(domain.NewLocalBranchName("f1")))
		lineageHave := runtime.Config.Lineage()
		lineageWant := config.Lineage{}
		lineageWant[domain.NewLocalBranchName("f1")] = domain.NewLocalBranchName("main")
		assert.Equal(t, lineageWant, lineageHave)
	})

	t.Run(".CurrentBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CheckoutBranch(domain.NewLocalBranchName("initial"))
		runtime.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		runtime.CheckoutBranch(domain.NewLocalBranchName("b1"))
		branch, err := runtime.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, domain.NewLocalBranchName("b1"), branch)
		runtime.CheckoutBranch(domain.NewLocalBranchName("initial"))
		branch, err = runtime.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, domain.NewLocalBranchName("initial"), branch)
	})

	t.Run(".HasLocalBranch()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		runner.CreateBranch(domain.NewLocalBranchName("b2"), domain.NewLocalBranchName("initial"))
		has, err := runner.Backend.HasLocalBranch(domain.NewLocalBranchName("b1"))
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = runner.Backend.HasLocalBranch(domain.NewLocalBranchName("b2"))
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = runner.Backend.HasLocalBranch(domain.NewLocalBranchName("b3"))
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
		runner.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		runner.CreateBranch(domain.NewLocalBranchName("b2"), domain.NewLocalBranchName("initial"))
		origin.CreateBranch(domain.NewLocalBranchName("b3"), domain.NewLocalBranchName("initial"))
		runner.Fetch()
		branches, err := runner.Backend.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		assert.NoError(t, err)
		want := domain.LocalBranchNamesFrom("initial", "b1", "b2")
		assert.Equal(t, want, branches)
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
				assert.Equal(t, domain.NewLocalBranchName("branch-1"), currentBranch)
			})
			t.Run("marker is at the middle entry", func(t *testing.T) {
				give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
* branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, domain.NewLocalBranchName("branch-2"), currentBranch)
			})
			t.Run("marker is at the last entry", func(t *testing.T) {
				give := `
  branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
  branch-2                     da796a69 [origin/branch-2] Commit message 2
* branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3a`[1:]
				_, currentBranch := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, domain.NewLocalBranchName("branch-3"), currentBranch)
			})
		})

		t.Run("recognizes the branch sync status", func(t *testing.T) {
			t.Parallel()
			t.Run("branch is ahead of its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1: ahead 1] Commit message 1a
  remotes/origin/branch-1      222222 Commit message 1b`[1:]
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
						SyncStatus: git.SyncStatusAhead,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("222222"),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})

			t.Run("branch is behind its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1: behind 2] Commit message 1
  remotes/origin/branch-1      222222 Commit message 1b`[1:]
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
						SyncStatus: git.SyncStatusBehind,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("222222"),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})

			t.Run("branch is ahead and behind its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1: ahead 31, behind 2] Commit message 1a
  remotes/origin/branch-1      222222 Commit message 1b`[1:]
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
						SyncStatus: git.SyncStatusAheadAndBehind,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("222222"),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})

			t.Run("branch is in sync with its remote branch", func(t *testing.T) {
				t.Parallel()
				give := `
  branch-1                     111111 [origin/branch-1] Commit message 1
  remotes/origin/branch-1      111111 Commit message 1`[1:]
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
						SyncStatus: git.SyncStatusUpToDate,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})

			t.Run("remote-only branch", func(t *testing.T) {
				t.Parallel()
				give := `
  remotes/origin/branch-1    222222 Commit message 2`[1:]
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.LocalBranchName{},
						InitialSHA: domain.SHA{},
						SyncStatus: git.SyncStatusRemoteOnly,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("222222"),
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})

			t.Run("local-only branch", func(t *testing.T) {
				t.Parallel()
				give := `  branch-1                     01a7eded Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("01a7eded"),
						SyncStatus: git.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})

			t.Run("branch is deleted at the remote", func(t *testing.T) {
				t.Parallel()
				give := `  branch-1                     01a7eded [origin/branch-1: gone] Commit message 1`
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("01a7eded"),
						SyncStatus: git.SyncStatusDeletedAtRemote,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.SHA{},
					},
				}
				have, _ := git.ParseVerboseBranchesOutput(give)
				assert.Equal(t, want, have)
			})
		})

		t.Run("branch with a different tracking branch name", func(t *testing.T) {
			t.Run("a branch uses a differently named tracking branch", func(t *testing.T) {
				give := `
  branch-1                     111111 [origin/branch-2] Commit message 1
  remotes/origin/branch-1      222222 Commit message 2
  remotes/origin/branch-2      111111 Commit message 1`[1:]
				want := git.BranchesSyncStatus{
					git.BranchSyncStatus{
						Name:       domain.NewLocalBranchName("branch-1"),
						InitialSHA: domain.NewSHA("111111"),
						SyncStatus: git.SyncStatusUpToDate,
						RemoteName: domain.NewRemoteBranchName("origin/branch-2"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
					git.BranchSyncStatus{
						Name:       domain.LocalBranchName{},
						InitialSHA: domain.SHA{},
						SyncStatus: git.SyncStatusRemoteOnly,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("222222"),
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
  main                         024df944 [origin/main] Commit message on main (#1234)
  branch-4                     e4d6bc09 [origin/branch-4: gone] Commit message 4
  remotes/origin/branch-1      307a7bf4 Commit message 1b
  remotes/origin/branch-2      da796a69 Commit message 2
  remotes/origin/branch-3      bc39378a Commit message 3b
  remotes/origin/HEAD          -> origin/initial
  remotes/origin/main          024df944 Commit message on main (#1234)
`[1:]
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:       domain.NewLocalBranchName("branch-1"),
					InitialSHA: domain.NewSHA("01a7eded"),
					SyncStatus: git.SyncStatusAhead,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("307a7bf4"),
				},
				git.BranchSyncStatus{
					Name:       domain.NewLocalBranchName("branch-2"),
					InitialSHA: domain.NewSHA("da796a69"),
					SyncStatus: git.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-2"),
					RemoteSHA:  domain.NewSHA("da796a69"),
				},
				git.BranchSyncStatus{
					Name:       domain.NewLocalBranchName("branch-3"),
					InitialSHA: domain.NewSHA("f4ebec0a"),
					SyncStatus: git.SyncStatusBehind,
					RemoteName: domain.NewRemoteBranchName("origin/branch-3"),
					RemoteSHA:  domain.NewSHA("bc39378a"),
				},
				git.BranchSyncStatus{
					Name:       domain.NewLocalBranchName("main"),
					InitialSHA: domain.NewSHA("024df944"),
					SyncStatus: git.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/main"),
					RemoteSHA:  domain.NewSHA("024df944"),
				},
				git.BranchSyncStatus{
					Name:       domain.NewLocalBranchName("branch-4"),
					InitialSHA: domain.NewSHA("e4d6bc09"),
					SyncStatus: git.SyncStatusDeletedAtRemote,
					RemoteName: domain.NewRemoteBranchName("origin/branch-4"),
					RemoteSHA:  domain.SHA{},
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, domain.NewLocalBranchName("branch-2"), currentBranch)
		})
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch(domain.NewLocalBranchName("feature1"), domain.NewLocalBranchName("initial"))
		runtime.CreateBranch(domain.NewLocalBranchName("feature2"), domain.NewLocalBranchName("initial"))
		runtime.CheckoutBranch(domain.NewLocalBranchName("feature1"))
		runtime.CheckoutBranch(domain.NewLocalBranchName("feature2"))
		have := runtime.Backend.PreviouslyCheckedOutBranch()
		assert.Equal(t, domain.NewLocalBranchName("feature1"), have)
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
				CurrentBranchCache: &cache.LocalBranch{},
				RemoteBranchCache:  &cache.RemoteBranch{},
				RemotesCache:       &cache.Strings{},
			}
			have := cmds.RootDirectory()
			assert.Empty(t, have)
		})
	})
}
