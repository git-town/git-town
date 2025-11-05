package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/filesystem"
	"github.com/git-town/git-town/v22/internal/test/fixture"
	"github.com/git-town/git-town/v22/internal/test/testgit"
	"github.com/git-town/git-town/v22/internal/test/testruntime"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestTestCommands(t *testing.T) {
	t.Parallel()

	t.Run("AddRemote", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		remotes, err := dev.Git.Remotes(dev)
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{}, remotes)
		origin := testruntime.Create(t)
		dev.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		remotes, err = dev.Git.Remotes(dev)
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{gitdomain.RemoteOrigin}, remotes)
	})

	t.Run("CommitSHAs", func(t *testing.T) {
		t.Parallel()
		t.Run("includes commits with empty messages", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			// Create a commit with a normal message
			runtime.CreateCommit(testgit.Commit{
				Branch:      "initial",
				FileContent: "content1",
				FileName:    "file1",
				Message:     "first commit",
			})
			// Create a commit with an empty message
			runtime.CreateFile("file2", "content2")
			runtime.StageFiles("file2")
			runtime.MustRun("git", "commit", "--allow-empty-message", "-m", "")
			// Get all commits
			commits := runtime.CommitSHAs()
			must.EqOp(t, "", commits[0].Message)
			must.EqOp(t, "first commit", commits[1].Message)
			must.EqOp(t, "initial commit", commits[2].Message)
			must.EqOp(t, 40, len(commits[0].SHA.String()))
			must.EqOp(t, 40, len(commits[1].SHA.String()))
			must.EqOp(t, 40, len(commits[2].SHA.String()))
			must.NotEqOp(t, commits[0].SHA.String(), commits[1].SHA.String())
			must.NotEqOp(t, commits[0].SHA.String(), commits[2].SHA.String())
			must.NotEqOp(t, commits[1].SHA.String(), commits[2].SHA.String())
		})
	})

	t.Run("Commits", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "hello",
			FileName:    "file1",
			Message:     "first commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "hello again",
			FileName:    "file2",
			Message:     "second commit",
		})
		commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, runtime.Config.NormalConfig.Lineage, configdomain.OrderAsc)
		must.Len(t, 2, commits)
		must.EqOp(t, "initial", commits[0].Branch)
		must.EqOp(t, "file1", commits[0].FileName)
		must.EqOp(t, "hello", commits[0].FileContent)
		must.EqOp(t, "first commit", commits[0].Message)
		must.EqOp(t, "initial", commits[1].Branch)
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
		runtime.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		runtime.Fetch()
		runtime.ConnectTrackingBranch("initial")
		runtime.PushBranch()
	})

	t.Run("CreateBranch", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("branch1", "initial")
			currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
			must.EqOp(t, "initial", currentBranch)
			branches, _, err := runtime.LocalBranchesMainFirst("initial")
			must.NoError(t, err)
			want := gitdomain.NewLocalBranchNames("initial", "branch1")
			must.Eq(t, want, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("my/feature", "initial")
			currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
			must.EqOp(t, "initial", currentBranch)
			branches, _, err := runtime.LocalBranchesMainFirst("initial")
			must.NoError(t, err)
			want := gitdomain.NewLocalBranchNames("initial", "my/feature")
			must.Eq(t, want, branches)
		})
	})

	t.Run("CreateCommit", func(t *testing.T) {
		t.Run("minimal arguments", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(testgit.Commit{
				Branch:      "initial",
				FileContent: "hello world",
				FileName:    "hello.txt",
				Message:     "test commit",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, runtime.Config.NormalConfig.Lineage, configdomain.OrderAsc)
			must.Len(t, 1, commits)
			must.EqOp(t, "hello.txt", commits[0].FileName)
			must.EqOp(t, "hello world", commits[0].FileContent)
			must.EqOp(t, "test commit", commits[0].Message)
			must.EqOp(t, "initial", commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateCommit(testgit.Commit{
				Author:      "developer <developer@example.com>",
				Branch:      "initial",
				FileContent: "hello world",
				FileName:    "hello.txt",
				Message:     "test commit",
			})
			commits := runtime.Commits([]string{"FILE NAME", "FILE CONTENT"}, runtime.Config.NormalConfig.Lineage, configdomain.OrderAsc)
			must.Len(t, 1, commits)
			must.EqOp(t, "hello.txt", commits[0].FileName)
			must.EqOp(t, "hello world", commits[0].FileContent)
			must.EqOp(t, "test commit", commits[0].Message)
			must.EqOp(t, "initial", commits[0].Branch)
			must.EqOp(t, "developer <developer@example.com>", commits[0].Author)
		})
	})

	t.Run("CreateFeatureBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("f1", "main")
		runtime.Config.Reload(runtime.TestRunner)
		must.False(t, runtime.Config.IsMainOrPerennialBranch("f1"))
		lineageHave := runtime.Config.NormalConfig.Lineage
		lineageWant := configdomain.NewLineageWith(configdomain.LineageData{
			"f1": "main",
		})
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

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		repo.Fetch()
	})

	t.Run("FileContentInCommit", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "hello world",
			FileName:    "hello.txt",
			Message:     "commit",
		})
		commits := runtime.CommitsInBranch("initial", None[gitdomain.BranchName](), []string{})
		must.Len(t, 1, commits)
		content, deleted := runtime.FileContentInCommit(commits[0].SHA.Location(), "hello.txt")
		must.False(t, deleted)
		must.EqOp(t, "hello world", content)
	})

	t.Run("FilesInCommit", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:    "initial",
			FileName:  "initial_file",
			Locations: []testgit.Location{testgit.LocationLocal},
			Message:   "initial file commit",
		})
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		runtime.StageFiles("f1.txt", "f2.txt")
		runtime.CommitStagedChanges("stuff")
		commits := runtime.Commits([]string{}, runtime.Config.NormalConfig.Lineage, configdomain.OrderAsc)
		must.Len(t, 2, commits)
		fileNames := runtime.FilesInCommit(commits[1].SHA)
		must.Eq(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run("HasBranchesOutOfSync", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateBranch("branch1", "main")
			devRepo.CheckoutBranch("branch1")
			devRepo.CreateFile("file1", "content")
			devRepo.StageFiles("file1")
			devRepo.CommitStagedChanges("stuff")
			devRepo.PushBranchToRemote("branch1", gitdomain.RemoteOrigin)
			have, _ := devRepo.HasBranchesOutOfSync()
			must.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateBranch("branch1", "main")
			devRepo.PushBranch()
			devRepo.CreateFile("file1", "content")
			devRepo.StageFiles("file1")
			devRepo.CommitStagedChanges("stuff")
			have, _ := devRepo.HasBranchesOutOfSync()
			must.True(t, have)
		})

		t.Run("branch is behind", func(t *testing.T) {
			t.Parallel()
			fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateBranch("branch1", "main")
			devRepo.PushBranch()
			originRepo := fixture.OriginRepo.GetOrPanic()
			originRepo.CheckoutBranch("main")
			originRepo.CreateFile("file1", "content")
			originRepo.StageFiles("file1")
			originRepo.CommitStagedChanges("stuff")
			originRepo.CheckoutBranch("initial")
			devRepo.Fetch()
			have, _ := devRepo.HasBranchesOutOfSync()
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
			t.Parallel()
			runtime := testruntime.Create(t)
			must.NoError(t, runtime.Config.SetMainBranch("main", runtime.TestRunner))
			must.Error(t, runtime.VerifyNoGitTownConfiguration())
		})
		t.Run("the perennial branches are configured", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			must.NoError(t, gitconfig.SetPerennialBranches(runtime.TestRunner, gitdomain.NewLocalBranchNames("qa"), configdomain.ConfigScopeLocal))
			must.Error(t, runtime.VerifyNoGitTownConfiguration())
		})
		t.Run("branch lineage is configured", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("main", "initial")
			runtime.CreateFeatureBranch("foo", "main")
			must.Error(t, runtime.VerifyNoGitTownConfiguration())
		})
	})

	t.Run("LocalBranchesMainFirst", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		repoDir := t.TempDir()
		runner := testruntime.Clone(origin.TestRunner, repoDir)
		initial := gitdomain.NewLocalBranchName("initial")
		runner.CreateBranch("b1", initial.BranchName())
		runner.CreateBranch("b2", initial.BranchName())
		origin.CreateBranch("b3", initial.BranchName())
		runner.Fetch()
		branches, _, err := runner.LocalBranchesMainFirst(initial)
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "b1", "b2")
		must.Eq(t, want, branches)
	})

	t.Run("PushBranchToRemote", func(t *testing.T) {
		t.Parallel()
		dev := testruntime.Create(t)
		origin := testruntime.Create(t)
		dev.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		dev.CreateBranch("b1", "initial")
		dev.PushBranchToRemote("b1", gitdomain.RemoteOrigin)
		branches, _, err := origin.LocalBranchesMainFirst("initial")
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "b1")
		must.Eq(t, want, branches)
	})

	t.Run("RemoveBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("b1", "initial")
		branches, _, err := runtime.LocalBranchesMainFirst("initial")
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "b1")
		must.Eq(t, want, branches)
		runtime.RemoveBranch("b1")
		branches, _, err = runtime.LocalBranchesMainFirst("initial")
		must.NoError(t, err)
		wantBranches := gitdomain.NewLocalBranchNames("initial")
		must.Eq(t, wantBranches, branches)
	})

	t.Run("RemoveRemote", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		repo.RemoveRemote(gitdomain.RemoteOrigin)
		remotes, err := repo.Git.Remotes(repo.TestRunner)
		must.NoError(t, err)
		must.Len(t, 0, remotes)
	})

	t.Run("SHAForCommit", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		repo.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "bar",
			FileName:    "foo",
			Message:     "commit",
		})
		shas := repo.SHAsForCommit("commit")
		must.EqOp(t, 1, len(shas))
		sha := shas.First()
		must.EqOp(t, 40, len(sha))
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
