package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/test/filesystem"
	"github.com/git-town/git-town/v11/test/fixture"
	"github.com/git-town/git-town/v11/test/git"
	"github.com/git-town/git-town/v11/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestTestCommands(t *testing.T) {
	t.Parallel()

	t.Run("AddRemote", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		remotes, err := dev.Remotes()
		must.NoError(t, err)
		must.Eq(t, domain.Remotes{}, remotes)
		origin := testruntime.Create(t)
		dev.AddRemote(domain.OriginRemote, origin.WorkingDir)
		remotes, err = dev.Remotes()
		must.NoError(t, err)
		must.Eq(t, domain.Remotes{domain.OriginRemote}, remotes)
	})

	t.Run("Commits", func(t *testing.T) {
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
		commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, domain.NewLocalBranchName("initial"))
		must.Len(t, 2, commits)
		must.EqOp(t, domain.NewLocalBranchName("initial"), commits[0].Branch)
		must.EqOp(t, "file1", commits[0].FileName)
		must.EqOp(t, "hello", commits[0].FileContent)
		must.EqOp(t, "first commit", commits[0].Message)
		must.EqOp(t, domain.NewLocalBranchName("initial"), commits[1].Branch)
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
		runtime.AddRemote(domain.OriginRemote, origin.WorkingDir)
		runtime.Fetch()
		runtime.ConnectTrackingBranch(domain.NewLocalBranchName("initial"))
		runtime.PushBranch()
	})

	t.Run("CreateBranch", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("initial"))
			currentBranch, err := runtime.CurrentBranch()
			must.NoError(t, err)
			must.EqOp(t, domain.NewLocalBranchName("initial"), currentBranch)
			branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
			must.NoError(t, err)
			want := domain.NewLocalBranchNames("initial", "branch1")
			must.Eq(t, want, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch(domain.NewLocalBranchName("my/feature"), domain.NewLocalBranchName("initial"))
			currentBranch, err := runtime.CurrentBranch()
			must.NoError(t, err)
			must.EqOp(t, domain.NewLocalBranchName("initial"), currentBranch)
			branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
			must.NoError(t, err)
			want := domain.NewLocalBranchNames("initial", "my/feature")
			must.Eq(t, want, branches)
		})
	})

	t.Run("CreateChildFeatureBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		err := runtime.CreateFeatureBranch(domain.NewLocalBranchName("f1"))
		must.NoError(t, err)
		runtime.CreateChildFeatureBranch(domain.NewLocalBranchName("f1a"), domain.NewLocalBranchName("f1"))
		output, err := runtime.BackendRunner.QueryTrim("git-town", "config")
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
				Branch:      domain.NewLocalBranchName("initial"),
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, domain.NewLocalBranchName("initial"))
			must.Len(t, 1, commits)
			must.EqOp(t, "hello.txt", commits[0].FileName)
			must.EqOp(t, "hello world", commits[0].FileContent)
			must.EqOp(t, "test commit", commits[0].Message)
			must.EqOp(t, domain.NewLocalBranchName("initial"), commits[0].Branch)
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
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, domain.NewLocalBranchName("initial"))
			must.Len(t, 1, commits)
			must.EqOp(t, "hello.txt", commits[0].FileName)
			must.EqOp(t, "hello world", commits[0].FileContent)
			must.EqOp(t, "test commit", commits[0].Message)
			must.EqOp(t, domain.NewLocalBranchName("initial"), commits[0].Branch)
			must.EqOp(t, "developer <developer@example.com>", commits[0].Author)
		})
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
		runtime.CreatePerennialBranches(domain.NewLocalBranchName("p1"), domain.NewLocalBranchName("p2"))
		branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("main"))
		must.NoError(t, err)
		want := domain.NewLocalBranchNames("main", "initial", "p1", "p2")
		must.Eq(t, want, branches)
		runtime.GitTown.Reload()
		branchTypes := runtime.GitTown.BranchTypes()
		must.True(t, branchTypes.IsPerennialBranch(domain.NewLocalBranchName("p1")))
		must.True(t, branchTypes.IsPerennialBranch(domain.NewLocalBranchName("p2")))
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(domain.OriginRemote, origin.WorkingDir)
		repo.Fetch()
	})

	t.Run("FileContentInCommit", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(git.Commit{
			Branch:      domain.NewLocalBranchName("initial"),
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		commits := runtime.CommitsInBranch(domain.NewLocalBranchName("initial"), []string{})
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
		commits := runtime.Commits([]string{}, domain.NewLocalBranchName("initial"))
		must.Len(t, 1, commits)
		fileNames := runtime.FilesInCommit(commits[0].SHA)
		must.Eq(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run("HasBranchesOutOfSync", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env := fixture.NewStandardFixture(t.TempDir())
			runner := env.DevRepo
			runner.CreateBranch(domain.NewLocalBranchName("branch1"), domain.NewLocalBranchName("main"))
			runner.CheckoutBranch(domain.NewLocalBranchName("branch1"))
			runner.CreateFile("file1", "content")
			runner.StageFiles("file1")
			runner.CommitStagedChanges("stuff")
			runner.PushBranchToRemote(domain.NewLocalBranchName("branch1"), domain.OriginRemote)
			have := runner.HasBranchesOutOfSync()
			must.False(t, have)
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
			must.True(t, have)
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
			res := runtime.HasGitTownConfigNow()
			must.False(t, res)
		})
		t.Run("the main branch is configured", func(t *testing.T) {
			runtime := testruntime.Create(t)
			must.NoError(t, runtime.GitTown.SetMainBranch(domain.NewLocalBranchName("main")))
			must.True(t, runtime.HasGitTownConfigNow())
		})
		t.Run("the perennial branches are configured", func(t *testing.T) {
			runtime := testruntime.Create(t)
			must.NoError(t, runtime.GitTown.SetPerennialBranches(domain.NewLocalBranchNames("qa")))
			must.True(t, runtime.HasGitTownConfigNow())
		})
		t.Run("branch lineage is configured", func(t *testing.T) {
			runtime := testruntime.Create(t)
			runtime.CreateBranch(domain.NewLocalBranchName("main"), domain.NewLocalBranchName("initial"))
			must.NoError(t, runtime.CreateFeatureBranch(domain.NewLocalBranchName("foo")))
			must.True(t, runtime.HasGitTownConfigNow())
		})
	})

	t.Run("LocalBranchesMainFirst", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		initial := domain.NewLocalBranchName("initial")
		runner.CreateBranch(domain.NewLocalBranchName("b1"), initial)
		runner.CreateBranch(domain.NewLocalBranchName("b2"), initial)
		origin.CreateBranch(domain.NewLocalBranchName("b3"), initial)
		runner.Fetch()
		branches, err := runner.LocalBranchesMainFirst(initial)
		must.NoError(t, err)
		want := domain.NewLocalBranchNames("initial", "b1", "b2")
		must.Eq(t, want, branches)
	})

	t.Run("PushBranchToRemote", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		origin := testruntime.Create(t)
		dev.AddRemote(domain.OriginRemote, origin.WorkingDir)
		dev.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		dev.PushBranchToRemote(domain.NewLocalBranchName("b1"), domain.OriginRemote)
		branches, err := origin.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		must.NoError(t, err)
		want := domain.NewLocalBranchNames("initial", "b1")
		must.Eq(t, want, branches)
	})

	t.Run("RemoveBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch(domain.NewLocalBranchName("b1"), domain.NewLocalBranchName("initial"))
		branches, err := runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		must.NoError(t, err)
		want := domain.NewLocalBranchNames("initial", "b1")
		must.Eq(t, want, branches)
		runtime.RemoveBranch(domain.NewLocalBranchName("b1"))
		branches, err = runtime.LocalBranchesMainFirst(domain.NewLocalBranchName("initial"))
		must.NoError(t, err)
		wantBranches := domain.NewLocalBranchNames("initial")
		must.Eq(t, wantBranches, branches)
	})

	t.Run("RemoveRemote", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(domain.OriginRemote, origin.WorkingDir)
		repo.RemoveRemote(domain.OriginRemote)
		remotes, err := repo.Remotes()
		must.NoError(t, err)
		must.Len(t, 0, remotes)
	})

	t.Run("SHAForCommit", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		repo.CreateCommit(git.Commit{Branch: domain.NewLocalBranchName("initial"), FileName: "foo", FileContent: "bar", Message: "commit"})
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
