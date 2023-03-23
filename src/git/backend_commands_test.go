package git_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/test"
	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	t.Parallel()

	t.Run("BranchAuthors", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		err := repo.CreateBranch("branch", "initial")
		assert.NoError(t, err)
		err = repo.CreateCommit(git.Commit{
			Branch:      "branch",
			FileName:    "file1",
			FileContent: "file1",
			Message:     "first commit",
		})
		assert.NoError(t, err)
		err = repo.CreateCommit(git.Commit{
			Branch:      "branch",
			FileName:    "file2",
			FileContent: "file2",
			Message:     "second commit",
		})
		assert.NoError(t, err)
		authors, err := repo.Backend.BranchAuthors("branch", "initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"user <email@example.com>"}, authors)
	})

	t.Run(".CheckoutBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		err := repo.CreateBranch("branch1", "initial")
		assert.NoError(t, err)
		err = repo.Backend.CheckoutBranch("branch1")
		assert.NoError(t, err)
		currentBranch, err := repo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "branch1", currentBranch)
		err = repo.CheckoutBranch("initial")
		assert.NoError(t, err)
		currentBranch, err = repo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", currentBranch)
	})

	t.Run(".CreateFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t)
		err := repo.Backend.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		repo.Config.Reload()
		assert.True(t, repo.Config.IsFeatureBranch("f1"))
		assert.Equal(t, []string{"main"}, repo.Config.AncestorBranches("f1"))
	})

	t.Run(".CurrentBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		err := repo.CheckoutBranch("initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("b1")
		assert.NoError(t, err)
		branch, err := repo.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "b1", branch)
		err = repo.CheckoutBranch("initial")
		assert.NoError(t, err)
		branch, err = repo.Backend.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch)
	})

	t.Run(".HasLocalBranch()", func(t *testing.T) {
		t.Parallel()
		origin := test.CreateRepo(t)
		repoDir := t.TempDir()
		repo, err := origin.Clone(repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		has, err := repo.Backend.HasLocalBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.Backend.HasLocalBranch("b2")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.Backend.HasLocalBranch("b3")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasOpenChanges()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		has, err := repo.Backend.HasOpenChanges()
		assert.NoError(t, err)
		assert.False(t, has)
		err = repo.CreateFile("foo", "bar")
		assert.NoError(t, err)
		has, err = repo.Backend.HasOpenChanges()
		assert.NoError(t, err)
		assert.True(t, has)
	})

	t.Run(".HasRebaseInProgress()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		has, err := repo.Backend.HasRebaseInProgress()
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasRemote()", func(t *testing.T) {
		t.Parallel()
		origin := test.CreateRepo(t)
		repoDir := t.TempDir()
		repo, err := origin.Clone(repoDir)
		assert.NoError(t, err)
		has, err := repo.HasOrigin()
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.Backend.HasRemote("zonk")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasTrackingBranch()", func(t *testing.T) {
		t.Parallel()
		origin := test.CreateRepo(t)
		err := origin.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		repoDir := t.TempDir()
		devRepo, err := origin.Clone(repoDir)
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
		origin := test.CreateRepo(t)
		repoDir := t.TempDir()
		repo, err := origin.Clone(repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		err = origin.CreateBranch("b3", "initial")
		assert.NoError(t, err)
		err = repo.Fetch()
		assert.NoError(t, err)
		branches, err := repo.Backend.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2"}, branches)
	})

	t.Run(".LocalAndOriginBranches()", func(t *testing.T) {
		t.Parallel()
		origin := test.CreateRepo(t)
		repoDir := t.TempDir()
		repo, err := origin.Clone(repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		err = origin.CreateBranch("b3", "initial")
		assert.NoError(t, err)
		err = repo.Fetch()
		assert.NoError(t, err)
		branches, err := repo.Backend.LocalAndOriginBranches("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1", "b2", "b3"}, branches)
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		err := repo.CreateBranch("feature1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("feature2", "initial")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("feature1")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("feature2")
		assert.NoError(t, err)
		have, err := repo.Backend.PreviouslyCheckedOutBranch()
		assert.NoError(t, err)
		assert.Equal(t, "feature1", have)
	})

	t.Run(".RemoteBranches()", func(t *testing.T) {
		t.Parallel()
		origin := test.CreateRepo(t)
		repoDir := t.TempDir()
		repo, err := origin.Clone(repoDir)
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		err = origin.CreateBranch("b3", "initial")
		assert.NoError(t, err)
		err = repo.Fetch()
		assert.NoError(t, err)
		branches, err := repo.Backend.RemoteBranches()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin/b3", "origin/initial"}, branches)
	})

	t.Run(".Remotes()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t)
		origin := test.CreateRepo(t)
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		remotes, err := repo.Backend.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{config.OriginRemote}, remotes)
	})
}
