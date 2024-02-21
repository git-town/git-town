package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/test/filesystem"
	"github.com/git-town/git-town/v12/test/fixture"
	"github.com/git-town/git-town/v12/test/git"
	"github.com/git-town/git-town/v12/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestTestCommands(t *testing.T) {
	t.Parallel()

	t.Run("AddRemote", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		remotes, err := dev.Remotes()
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{}, remotes)
		origin := testruntime.Create(t)
		dev.AddRemote(gitdomain.OriginRemote, origin.WorkingDir)
		remotes, err = dev.Remotes()
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{gitdomain.OriginRemote}, remotes)
	})

	t.Run("Commits", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(git.Commit{
			Branch:      gitdomain.NewLocalBranchName("initial"),
			FileContent: "hello",
			FileName:    "file1",
			Message:     "first commit",
		})
		runtime.CreateCommit(git.Commit{
			Branch:      gitdomain.NewLocalBranchName("initial"),
			FileContent: "hello again",
			FileName:    "file2",
			Message:     "second commit",
		})
		commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, gitdomain.NewLocalBranchName("initial"))
		must.Len(t, 2, commits)
		must.EqOp(t, gitdomain.NewLocalBranchName("initial"), commits[0].Branch)
		must.EqOp(t, "file1", commits[0].FileName)
		must.EqOp(t, "hello", commits[0].FileContent)
		must.EqOp(t, "first commit", commits[0].Message)
		must.EqOp(t, gitdomain.NewLocalBranchName("initial"), commits[1].Branch)
		must.EqOp(t, "file2", commits[1].FileName)
		must.EqOp(t, "hello again", commits[1].FileContent)
		must.EqOp(t, "second commit", commits[1].Message)
	})

	t.Run("ConnectTrackingBranch", func(t *testing.T) {
		t.Parallel()
		// replicating the situation this is used in,
		// connecting branches of repos with the same commits in them
		origin := testruntime.Create(t)
		repoDir := filepath.Join(t.TempDir(), "repo") // need a non-existing directory
		filesystem.CopyDirectory(origin.WorkingDir, repoDir)
		runtime := testruntime.New(repoDir, repoDir, "")
		runtime.AddRemote(gitdomain.OriginRemote, origin.WorkingDir)
		runtime.Fetch()
		runtime.ConnectTrackingBranch(gitdomain.NewLocalBranchName("initial"))
		runtime.PushBranch()
	})

	t.Run("CreateBranch", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch(gitdomain.NewLocalBranchName("branch1"), gitdomain.NewLocalBranchName("initial"))
			currentBranch, err := runtime.CurrentBranch()
			must.NoError(t, err)
			must.EqOp(t, gitdomain.NewLocalBranchName("initial"), currentBranch)
			branches, err := runtime.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("initial"))
			must.NoError(t, err)
			want := gitdomain.NewLocalBranchNames("initial", "branch1")
			must.Eq(t, want, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch(gitdomain.NewLocalBranchName("my/feature"), gitdomain.NewLocalBranchName("initial"))
			currentBranch, err := runtime.CurrentBranch()
			must.NoError(t, err)
			must.EqOp(t, gitdomain.NewLocalBranchName("initial"), currentBranch)
			branches, err := runtime.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("initial"))
			must.NoError(t, err)
			want := gitdomain.NewLocalBranchNames("initial", "my/feature")
			must.Eq(t, want, branches)
		})
	})

	t.Run("CreateChildFeatureBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch(gitdomain.NewLocalBranchName("f1"))
		runtime.CreateChildFeatureBranch(gitdomain.NewLocalBranchName("f1a"), gitdomain.NewLocalBranchName("f1"))
		output, err := runtime.Runner.QueryTrim("git-town", "config")
		must.NoError(t, err)
		output = stripansi.Strip(output)
		if !strings.Contains(output, "Branch Lineage:\n  main\n    f1\n      f1a") {
			t.Fatalf("unexpected output:\n%s", output)
		}
	})

	t.Run("CreateCommit", func(t *testing.T) {
		t.Run("minimal arguments", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(git.Commit{
				Branch:      gitdomain.NewLocalBranchName("initial"),
				FileContent: "hello world",
				FileName:    "hello.txt",
				Message:     "test commit",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, gitdomain.NewLocalBranchName("initial"))
			must.Len(t, 1, commits)
			must.EqOp(t, "hello.txt", commits[0].FileName)
			must.EqOp(t, "hello world", commits[0].FileContent)
			must.EqOp(t, "test commit", commits[0].Message)
			must.EqOp(t, gitdomain.NewLocalBranchName("initial"), commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(git.Commit{
				Author:      "developer <developer@example.com>",
				Branch:      gitdomain.NewLocalBranchName("initial"),
				FileContent: "hello world",
				FileName:    "hello.txt",
				Message:     "test commit",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, gitdomain.NewLocalBranchName("initial"))
			must.Len(t, 1, commits)
			must.EqOp(t, "hello.txt", commits[0].FileName)
			must.EqOp(t, "hello world", commits[0].FileContent)
			must.EqOp(t, "test commit", commits[0].Message)
			must.EqOp(t, gitdomain.NewLocalBranchName("initial"), commits[0].Branch)
			must.EqOp(t, "developer <developer@example.com>", commits[0].Author)
		})
	})

	t.Run("CreateFeatureBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch(gitdomain.NewLocalBranchName("f1"))
		runtime.Config.Reload()
		must.True(t, runtime.Config.FullConfig.IsFeatureBranch(gitdomain.NewLocalBranchName("f1")))
		lineageHave := runtime.Config.FullConfig.Lineage
		lineageWant := configdomain.Lineage{}
		lineageWant[gitdomain.NewLocalBranchName("f1")] = gitdomain.NewLocalBranchName("main")
		must.Eq(t, lineageWant, lineageHave)
	})

	t.Run("CreateFile", func(t *testing.T) {
		t.Run("simple example", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("filename", "content")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "filename"))
			must.Nil(t, err)
			must.EqOp(t, "content", string(content))
		})

		t.Run("create file in subfolder", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("folder/filename", "content")
			content, err := os.ReadFile(filepath.Join(runtime.WorkingDir, "folder/filename"))
			must.Nil(t, err)
			must.EqOp(t, "content", string(content))
		})
	})

	t.Run("CreatePerennialBranches", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreatePerennialBranches(gitdomain.NewLocalBranchName("p1"), gitdomain.NewLocalBranchName("p2"))
		branches, err := runtime.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("main"))
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("main", "initial", "p1", "p2")
		must.Eq(t, want, branches)
		runtime.Config.Reload()
		must.True(t, runtime.Config.FullConfig.IsPerennialBranch(gitdomain.NewLocalBranchName("p1")))
		must.True(t, runtime.Config.FullConfig.IsPerennialBranch(gitdomain.NewLocalBranchName("p2")))
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(gitdomain.OriginRemote, origin.WorkingDir)
		repo.Fetch()
	})

	t.Run("FileContentInCommit", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(git.Commit{
			Branch:      gitdomain.NewLocalBranchName("initial"),
			FileContent: "hello world",
			FileName:    "hello.txt",
			Message:     "commit",
		})
		commits := runtime.CommitsInBranch(gitdomain.NewLocalBranchName("initial"), []string{})
		must.Len(t, 1, commits)
		content := runtime.FileContentInCommit(commits[0].SHA.Location(), "hello.txt")
		must.EqOp(t, "hello world", content)
	})

	t.Run("FilesInCommit", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		runtime.StageFiles("f1.txt", "f2.txt")
		runtime.CommitStagedChanges("stuff")
		commits := runtime.Commits([]string{}, gitdomain.NewLocalBranchName("initial"))
		must.Len(t, 1, commits)
		fileNames := runtime.FilesInCommit(commits[0].SHA)
		must.Eq(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run("HasBranchesOutOfSync", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			runner := env.DevRepo
			runner.CreateBranch(gitdomain.NewLocalBranchName("branch1"), gitdomain.NewLocalBranchName("main"))
			runner.CheckoutBranch(gitdomain.NewLocalBranchName("branch1"))
			runner.CreateFile("file1", "content")
			runner.StageFiles("file1")
			runner.CommitStagedChanges("stuff")
			runner.PushBranchToRemote(gitdomain.NewLocalBranchName("branch1"), gitdomain.OriginRemote)
			have := runner.HasBranchesOutOfSync()
			must.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			env.DevRepo.CreateBranch(gitdomain.NewLocalBranchName("branch1"), gitdomain.NewLocalBranchName("main"))
			env.DevRepo.PushBranch()
			env.DevRepo.CreateFile("file1", "content")
			env.DevRepo.StageFiles("file1")
			env.DevRepo.CommitStagedChanges("stuff")
			have := env.DevRepo.HasBranchesOutOfSync()
			must.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			env.DevRepo.CreateBranch(gitdomain.NewLocalBranchName("branch1"), gitdomain.NewLocalBranchName("main"))
			env.DevRepo.PushBranch()
			env.OriginRepo.CheckoutBranch(gitdomain.NewLocalBranchName("main"))
			env.OriginRepo.CreateFile("file1", "content")
			env.OriginRepo.StageFiles("file1")
			env.OriginRepo.CommitStagedChanges("stuff")
			env.OriginRepo.CheckoutBranch(gitdomain.NewLocalBranchName("initial"))
			env.DevRepo.Fetch()
			have := env.DevRepo.HasBranchesOutOfSync()
			must.True(t, have)
		})
	})

	t.Run("HasFile", func(t *testing.T) {
		t.Parallel()
		t.Run("filename and content match", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("f1.txt", "one")
			must.EqOp(t, "", runtime.HasFile("f1.txt", "one"))
		})
		t.Run("filename matches, content doesn't match", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("f1.txt", "one")
			must.EqOp(t, "file \"f1.txt\" should have content \"zonk\" but has \"one\"", runtime.HasFile("f1.txt", "zonk"))
		})
		t.Run("filename doesn't match but content matches", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("f1.txt", "one")
			must.EqOp(t, "repo doesn't have file \"zonk.txt\"", runtime.HasFile("zonk.txt", "one"))
		})
	})

	t.Run("HasGitTownConfigNow", func(t *testing.T) {
		t.Parallel()
		t.Run("no config exists", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			must.NoError(t, runtime.VerifyNoGitTownConfiguration())
		})
		t.Run("the main branch is configured", func(t *testing.T) {
			runtime := testruntime.Create(t)
			must.NoError(t, runtime.Config.SetMainBranch(gitdomain.NewLocalBranchName("main")))
			must.Error(t, runtime.VerifyNoGitTownConfiguration())
		})
		t.Run("the perennial branches are configured", func(t *testing.T) {
			runtime := testruntime.Create(t)
			must.NoError(t, runtime.Config.SetPerennialBranches(gitdomain.NewLocalBranchNames("qa")))
			must.Error(t, runtime.VerifyNoGitTownConfiguration())
		})
		t.Run("branch lineage is configured", func(t *testing.T) {
			runtime := testruntime.Create(t)
			runtime.CreateBranch(gitdomain.NewLocalBranchName("main"), gitdomain.NewLocalBranchName("initial"))
			runtime.CreateFeatureBranch(gitdomain.NewLocalBranchName("foo"))
			must.Error(t, runtime.VerifyNoGitTownConfiguration())
		})
	})

	t.Run("LocalBranchesMainFirst", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		initial := gitdomain.NewLocalBranchName("initial")
		runner.CreateBranch(gitdomain.NewLocalBranchName("b1"), initial)
		runner.CreateBranch(gitdomain.NewLocalBranchName("b2"), initial)
		origin.CreateBranch(gitdomain.NewLocalBranchName("b3"), initial)
		runner.Fetch()
		branches, err := runner.LocalBranchesMainFirst(initial)
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "b1", "b2")
		must.Eq(t, want, branches)
	})

	t.Run("PushBranchToRemote", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		origin := testruntime.Create(t)
		dev.AddRemote(gitdomain.OriginRemote, origin.WorkingDir)
		dev.CreateBranch(gitdomain.NewLocalBranchName("b1"), gitdomain.NewLocalBranchName("initial"))
		dev.PushBranchToRemote(gitdomain.NewLocalBranchName("b1"), gitdomain.OriginRemote)
		branches, err := origin.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("initial"))
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "b1")
		must.Eq(t, want, branches)
	})

	t.Run("RemoveBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch(gitdomain.NewLocalBranchName("b1"), gitdomain.NewLocalBranchName("initial"))
		branches, err := runtime.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("initial"))
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "b1")
		must.Eq(t, want, branches)
		runtime.RemoveBranch(gitdomain.NewLocalBranchName("b1"))
		branches, err = runtime.LocalBranchesMainFirst(gitdomain.NewLocalBranchName("initial"))
		must.NoError(t, err)
		wantBranches := gitdomain.NewLocalBranchNames("initial")
		must.Eq(t, wantBranches, branches)
	})

	t.Run("RemoveRemote", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(gitdomain.OriginRemote, origin.WorkingDir)
		repo.RemoveRemote(gitdomain.OriginRemote)
		remotes, err := repo.Remotes()
		must.NoError(t, err)
		must.Len(t, 0, remotes)
	})

	t.Run("SHAForCommit", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		repo.CreateCommit(git.Commit{
			Branch:      gitdomain.NewLocalBranchName("initial"),
			FileContent: "bar",
			FileName:    "foo",
			Message:     "commit",
		})
		sha := repo.SHAForCommit("commit")
		must.EqOp(t, 7, len(sha))
	})

	t.Run("UncommittedFiles", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		files := runtime.UncommittedFiles()
		must.Eq(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
