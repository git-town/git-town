package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Parallel()

	t.Run("BranchAuthors", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateBranch("branch", "initial")
		assert.NoError(t, err)
		err = runtime.CreateCommit(git.Commit{
			Branch:      "branch",
			FileName:    "file1",
			FileContent: "file1",
			Message:     "first commit",
		})
		assert.NoError(t, err)
		err = runtime.CreateCommit(git.Commit{
			Branch:      "branch",
			FileName:    "file2",
			FileContent: "file2",
			Message:     "second commit",
		})
		assert.NoError(t, err)
		authors, err := runtime.Backend.BranchAuthors("branch", "initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"user <email@example.com>"}, authors)
	})

	t.Run(".CheckoutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateBranch("branch1", "initial")
		assert.NoError(t, err)
		err = runtime.Backend.CheckoutBranch("branch1")
		assert.NoError(t, err)
		currentBranch, err := runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "branch1", currentBranch)
		err = runtime.CheckoutBranch("initial")
		assert.NoError(t, err)
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
		assert.Equal(t, []string{"main"}, runtime.Config.Lineage.Ancestors("f1"))
	})

	t.Run(".CurrentBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CheckoutBranch("initial")
		assert.NoError(t, err)
		err = runtime.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = runtime.CheckoutBranch("b1")
		assert.NoError(t, err)
		branch, err := runtime.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "b1", branch)
		err = runtime.CheckoutBranch("initial")
		assert.NoError(t, err)
		branch, err = runtime.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch)
	})

	t.Run(".HasLocalBranch()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner, err := testruntime.Clone(origin.TestRunner, repoDir)
		assert.NoError(t, err)
		err = runner.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = runner.CreateBranch("b2", "initial")
		assert.NoError(t, err)
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
		err = runtime.CreateFile("foo", "bar")
		assert.NoError(t, err)
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
		runner, err := testruntime.Clone(origin.TestRunner, repoDir)
		assert.NoError(t, err)
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
		err := origin.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		repoDir := t.TempDir()
		devRepo, err := testruntime.Clone(origin.TestRunner, repoDir)
		assert.NoError(t, err)
		err = devRepo.CheckoutBranch("b1")
		assert.NoError(t, err)
		err = devRepo.CreateBranch("b2", "initial")
		assert.NoError(t, err)
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
		runner, err := testruntime.Clone(origin.TestRunner, repoDir)
		assert.NoError(t, err)
		err = runner.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = runner.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		err = origin.CreateBranch("b3", "initial")
		assert.NoError(t, err)
		err = runner.Fetch()
		assert.NoError(t, err)
		branches, err := runner.Backend.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2"}, branches)
	})

	t.Run(".LocalAndOriginBranches()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner, err := testruntime.Clone(origin.TestRunner, repoDir)
		assert.NoError(t, err)
		err = runner.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = runner.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		err = origin.CreateBranch("b3", "initial")
		assert.NoError(t, err)
		err = runner.Fetch()
		assert.NoError(t, err)
		branches, err := runner.Backend.LocalAndOriginBranches("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2", "b3"}, branches)
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateBranch("feature1", "initial")
		assert.NoError(t, err)
		err = runtime.CreateBranch("feature2", "initial")
		assert.NoError(t, err)
		err = runtime.CheckoutBranch("feature1")
		assert.NoError(t, err)
		err = runtime.CheckoutBranch("feature2")
		assert.NoError(t, err)
		have, err := runtime.Backend.PreviouslyCheckedOutBranch()
		assert.NoError(t, err)
		assert.Equal(t, "feature1", have)
	})

	t.Run(".RemoteBranches()", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner, err := testruntime.Clone(origin.TestRunner, repoDir)
		assert.NoError(t, err)
		err = runner.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = runner.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		err = origin.CreateBranch("b3", "initial")
		assert.NoError(t, err)
		err = runner.Fetch()
		assert.NoError(t, err)
		branches, err := runner.Backend.RemoteBranches()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin/b3", "origin/initial"}, branches)
	})

	t.Run(".Remotes()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		err := runtime.AddRemote(config.OriginRemote, origin.WorkingDir)
		assert.NoError(t, err)
		remotes, err := runtime.Backend.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{config.OriginRemote}, remotes)
	})
}
