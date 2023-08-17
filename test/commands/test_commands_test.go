package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
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
		assert.Equal(t, config.Remotes{}, remotes)
		origin := testruntime.Create(t)
		dev.AddRemote(config.OriginRemote, origin.WorkingDir)
		remotes, err = dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, config.Remotes{"origin"}, remotes)
	})

	t.Run(".Commits()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(git.Commit{
			Branch:      domain.NewLocalBranchName("initial"),
			FileName:    "file1",
			FileContent: "hello",
			Message:     "first commit",
		})
		runtime.CreateCommit(git.Commit{
			Branch:      domain.NewLocalBranchName("initial"),
			FileName:    "file2",
			FileContent: "hello again",
			Message:     "second commit",
		})
		commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
		assert.Len(t, commits, 2)
		assert.Equal(t, domain.NewLocalBranchName("initial"), commits[0].Branch)
		assert.Equal(t, "file1", commits[0].FileName)
		assert.Equal(t, "hello", commits[0].FileContent)
		assert.Equal(t, "first commit", commits[0].Message)
		assert.Equal(t, domain.NewLocalBranchName("initial"), commits[1].Branch)
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
		runtime.AddRemote(config.OriginRemote, origin.WorkingDir)
		runtime.Fetch()
		runtime.ConnectTrackingBranch("initial")
		runtime.PushBranch()
	})

	t.Run(".CreateBranch()", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("initial"))
			currentBranch, err := runtime.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, domain.NewLocalBranchName("initial"), currentBranch)
			branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
			assert.NoError(t, err)
			want := domain.LocalBranchNamesFrom("initial", "branch1")
			assert.Equal(t, want, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch(domain.NewLocalBranchName("my/feature"), domain.NewLocalBranchName("initial"))
			currentBranch, err := runtime.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, domain.NewLocalBranchName("initial"), currentBranch)
			branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
			assert.NoError(t, err)
			want := domain.LocalBranchNamesFrom("initial", "my/feature")
			assert.Equal(t, want, branches)
		})
	})

	t.Run(".CreateChildFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		err := runtime.CreateFeatureBranch(domain.NewLocalBranchName("f1"))
		assert.NoError(t, err)
		runtime.CreateChildFeatureBranch(domain.NewLocalBranchName("f1a"), domain.NewLocalBranchName("f1"))
		output, err := runtime.BackendRunner.QueryTrim("git-town", "config")
		assert.NoError(t, err)
		output = stripansi.Strip(output)
		if !strings.Contains(output, "Branch Lineage:\n  main\n    f1\n      f1a") {
			t.Fatalf("unexpected output:\n%s", output)
		}
	})

	t.Run(".CreateCommit()", func(t *testing.T) {
		t.Run("minimal arguments", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(git.Commit{
				Branch:      domain.NewLocalBranchName("initial"),
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, domain.NewLocalBranchName("initial"), commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(git.Commit{
				Branch:      domain.NewLocalBranchName("initial"),
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
				Author:      "developer <developer@example.com>",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, domain.NewLocalBranchName("initial"), commits[0].Branch)
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
		runtime.CreatePerennialBranches(domain.NewLocalBranchName("p1"), domain.NewLocalBranchName("p2"))
		branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("main"))
		assert.NoError(t, err)
		want := domain.LocalBranchNamesFrom("main", "initial", "p1", "p2")
		assert.Equal(t, want, branches)
		runtime.Config.Reload()
		durations := runtime.Config.BranchDurations()
		assert.True(t, durations.IsPerennialBranch(domain.NewLocalBranchName("p1")))
		assert.True(t, durations.IsPerennialBranch(domain.NewLocalBranchName("p2")))
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
			Branch:      domain.NewLocalBranchName("initial"),
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		commits := runtime.CommitsInBranch(domain.NewLocalBranchName("initial"), []string{})
		assert.Len(t, commits, 1)
		content := runtime.FileContentInCommit(commits[0].SHA.Location, "hello.txt")
		assert.Equal(t, "hello world", content)
	})

	t.Run(".FilesInCommit()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		runtime.StageFiles("f1.txt", "f2.txt")
		runtime.CommitStagedChanges("stuff")
		commits := runtime.Commits([]string{}, "initial")
		assert.Len(t, commits, 1)
		fileNames := runtime.FilesInCommit(commits[0].SHA)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run(".HasBranchesOutOfSync()", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			runner := env.DevRepo
			runner.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("main"))
			runner.CheckoutBranch(domain.NewLocalBranchName("branch1"))
			runner.CreateFile("file1", "content")
			runner.StageFiles("file1")
			runner.CommitStagedChanges("stuff")
			runner.PushBranchToRemote(domain.NewLocalBranchName("branch1"), config.OriginRemote)
			have := runner.HasBranchesOutOfSync()
			assert.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			env.DevRepo.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("main"))
			env.DevRepo.PushBranch()
			env.DevRepo.CreateFile("file1", "content")
			env.DevRepo.StageFiles("file1")
			env.DevRepo.CommitStagedChanges("stuff")
			have := env.DevRepo.HasBranchesOutOfSync()
			assert.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			env.DevRepo.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("main"))
			env.DevRepo.PushBranch()
			env.OriginRepo.CheckoutBranch(domain.NewLocalBranchName("main"))
			env.OriginRepo.CreateFile("file1", "content")
			env.OriginRepo.StageFiles("file1")
			env.OriginRepo.CommitStagedChanges("stuff")
			env.OriginRepo.CheckoutBranch(domain.NewLocalBranchName("initial"))
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
		runtime.CreateBranch(domain.NewLocalBranchName("main"), domain.NewLocalBranchName("initial"))
		err := runtime.CreateFeatureBranch(domain.NewLocalBranchName("foo"))
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
		dev.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		dev.PushBranchToRemote(domain.NewLocalBranchName("b1"), config.OriginRemote)
		branches, err := origin.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		assert.NoError(t, err)
		want := domain.LocalBranchNamesFrom("initial", "b1")
		assert.Equal(t, want, branches)
	})

	t.Run(".RemoveBranch()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		assert.NoError(t, err)
		want := domain.LocalBranchNamesFrom("initial", "b1")
		assert.Equal(t, want, branches)
		runtime.RemoveBranch("b1")
		branches, err = runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		assert.NoError(t, err)
		wantBranches := domain.LocalBranchNamesFrom("initial")
		assert.Equal(t, wantBranches, branches)
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
		repo.CreateCommit(git.Commit{Branch: domain.NewLocalBranchName("initial"), FileName: "foo", FileContent: "bar", Message: "commit"})
		sha := repo.ShaForCommit("commit")
		assert.Len(t, sha, 40)
	})

	t.Run(".UncommittedFiles()", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		files := runtime.UncommittedFiles()
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
