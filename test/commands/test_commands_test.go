package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/test/fixture"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestTestCommands(t *testing.T) {
	t.Parallel()

	t.Run(".AddRemote()", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		remotes, err := dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, remotes)
		origin := testruntime.Create(t)
		dev.AddRemoteX(config.OriginRemote, origin.WorkingDir)
		remotes, err = dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})

	t.Run(".Commits()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommitX(git.Commit{
			Branch:      "initial",
			FileName:    "file1",
			FileContent: "hello",
			Message:     "first commit",
		})
		runtime.CreateCommitX(git.Commit{
			Branch:      "initial",
			FileName:    "file2",
			FileContent: "hello again",
			Message:     "second commit",
		})
		commits := runtime.CommitsX([]string{"FILE NAME", "FILE CONTENT"}, "initial")
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
		origin := testruntime.Create(t)
		repoDir := filepath.Join(t.TempDir(), "repo") // need a non-existing directory
		helpers.CopyDirectory(origin.WorkingDir, repoDir)
		runtime := testruntime.New(repoDir, repoDir, "")
		runtime.AddRemoteX(config.OriginRemote, origin.WorkingDir)
		runtime.FetchX()
		runtime.ConnectTrackingBranchX("initial")
		runtime.PushBranchX()
	})

	t.Run(".CreateBranch()", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranchX("branch1", "initial")
			currentBranch, err := runtime.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, "initial", currentBranch)
			branches, err := runtime.LocalBranchesMainFirst("initial")
			assert.NoError(t, err)
			assert.Equal(t, []string{"initial", "branch1"}, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranchX("my/feature", "initial")
			currentBranch, err := runtime.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, "initial", currentBranch)
			branches, err := runtime.LocalBranchesMainFirst("initial")
			assert.NoError(t, err)
			assert.Equal(t, []string{"initial", "my/feature"}, branches)
		})
	})

	t.Run(".CreateChildFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		err := runtime.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		runtime.CreateChildFeatureBranchX("f1a", "f1")
		output, err := runtime.BackendRunner.Query("git-town", "config")
		assert.NoError(t, err)
		output = stripansi.Strip(output)
		if !strings.Contains(output, "Branch Lineage:\n  main\n    f1\n      f1a") {
			t.Fatalf("unexpected output: %s", output)
		}
	})

	t.Run(".CreateCommit()", func(t *testing.T) {
		t.Run("minimal arguments", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommitX(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			commits := runtime.CommitsX([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, "initial", commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommitX(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
				Author:      "developer <developer@example.com>",
			})
			commits := runtime.CommitsX([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, "initial", commits[0].Branch)
			assert.Equal(t, "developer <developer@example.com>", commits[0].Author)
		})
	})

	t.Run(".CreateFile()", func(t *testing.T) {
		t.Run("simple example", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFileX("filename", "content")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})

		t.Run("create file in subfolder", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFileX("folder/filename", "content")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "folder/filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})
	})

	t.Run(".CreatePerennialBranches()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreatePerennialBranchesX("p1", "p2")
		branches, err := runtime.LocalBranchesMainFirst("main")
		assert.NoError(t, err)
		assert.Equal(t, []string{"main", "initial", "p1", "p2"}, branches)
		runtime.Config.Reload()
		assert.True(t, runtime.Config.IsPerennialBranch("p1"))
		assert.True(t, runtime.Config.IsPerennialBranch("p2"))
	})

	t.Run(".Fetch()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemoteX(config.OriginRemote, origin.WorkingDir)
		repo.FetchX()
	})

	t.Run(".FileContentInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommitX(git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		commits := runtime.CommitsInBranchX("initial", []string{})
		assert.Len(t, commits, 1)
		content := runtime.FileContentInCommitX(commits[0].SHA, "hello.txt")
		assert.Equal(t, "hello world", content)
	})

	t.Run(".FilesInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFileX("f1.txt", "one")
		runtime.CreateFileX("f2.txt", "two")
		runtime.StageFilesX("f1.txt", "f2.txt")
		runtime.CommitStagedChangesX("stuff")
		commits := runtime.CommitsX([]string{}, "initial")
		assert.Len(t, commits, 1)
		fileNames := runtime.FilesInCommitX(commits[0].SHA)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run(".HasBranchesOutOfSync()", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			runner := env.DevRepo
			runner.CreateBranchX("branch1", "main")
			runner.CheckoutBranchX("branch1")
			runner.CreateFileX("file1", "content")
			runner.StageFilesX("file1")
			runner.CommitStagedChangesX("stuff")
			runner.PushBranchToRemoteX("branch1", config.OriginRemote)
			have := runner.HasBranchesOutOfSyncX()
			assert.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			env.DevRepo.CreateBranchX("branch1", "main")
			env.DevRepo.PushBranchX()
			env.DevRepo.CreateFileX("file1", "content")
			env.DevRepo.StageFilesX("file1")
			env.DevRepo.CommitStagedChangesX("stuff")
			have := env.DevRepo.HasBranchesOutOfSyncX()
			assert.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			env.DevRepo.CreateBranchX("branch1", "main")
			env.DevRepo.PushBranchX()
			env.OriginRepo.CheckoutBranchX("main")
			env.OriginRepo.CreateFileX("file1", "content")
			env.OriginRepo.StageFilesX("file1")
			env.OriginRepo.CommitStagedChangesX("stuff")
			env.OriginRepo.CheckoutBranchX("initial")
			env.DevRepo.FetchX()
			have := env.DevRepo.HasBranchesOutOfSyncX()
			assert.True(t, have)
		})
	})

	t.Run(".HasFile()", func(t *testing.T) {
		t.Parallel()
		t.Run("filename and content match", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFileX("f1.txt", "one")
			assert.Equal(t, "", runtime.HasFile("f1.txt", "one"))
		})
		t.Run("filename matches, content doesn't match", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFileX("f1.txt", "one")
			assert.Equal(t, "file \"f1.txt\" should have content \"zonk\" but has \"one\"", runtime.HasFile("f1.txt", "zonk"))
		})
		t.Run("filename doesn't match but content matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFileX("f1.txt", "one")
			assert.Equal(t, "repo doesn't have file \"zonk.txt\"", runtime.HasFile("zonk.txt", "one"))
		})
	})

	t.Run(".HasGitTownConfigNow()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		res := runtime.HasGitTownConfigNow()
		assert.False(t, res)
		runtime.CreateBranchX("main", "initial")
		err := runtime.CreateFeatureBranch("foo")
		assert.NoError(t, err)
		res = runtime.HasGitTownConfigNow()
		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run(".PushBranchToRemote()", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		origin := testruntime.Create(t)
		dev.AddRemoteX(config.OriginRemote, origin.WorkingDir)
		dev.CreateBranchX("b1", "initial")
		dev.PushBranchToRemoteX("b1", config.OriginRemote)
		branches, err := origin.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
	})

	t.Run(".RemoveBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranchX("b1", "initial")
		branches, err := runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
		runtime.RemoveBranchX("b1")
		branches, err = runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial"}, branches)
	})

	t.Run(".RemoveRemote()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemoteX(config.OriginRemote, origin.WorkingDir)
		repo.RemoveRemoteX(config.OriginRemote)
		remotes, err := repo.Remotes()
		assert.NoError(t, err)
		assert.Len(t, remotes, 0)
	})

	t.Run(".ShaForCommit()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		repo.CreateCommitX(git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
		sha := repo.ShaForCommitX("commit")
		assert.Len(t, sha, 40)
	})

	t.Run(".UncommittedFiles()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFileX("f1.txt", "one")
		runtime.CreateFileX("f2.txt", "two")
		files := runtime.UncommittedFilesX()
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
