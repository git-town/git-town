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
		dev.AddRemote(config.OriginRemote, origin.WorkingDir)
		remotes, err = dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})

	t.Run(".Commits()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file1",
			FileContent: "hello",
			Message:     "first commit",
		})
		runtime.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file2",
			FileContent: "hello again",
			Message:     "second commit",
		})
		commits, err := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
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
		origin := testruntime.Create(t)
		repoDir := filepath.Join(t.TempDir(), "repo") // need a non-existing directory
		err := helpers.CopyDirectory(origin.WorkingDir, repoDir)
		assert.NoError(t, err)
		runtime := testruntime.New(repoDir, repoDir, "")
		runtime.AddRemote(config.OriginRemote, origin.WorkingDir)
		runtime.Fetch()
		runtime.ConnectTrackingBranch("initial")
		runtime.PushBranch()
	})

	t.Run(".CreateBranch()", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("branch1", "initial")
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
			runtime.CreateBranch("my/feature", "initial")
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
		runtime.CreateChildFeatureBranch("f1a", "f1")
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
			runtime.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			commits, err := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.NoError(t, err)
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, "initial", commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
				Author:      "developer <developer@example.com>",
			})
			commits, err := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.NoError(t, err)
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
			runtime.CreateFile("filename", "content")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})

		t.Run("create file in subfolder", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("folder/filename", "content")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "folder/filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})
	})

	t.Run(".CreatePerennialBranches()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreatePerennialBranches("p1", "p2")
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
		repo.AddRemote(config.OriginRemote, origin.WorkingDir)
		repo.Fetch()
	})

	t.Run(".FileContentInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		commits, err := runtime.CommitsInBranch("initial", []string{})
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		content := runtime.FileContentInCommit(commits[0].SHA, "hello.txt")
		assert.Equal(t, "hello world", content)
	})

	t.Run(".FilesInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		err := runtime.StageFiles("f1.txt", "f2.txt")
		assert.NoError(t, err)
		runtime.CommitStagedChanges("stuff")
		commits, err := runtime.Commits([]string{}, "initial")
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		fileNames := runtime.FilesInCommit(commits[0].SHA)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run(".HasBranchesOutOfSync()", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env, err := fixture.NewStandardFixture(t.TempDir())
			assert.NoError(t, err)
			runner := env.DevRepo
			runner.CreateBranch("branch1", "main")
			runner.CheckoutBranch("branch1")
			runner.CreateFile("file1", "content")
			err = runner.StageFiles("file1")
			assert.NoError(t, err)
			runner.CommitStagedChanges("stuff")
			runner.PushBranchToRemote("branch1", config.OriginRemote)
			have := runner.HasBranchesOutOfSync()
			assert.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			env, err := fixture.NewStandardFixture(t.TempDir())
			assert.NoError(t, err)
			env.DevRepo.CreateBranch("branch1", "main")
			env.DevRepo.PushBranch()
			env.DevRepo.CreateFile("file1", "content")
			err = env.DevRepo.StageFiles("file1")
			assert.NoError(t, err)
			env.DevRepo.CommitStagedChanges("stuff")
			have := env.DevRepo.HasBranchesOutOfSync()
			assert.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			env, err := fixture.NewStandardFixture(t.TempDir())
			assert.NoError(t, err)
			env.DevRepo.CreateBranch("branch1", "main")
			env.DevRepo.PushBranch()
			env.OriginRepo.CheckoutBranch("main")
			env.OriginRepo.CreateFile("file1", "content")
			err = env.OriginRepo.StageFiles("file1")
			assert.NoError(t, err)
			env.OriginRepo.CommitStagedChanges("stuff")
			env.OriginRepo.CheckoutBranch("initial")
			env.DevRepo.Fetch()
			have := env.DevRepo.HasBranchesOutOfSync()
			assert.True(t, have)
		})
	})

	t.Run(".HasFile()", func(t *testing.T) {
		t.Parallel()
		t.Run("filename and content match", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("f1.txt", "one")
			assert.Equal(t, "", runtime.HasFile("f1.txt", "one"))
		})
		t.Run("filename matches, content doesn't match", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("f1.txt", "one")
			assert.Equal(t, "file \"f1.txt\" should have content \"zonk\" but has \"one\"", runtime.HasFile("f1.txt", "zonk"))
		})
		t.Run("filename doesn't match but content matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("f1.txt", "one")
			assert.Equal(t, "repo doesn't have file \"zonk.txt\"", runtime.HasFile("zonk.txt", "one"))
		})
	})

	t.Run(".HasGitTownConfigNow()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		res := runtime.HasGitTownConfigNow()
		assert.False(t, res)
		runtime.CreateBranch("main", "initial")
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
		dev.AddRemote(config.OriginRemote, origin.WorkingDir)
		dev.CreateBranch("b1", "initial")
		dev.PushBranchToRemote("b1", config.OriginRemote)
		branches, err := origin.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
	})

	t.Run(".RemoveBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("b1", "initial")
		branches, err := runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
		runtime.RemoveBranch("b1")
		branches, err = runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial"}, branches)
	})

	t.Run(".RemoveRemote()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(config.OriginRemote, origin.WorkingDir)
		repo.RemoveRemote(config.OriginRemote)
		remotes, err := repo.Remotes()
		assert.NoError(t, err)
		assert.Len(t, remotes, 0)
	})

	t.Run(".ShaForCommit()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		repo.CreateCommit(git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
		sha, err := repo.ShaForCommit("commit")
		assert.NoError(t, err)
		assert.Len(t, sha, 40)
	})

	t.Run(".UncommittedFiles()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		files, err := runtime.UncommittedFiles()
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
