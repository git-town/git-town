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
		err = dev.AddRemote(config.OriginRemote, origin.WorkingDir)
		assert.NoError(t, err)
		remotes, err = dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})

	t.Run(".Commits()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file1",
			FileContent: "hello",
			Message:     "first commit",
		})
		assert.NoError(t, err)
		err = runtime.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file2",
			FileContent: "hello again",
			Message:     "second commit",
		})
		assert.NoError(t, err)
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
		err = runtime.AddRemote(config.OriginRemote, origin.WorkingDir)
		assert.NoError(t, err)
		err = runtime.Fetch()
		assert.NoError(t, err)
		err = runtime.ConnectTrackingBranch("initial")
		assert.NoError(t, err)
		err = runtime.PushBranch()
		assert.NoError(t, err)
	})

	t.Run(".CreateBranch()", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			err := runtime.CreateBranch("branch1", "initial")
			assert.NoError(t, err)
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
			err := runtime.CreateBranch("my/feature", "initial")
			assert.NoError(t, err)
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
		err = runtime.CreateChildFeatureBranch("f1a", "f1")
		assert.NoError(t, err)
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
			err := runtime.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			assert.NoError(t, err)
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
			err := runtime.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
				Author:      "developer <developer@example.com>",
			})
			assert.NoError(t, err)
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
			err := runtime.CreateFile("filename", "content")
			assert.Nil(t, err, "cannot create file in repo")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})

		t.Run("create file in subfolder", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			err := runtime.CreateFile("folder/filename", "content")
			assert.Nil(t, err, "cannot create file in repo")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "folder/filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})
	})

	t.Run(".CreatePerennialBranches()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		err := runtime.CreatePerennialBranches("p1", "p2")
		assert.NoError(t, err)
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
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir)
		assert.NoError(t, err)
		err = repo.Fetch()
		assert.NoError(t, err)
	})

	t.Run(".FileContentInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		assert.NoError(t, err)
		commits, err := runtime.CommitsInBranch("initial", []string{})
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		content, err := runtime.FileContentInCommit(commits[0].SHA, "hello.txt")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", content)
	})

	t.Run(".FilesInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		err = runtime.CreateFile("f2.txt", "two")
		assert.NoError(t, err)
		err = runtime.StageFiles("f1.txt", "f2.txt")
		assert.NoError(t, err)
		err = runtime.CommitStagedChanges("stuff")
		assert.NoError(t, err)
		commits, err := runtime.Commits([]string{}, "initial")
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		fileNames, err := runtime.FilesInCommit(commits[0].SHA)
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run(".HasBranchesOutOfSync()", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env, err := fixture.NewStandardFixture(t.TempDir())
			assert.NoError(t, err)
			runner := env.DevRepo
			err = runner.CreateBranch("branch1", "main")
			assert.NoError(t, err)
			err = runner.CheckoutBranch("branch1")
			assert.NoError(t, err)
			err = runner.CreateFile("file1", "content")
			assert.NoError(t, err)
			err = runner.StageFiles("file1")
			assert.NoError(t, err)
			err = runner.CommitStagedChanges("stuff")
			assert.NoError(t, err)
			err = runner.PushBranchToRemote("branch1", config.OriginRemote)
			assert.NoError(t, err)
			have, err := runner.HasBranchesOutOfSync()
			assert.NoError(t, err)
			assert.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			env, err := fixture.NewStandardFixture(t.TempDir())
			assert.NoError(t, err)
			err = env.DevRepo.CreateBranch("branch1", "main")
			assert.NoError(t, err)
			err = env.DevRepo.PushBranch()
			assert.NoError(t, err)
			err = env.DevRepo.CreateFile("file1", "content")
			assert.NoError(t, err)
			err = env.DevRepo.StageFiles("file1")
			assert.NoError(t, err)
			err = env.DevRepo.CommitStagedChanges("stuff")
			assert.NoError(t, err)
			have, err := env.DevRepo.HasBranchesOutOfSync()
			assert.NoError(t, err)
			assert.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			env, err := fixture.NewStandardFixture(t.TempDir())
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
			have, err := env.DevRepo.HasBranchesOutOfSync()
			assert.NoError(t, err)
			assert.True(t, have)
		})
	})

	t.Run(".HasFile()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		has, err := runtime.HasFile("f1.txt", "one")
		assert.NoError(t, err)
		assert.True(t, has)
		_, err = runtime.HasFile("f1.txt", "zonk")
		assert.Error(t, err)
		_, err = runtime.HasFile("zonk.txt", "one")
		assert.Error(t, err)
	})

	t.Run(".HasGitTownConfigNow()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		res := runtime.HasGitTownConfigNow()
		assert.False(t, res)
		err := runtime.CreateBranch("main", "initial")
		assert.NoError(t, err)
		err = runtime.CreateFeatureBranch("foo")
		assert.NoError(t, err)
		res = runtime.HasGitTownConfigNow()
		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run(".PushBranchToRemote()", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		origin := testruntime.Create(t)
		err := dev.AddRemote(config.OriginRemote, origin.WorkingDir)
		assert.NoError(t, err)
		err = dev.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = dev.PushBranchToRemote("b1", config.OriginRemote)
		assert.NoError(t, err)
		branches, err := origin.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
	})

	t.Run(".RemoveBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		branches, err := runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
		err = runtime.RemoveBranch("b1")
		assert.NoError(t, err)
		branches, err = runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial"}, branches)
	})

	t.Run(".RemoveRemote()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		err := repo.AddRemote(config.OriginRemote, origin.WorkingDir)
		assert.NoError(t, err)
		err = repo.RemoveRemote(config.OriginRemote)
		assert.NoError(t, err)
		remotes, err := repo.Remotes()
		assert.NoError(t, err)
		assert.Len(t, remotes, 0)
	})

	t.Run(".ShaForCommit()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		err := repo.CreateCommit(git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
		assert.NoError(t, err)
		sha, err := repo.ShaForCommit("commit")
		assert.NoError(t, err)
		assert.Len(t, sha, 40)
	})

	t.Run(".UncommittedFiles()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		err := runtime.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		err = runtime.CreateFile("f2.txt", "two")
		assert.NoError(t, err)
		files, err := runtime.UncommittedFiles()
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
