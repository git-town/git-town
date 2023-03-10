package git_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/test"
	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	t.Parallel()
	t.Run(".AddRemote()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		remotes, err := repo.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, remotes)
		origin := test.CreateRepo(t)
		err = repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		remotes, err = repo.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})

	t.Run(".CheckoutBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateBranch("branch1", "initial")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("branch1")
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

	t.Run(".Commits()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file1",
			FileContent: "hello",
			Message:     "first commit",
		})
		assert.NoError(t, err)
		err = repo.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file2",
			FileContent: "hello again",
			Message:     "second commit",
		})
		assert.NoError(t, err)
		commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
		assert.NoError(t, err)
		assert.Len(t, commits, 2)
		assert.Equal(t, "initial", commits[0].Branch)
		assert.Equal(t, "file1", commits[0].FileName)
		assert.Equal(t, "hello", commits[0].FileContent)
		assert.Equal(t, "first commit", commits[0].Message)
		assert.Equal(t, "initial", commits[1].Branch)
		assert.Equal(t, "file2", commits[1].FileName)
		assert.Equal(t, "hello again", commits[1].FileContent)
		assert.Equal(t, "second commit", commits[1].Message)
	})

	t.Run(".ConnectTrackingBranch()", func(t *testing.T) {
		t.Parallel()
		// replicating the situation this is used in,
		// connecting branches of repos with the same commits in them
		origin := test.CreateRepo(t)
		repoDir := filepath.Join(t.TempDir(), "repo") // need a non-existing directory
		err := test.CopyDirectory(origin.WorkingDir(), repoDir)
		assert.NoError(t, err)
		repo := test.NewRepo(repoDir, repoDir, "").Repo
		err = repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		err = repo.Fetch()
		assert.NoError(t, err)
		err = repo.ConnectTrackingBranch("initial")
		assert.NoError(t, err)
		err = repo.PushBranch()
		assert.NoError(t, err)
	})

	t.Run(".CreateBranch()", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			repo := test.CreateRepo(t).Repo
			err := repo.CreateBranch("branch1", "initial")
			assert.NoError(t, err)
			currentBranch, err := repo.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, "initial", currentBranch)
			branches, err := repo.LocalBranchesMainFirst()
			assert.NoError(t, err)
			assert.Equal(t, []string{"branch1", "initial"}, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			repo := test.CreateRepo(t).Repo
			err := repo.CreateBranch("my/feature", "initial")
			assert.NoError(t, err)
			currentBranch, err := repo.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, "initial", currentBranch)
			branches, err := repo.LocalBranchesMainFirst()
			assert.NoError(t, err)
			assert.Equal(t, []string{"initial", "my/feature"}, branches)
		})
	})

	t.Run(".CreateChildFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t).Repo
		err := repo.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		err = repo.CreateChildFeatureBranch("f1a", "f1")
		assert.NoError(t, err)
		res, err := repo.Run("git", "town", "config")
		assert.NoError(t, err)
		has := strings.Contains(res.OutputSanitized(), "Branch Ancestry:\n  main\n    f1\n      f1a")
		assert.True(t, has)
	})

	t.Run(".CreateCommit()", func(t *testing.T) {
		t.Run("minimal arguments", func(t *testing.T) {
			t.Parallel()
			repo := test.CreateRepo(t).Repo
			err := repo.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			assert.NoError(t, err)
			commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
			assert.NoError(t, err)
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, "initial", commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			repo := test.CreateRepo(t).Repo
			err := repo.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
				Author:      "developer <developer@example.com>",
			})
			assert.NoError(t, err)
			commits, err := repo.Commits([]string{"FILE NAME", "FILE CONTENT"})
			assert.NoError(t, err)
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, "initial", commits[0].Branch)
			assert.Equal(t, "developer <developer@example.com>", commits[0].Author)
		})
	})

	t.Run(".CreateFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t).Repo
		err := repo.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		repo.Config.Reload()
		assert.True(t, repo.Config.IsFeatureBranch("f1"))
		assert.Equal(t, []string{"main"}, repo.Config.AncestorBranches("f1"))
	})

	t.Run(".CreateFeatureBranchNoParent()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t).Repo
		err := repo.CreateFeatureBranchNoParent("f1")
		assert.NoError(t, err)
		repo.Config.Reload()
		assert.True(t, repo.Config.IsFeatureBranch("f1"))
		assert.Equal(t, []string{}, repo.Config.AncestorBranches("f1"))
	})

	t.Run(".CreateFile()", func(t *testing.T) {
		t.Run("simple example", func(t *testing.T) {
			t.Parallel()
			repo := test.CreateRepo(t).Repo
			err := repo.CreateFile("filename", "content")
			assert.Nil(t, err, "cannot create file in repo")
			content, err := os.ReadFile(filepath.Join(repo.WorkingDir(), "filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})

		t.Run("create file in subfolder", func(t *testing.T) {
			t.Parallel()
			repo := test.CreateRepo(t).Repo
			err := repo.CreateFile("folder/filename", "content")
			assert.Nil(t, err, "cannot create file in repo")
			content, err := os.ReadFile(filepath.Join(repo.WorkingDir(), "folder/filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})
	})

	t.Run(".CreatePerennialBranches()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t).Repo
		err := repo.CreatePerennialBranches("p1", "p2")
		assert.NoError(t, err)
		branches, err := repo.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.Equal(t, []string{"main", "initial", "p1", "p2"}, branches)
		repo.Config.Reload()
		assert.True(t, repo.Config.IsPerennialBranch("p1"))
		assert.True(t, repo.Config.IsPerennialBranch("p2"))
	})

	t.Run(".CurrentBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CheckoutBranch("initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("b1")
		assert.NoError(t, err)
		branch, err := repo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "b1", branch)
		err = repo.CheckoutBranch("initial")
		assert.NoError(t, err)
		branch, err = repo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch)
	})

	t.Run(".Fetch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		origin := test.CreateRepo(t)
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		err = repo.Fetch()
		assert.NoError(t, err)
	})

	t.Run(".FileContentInCommit()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		assert.NoError(t, err)
		commits, err := repo.CommitsInBranch("initial", []string{})
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		content, err := repo.FileContentInCommit(commits[0].SHA, "hello.txt")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", content)
	})

	t.Run(".FilesInCommit()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		err = repo.CreateFile("f2.txt", "two")
		assert.NoError(t, err)
		err = repo.StageFiles("f1.txt", "f2.txt")
		assert.NoError(t, err)
		err = repo.CommitStagedChanges("stuff")
		assert.NoError(t, err)
		commits, err := repo.Commits([]string{})
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		fileNames, err := repo.FilesInCommit(commits[0].SHA)
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run(".HasBranchesOutOfSync()", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env, err := test.NewStandardGitEnvironment(t.TempDir())
			assert.NoError(t, err)
			repo := env.DevRepo.Repo
			err = repo.CreateBranch("branch1", "main")
			assert.NoError(t, err)
			err = repo.CheckoutBranch("branch1")
			assert.NoError(t, err)
			err = repo.CreateFile("file1", "content")
			assert.NoError(t, err)
			err = repo.StageFiles("file1")
			assert.NoError(t, err)
			err = repo.CommitStagedChanges("stuff")
			assert.NoError(t, err)
			err = repo.PushBranch(git.PushArgs{Branch: "branch1", Remote: config.OriginRemote})
			assert.NoError(t, err)
			have, err := repo.HasBranchesOutOfSync()
			assert.NoError(t, err)
			assert.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			env, err := test.NewStandardGitEnvironment(t.TempDir())
			assert.NoError(t, err)
			repo := env.DevRepo.Repo
			err = repo.CreateBranch("branch1", "main")
			assert.NoError(t, err)
			err = repo.PushBranch()
			assert.NoError(t, err)
			err = repo.CreateFile("file1", "content")
			assert.NoError(t, err)
			err = repo.StageFiles("file1")
			assert.NoError(t, err)
			err = repo.CommitStagedChanges("stuff")
			assert.NoError(t, err)
			have, err := repo.HasBranchesOutOfSync()
			assert.NoError(t, err)
			assert.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			env, err := test.NewStandardGitEnvironment(t.TempDir())
			assert.NoError(t, err)
			err = env.DevRepo.CreateBranch("branch1", "main")
			assert.NoError(t, err)
			err = env.DevRepo.PushBranch()
			assert.NoError(t, err)
			err = env.OriginRepo.CheckoutBranch("main")
			assert.NoError(t, err)
			err = env.OriginRepo.CreateFile("file1", "content")
			assert.NoError(t, err)
			err = env.OriginRepo.StageFiles("file1")
			assert.NoError(t, err)
			err = env.OriginRepo.CommitStagedChanges("stuff")
			assert.NoError(t, err)
			err = env.OriginRepo.CheckoutBranch("initial")
			assert.NoError(t, err)
			err = env.DevRepo.Fetch()
			assert.NoError(t, err)
			have, err := env.DevRepo.Repo.HasBranchesOutOfSync()
			assert.NoError(t, err)
			assert.True(t, have)
		})
	})

	t.Run(".HasFile()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		has, err := repo.HasFile("f1.txt", "one")
		assert.NoError(t, err)
		assert.True(t, has)
		_, err = repo.HasFile("f1.txt", "zonk")
		assert.Error(t, err)
		_, err = repo.HasFile("zonk.txt", "one")
		assert.Error(t, err)
	})

	t.Run(".HasGitTownConfigNow()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		res, err := repo.HasGitTownConfigNow()
		assert.NoError(t, err)
		assert.False(t, res)
		err = repo.CreateBranch("main", "initial")
		assert.NoError(t, err)
		err = repo.CreateFeatureBranch("foo")
		assert.NoError(t, err)
		res, err = repo.HasGitTownConfigNow()
		assert.NoError(t, err)
		assert.True(t, res)
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
		has, err := repo.HasLocalBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.HasLocalBranch("b2")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.HasLocalBranch("b3")
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasOpenChanges()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		has, err := repo.HasOpenChanges()
		assert.NoError(t, err)
		assert.False(t, has)
		err = repo.CreateFile("foo", "bar")
		assert.NoError(t, err)
		has, err = repo.HasOpenChanges()
		assert.NoError(t, err)
		assert.True(t, has)
	})

	t.Run(".HasRebaseInProgress()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		has, err := repo.HasRebaseInProgress()
		assert.NoError(t, err)
		assert.False(t, has)
	})

	t.Run(".HasRemote()", func(t *testing.T) {
		t.Parallel()
		origin := test.CreateRepo(t)
		repoDir := t.TempDir()
		repo, err := origin.Clone(repoDir)
		assert.NoError(t, err)
		has, err := repo.Repo.HasOrigin()
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.Repo.HasRemote("zonk")
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
		repo := devRepo.Repo
		err = repo.CheckoutBranch("b1")
		assert.NoError(t, err)
		err = repo.CreateBranch("b2", "initial")
		assert.NoError(t, err)
		has, err := repo.HasTrackingBranch("b1")
		assert.NoError(t, err)
		assert.True(t, has)
		has, err = repo.HasTrackingBranch("b2")
		assert.NoError(t, err)
		assert.False(t, has)
		has, err = repo.HasTrackingBranch("b3")
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
		branches, err := repo.Repo.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.Equal(t, []string{"b1", "b2", "initial"}, branches)
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
		branches, err := repo.Repo.LocalAndOriginBranches()
		assert.NoError(t, err)
		assert.Equal(t, []string{"b1", "b2", "b3", "initial"}, branches)
	})

	t.Run(".PreviouslyCheckedOutBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateBranch("feature1", "initial")
		assert.NoError(t, err)
		err = repo.CreateBranch("feature2", "initial")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("feature1")
		assert.NoError(t, err)
		err = repo.CheckoutBranch("feature2")
		assert.NoError(t, err)
		have, err := repo.PreviouslyCheckedOutBranch()
		assert.NoError(t, err)
		assert.Equal(t, "feature1", have)
	})

	t.Run(".PushBranchToOrigin()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		origin := test.CreateRepo(t)
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		err = repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = repo.PushBranch(git.PushArgs{Branch: "b1", Remote: config.OriginRemote})
		assert.NoError(t, err)
		branches, err := origin.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.Equal(t, []string{"b1", "initial"}, branches)
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
		branches, err := repo.Repo.RemoteBranches()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin/b3", "origin/initial"}, branches)
	})

	t.Run(".Remotes()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		origin := test.CreateRepo(t)
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		remotes, err := repo.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{config.OriginRemote}, remotes)
	})

	t.Run(".RemoveBranch()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		branches, err := repo.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.Equal(t, []string{"b1", "initial"}, branches)
		err = repo.RemoveBranch("b1")
		assert.NoError(t, err)
		branches, err = repo.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial"}, branches)
	})

	t.Run(".RemoveRemote()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		origin := test.CreateRepo(t)
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir())
		assert.NoError(t, err)
		err = repo.RemoveRemote(config.OriginRemote)
		assert.NoError(t, err)
		remotes, err := repo.Remotes()
		assert.NoError(t, err)
		assert.Len(t, remotes, 0)
	})

	t.Run(".ShaForCommit()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateCommit(git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
		assert.NoError(t, err)
		sha, err := repo.ShaForCommit("commit")
		assert.NoError(t, err)
		assert.Len(t, sha, 40)
	})

	t.Run(".Stash()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		stashSize, err := repo.StashSize()
		assert.NoError(t, err)
		assert.Zero(t, stashSize)
		err = repo.CreateFile("f1.txt", "hello")
		assert.NoError(t, err)
		err = repo.Stash()
		assert.NoError(t, err)
		stashSize, err = repo.StashSize()
		assert.NoError(t, err)
		assert.Equal(t, 1, stashSize)
	})

	t.Run(".UncommittedFiles()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateRepo(t).Repo
		err := repo.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		err = repo.CreateFile("f2.txt", "two")
		assert.NoError(t, err)
		files, err := repo.UncommittedFiles()
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
