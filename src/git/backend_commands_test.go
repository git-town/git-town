package git_test

import (
	"strings"
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
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
		assert.True(t, runtime.Config.IsFeatureBranch("f1"))
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

	t.Run(".HasRemote()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		has, err := runner.HasOrigin()
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = runner.Backend.HasRemote("zonk")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasTrackingBranch()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		origin.CreateBranch("b1", "initial")
		repoDir := t.TempDir()
		devRepo := testruntime.Clone(origin.TestRunner, repoDir)
		devRepo.CheckoutBranch("b1")
		devRepo.CreateBranch("b2", "initial")
		has, err := devRepo.Backend.HasTrackingBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = devRepo.Backend.HasTrackingBranch("b2")
		assert.NoError(t, err)
		assert.False(t, has)
		has, err = devRepo.Backend.HasTrackingBranch("b3")
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

	t.Run(".LocalAndOriginBranches()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch("b1", "initial")
		runner.CreateBranch("b2", "initial")
		origin.CreateBranch("b3", "initial")
		runner.Fetch()
		branches, err := runner.Backend.LocalAndOriginBranches("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2", "b3"}, branches)
	})

	t.Run("ParseVerboseBranchesOutput", func(t *testing.T) {
		t.Parallel()
		t.Run("recognizes branches that are ahead of their remote branch", func(t *testing.T) {
			give := strings.TrimPrefix(`
* branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusAhead,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-1", currentBranch)
		})
		t.Run("recognizes branches that are behind their remote branch", func(t *testing.T) {
			give := strings.TrimPrefix(`
* branch-1                     01a7eded [origin/branch-1: behind 2] Commit message 1
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusBehind,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-1", currentBranch)
		})
		t.Run("recognizes branches that are in sync with their remote branch", func(t *testing.T) {
			give := strings.TrimPrefix(`
* branch-1                     01a7eded [origin/branch-1] Commit message 1
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusUpToDate,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-1", currentBranch)
		})
		t.Run("recognizes remote-only branches", func(t *testing.T) {
			give := strings.TrimPrefix(`
  remotes/origin/branch-1                     01a7eded Commit message 1
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusRemoteOnly,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "", currentBranch)
		})
		t.Run("recognizes local-only branches", func(t *testing.T) {
			give := strings.TrimPrefix(`
* branch-1                     01a7eded Commit message 1
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusLocalOnly,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-1", currentBranch)
		})
		t.Run("recognizes branches that got deleted at the remote", func(t *testing.T) {
			give := strings.TrimPrefix(`
* branch-1                     01a7eded [origin/branch-1: gone] Commit message 1
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusDeletedAtRemote,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-1", currentBranch)
		})
		t.Run("complex example", func(t *testing.T) {
			give := strings.TrimPrefix(`
* branch-1                     01a7eded [origin/branch-1: ahead 1] Commit message 1
  branch-2                     da796a69 [origin/branch-2] Commit message 2
  branch-3                     f4ebec0a [origin/branch-3: behind 2] Commit message 3
  main                         024df944 [origin/main] Commit message on main (#1234)
  branch-4                     e4d6bc09 [origin/branch-4: gone] Commit message 4
  remotes/origin/branch-5      307a7bf4 Commit message 5
`, "\n")
			want := git.BranchesWithSyncStatus{
				git.BranchWithSyncStatus{
					Name:       "branch-1",
					SyncStatus: git.SyncStatusAhead,
				},
				git.BranchWithSyncStatus{
					Name:       "branch-2",
					SyncStatus: git.SyncStatusUpToDate,
				},
				git.BranchWithSyncStatus{
					Name:       "branch-3",
					SyncStatus: git.SyncStatusBehind,
				},
				git.BranchWithSyncStatus{
					Name:       "main",
					SyncStatus: git.SyncStatusUpToDate,
				},
				git.BranchWithSyncStatus{
					Name:       "branch-4",
					SyncStatus: git.SyncStatusDeletedAtRemote,
				},
				git.BranchWithSyncStatus{
					Name:       "branch-5",
					SyncStatus: git.SyncStatusRemoteOnly,
				},
			}
			have, currentBranch := git.ParseVerboseBranchesOutput(give)
			assert.Equal(t, want, have)
			assert.Equal(t, "branch-1", currentBranch)
		})
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("feature1", "initial")
		runtime.CreateBranch("feature2", "initial")
		runtime.CheckoutBranch("feature1")
		runtime.CheckoutBranch("feature2")
		have, err := runtime.Backend.PreviouslyCheckedOutBranch()
		assert.NoError(t, err)
		assert.Equal(t, "feature1", have)
	})

	t.Run(".RemoteBranches()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		runner.CreateBranch("b1", "initial")
		runner.CreateBranch("b2", "initial")
		origin.CreateBranch("b3", "initial")
		runner.Fetch()
		branches, err := runner.Backend.RemoteBranches()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin/b3", "origin/initial"}, branches)
	})

	t.Run(".Remotes()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		runtime.AddRemote(config.OriginRemote, origin.WorkingDir)
		remotes, err := runtime.Backend.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{config.OriginRemote}, remotes)
	})
}
