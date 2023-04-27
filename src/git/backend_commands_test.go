package git_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/fs"
	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Parallel()

	t.Run("BranchAuthors", func(t *testing.T) {
		t.Parallel()
		dev := repo.Create(t)
		err := repo.CreateBranch(&dev, "branch", "initial")
		assert.NoError(t, err)
		err = repo.CreateCommit(&dev, git.Commit{
			Branch:      "branch",
			FileName:    "file1",
			FileContent: "file1",
			Message:     "first commit",
		})
		assert.NoError(t, err)
		err = repo.CreateCommit(&dev, git.Commit{
			Branch:      "branch",
			FileName:    "file2",
			FileContent: "file2",
			Message:     "second commit",
		})
		assert.NoError(t, err)
		authors, err := dev.BranchAuthors("branch", "initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"user <email@example.com>"}, authors)
	})

	t.Run(".CheckoutBranch()", func(t *testing.T) {
		t.Parallel()
		dev := repo.Create(t)
		err := repo.CreateBranch(&dev, "branch1", "initial")
		assert.NoError(t, err)
		err = dev.CheckoutBranch("branch1")
		assert.NoError(t, err)
		currentBranch, err := dev.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "branch1", currentBranch)
		err = dev.CheckoutBranch("initial")
		assert.NoError(t, err)
		currentBranch, err = dev.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", currentBranch)
	})

	t.Run(".CreateFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		repo := repo.CreateGitTown(t)
		err := repo.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		repo.Reload()
		assert.True(t, repo.IsFeatureBranch("f1"))
		assert.Equal(t, []string{"main"}, repo.AncestorBranches("f1"))
	})

	t.Run(".CurrentBranch()", func(t *testing.T) {
		t.Parallel()
		dev := repo.Create(t)
		err := dev.CheckoutBranch("initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&dev, "b1", "initial")
		assert.NoError(t, err)
		err = dev.CheckoutBranch("b1")
		assert.NoError(t, err)
		branch, err := dev.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "b1", branch)
		err = dev.CheckoutBranch("initial")
		assert.NoError(t, err)
		branch, err = dev.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch)
	})

	t.Run(".HasLocalBranch()", func(t *testing.T) {
		t.Parallel()
		origin := repo.Create(t)
		repoDir := t.TempDir()
		dev, err := repo.Clone(&origin, repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch(&dev, "b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&dev, "b2", "initial")
		assert.NoError(t, err)
		has, err := dev.HasLocalBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = dev.HasLocalBranch("b2")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = dev.HasLocalBranch("b3")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasOpenChanges()", func(t *testing.T) {
		t.Parallel()
		repo := repo.Create(t)
		has, err := repo.HasOpenChanges()
		assert.NoError(t, err)
		assert.False(t, has)
		err = fs.CreateFile(repo.Dir(), "foo", "bar")
		assert.NoError(t, err)
		has, err = repo.HasOpenChanges()
		assert.NoError(t, err)
		assert.True(t, has)
	})

	t.Run(".HasRebaseInProgress()", func(t *testing.T) {
		t.Parallel()
		repo := repo.Create(t)
		has, err := repo.HasRebaseInProgress()
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasRemote()", func(t *testing.T) {
		t.Parallel()
		origin := repo.Create(t)
		repoDir := t.TempDir()
		runner, err := repo.Clone(&origin, repoDir)
		assert.NoError(t, err)
		has, err := runner.HasOrigin()
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = runner.HasRemote("zonk")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasTrackingBranch()", func(t *testing.T) {
		t.Parallel()
		origin := repo.Create(t)
		err := repo.CreateBranch(&origin, "b1", "initial")
		assert.NoError(t, err)
		repoDir := t.TempDir()
		devRepo, err := repo.Clone(&origin, repoDir)
		assert.NoError(t, err)
		err = devRepo.CheckoutBranch("b1")
		assert.NoError(t, err)
		err = repo.CreateBranch(&devRepo, "b2", "initial")
		assert.NoError(t, err)
		has, err := devRepo.HasTrackingBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = devRepo.HasTrackingBranch("b2")
		assert.NoError(t, err)
		assert.False(t, has)
		has, err = devRepo.HasTrackingBranch("b3")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".LocalBranchesMainFirst()", func(t *testing.T) {
		t.Parallel()
		origin := repo.Create(t)
		repoDir := t.TempDir()
		runner, err := repo.Clone(&origin, repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch(&runner, "b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&runner, "b2", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&origin, "b3", "initial")
		assert.NoError(t, err)
		err = repo.Fetch(&runner)
		assert.NoError(t, err)
		branches, err := runner.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2"}, branches)
	})

	t.Run(".LocalAndOriginBranches()", func(t *testing.T) {
		t.Parallel()
		origin := repo.Create(t)
		repoDir := t.TempDir()
		runner, err := repo.Clone(&origin, repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch(&runner, "b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&runner, "b2", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&origin, "b3", "initial")
		assert.NoError(t, err)
		err = repo.Fetch(&runner)
		assert.NoError(t, err)
		branches, err := runner.LocalAndOriginBranches("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2", "b3"}, branches)
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		dev := repo.Create(t)
		err := repo.CreateBranch(&dev, "feature1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&dev, "feature2", "initial")
		assert.NoError(t, err)
		err = dev.CheckoutBranch("feature1")
		assert.NoError(t, err)
		err = dev.CheckoutBranch("feature2")
		assert.NoError(t, err)
		have, err := dev.PreviouslyCheckedOutBranch()
		assert.NoError(t, err)
		assert.Equal(t, "feature1", have)
	})

	t.Run(".RemoteBranches()", func(t *testing.T) {
		t.Parallel()
		origin := repo.Create(t)
		repoDir := t.TempDir()
		runner, err := repo.Clone(&origin, repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch(&runner, "b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&runner, "b2", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch(&origin, "b3", "initial")
		assert.NoError(t, err)
		err = repo.Fetch(&runner)
		assert.NoError(t, err)
		branches, err := runner.RemoteBranches()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin/b3", "origin/initial"}, branches)
	})

	t.Run(".Remotes()", func(t *testing.T) {
		t.Parallel()
		dev := repo.Create(t)
		origin := t.TempDir()
		err := repo.AddRemote(&dev, config.OriginRemote, origin)
		assert.NoError(t, err)
		remotes, err := dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{config.OriginRemote}, remotes)
	})
}
