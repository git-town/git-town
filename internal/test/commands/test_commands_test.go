package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/datatable"
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

	t.Run("AddSubmodule", func(t *testing.T) {
		t.Skip("AddSubmodule requires git config protocol.file.allow=always which may not be set")
		t.Parallel()
		runtime := testruntime.Create(t)
		submoduleOrigin := testruntime.Create(t)
		runtime.AddSubmodule(submoduleOrigin.WorkingDir)
		output := runtime.MustQuery("git", "submodule", "status")
		must.StrContains(t, output, submoduleOrigin.WorkingDir)
	})

	t.Run("AddWorktree", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch-1", "initial")
		worktreePath := filepath.Join(t.TempDir(), "worktree1")
		runtime.AddWorktree(worktreePath, "branch-1")
		output := runtime.MustQuery("git", "worktree", "list")
		must.StrContains(t, output, worktreePath)
	})

	t.Run("AmendCommit", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "initial content",
			FileName:    "file1.txt",
			Message:     "initial commit message",
		})
		runtime.CreateFile("file2.txt", "new content")
		runtime.StageFiles("file2.txt")
		runtime.AmendCommit()
		commits := runtime.Commits([]string{"FILE NAME"}, runtime.Config.NormalConfig.Lineage, configdomain.OrderAsc)
		must.Len(t, 1, commits)
		files := runtime.FilesInCommit(commits[0].SHA)
		must.Eq(t, []string{"file1.txt", "file2.txt"}, files)
	})

	t.Run("CheckoutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch-1", "initial")
		runtime.CheckoutBranch("branch-1")
		currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
		must.EqOp(t, "branch-1", currentBranch)
	})

	t.Run("CommitSHA", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("feature1", "main")
		runtime.CheckoutBranch("feature1")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "feature1",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "test commit",
		})
		sha := runtime.CommitSHA(runtime, "test commit", "feature1", "main")
		must.EqOp(t, 40, len(sha.String()))
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

	t.Run("CommitStagedChanges", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("file1.txt", "content")
		runtime.StageFiles("file1.txt")
		runtime.CommitStagedChanges("test commit")
		commits := runtime.Commits([]string{}, runtime.Config.NormalConfig.Lineage, configdomain.OrderAsc)
		must.Len(t, 1, commits)
		must.EqOp(t, "test commit", commits[0].Message)
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

	t.Run("CommitsInBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch-1", "initial")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch-1",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "first commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch-1",
			FileContent: "content",
			FileName:    "file2.txt",
			Message:     "second commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content",
			FileName:    "file3.txt",
			Message:     "third commit",
		})
		commits := runtime.CommitsInBranch("branch-1", None[gitdomain.BranchName](), []string{})
		must.Len(t, 2, commits)
		must.EqOp(t, "branch-1", commits[0].Branch)
		must.EqOp(t, "first commit", commits[0].Message)
		must.EqOp(t, "branch-1", commits[1].Branch)
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

	t.Run("CreateAndCheckoutBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateAndCheckoutBranch("branch-1", "initial")
		currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
		must.EqOp(t, "branch-1", currentBranch)
		branches, _, err := runtime.LocalBranchesMainFirst("initial")
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("initial", "branch-1")
		must.Eq(t, want, branches)
	})

	t.Run("CreateAndCheckoutFeatureBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateAndCheckoutFeatureBranch("feature1", "main")
		currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
		must.EqOp(t, "feature1", currentBranch)
		runtime.Config.Reload(runtime.TestRunner)
		parent := runtime.Config.NormalConfig.Lineage.Parent("feature1")
		must.EqOp(t, "main", parent.GetOrPanic())
		must.EqOp(t, configdomain.BranchTypeFeatureBranch, runtime.Config.BranchType("feature1"))
	})

	t.Run("CreateBranch", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateBranch("branch-1", "initial")
			currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
			must.EqOp(t, "initial", currentBranch)
			branches, _, err := runtime.LocalBranchesMainFirst("initial")
			must.NoError(t, err)
			want := gitdomain.NewLocalBranchNames("initial", "branch-1")
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

	t.Run("CreateBranchOfType", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateBranchOfType("observed-1", Some(gitdomain.NewLocalBranchName("main")), configdomain.BranchTypeObservedBranch)
		runtime.Config.Reload(runtime.TestRunner)
		parent := runtime.Config.NormalConfig.Lineage.Parent("observed-1")
		must.EqOp(t, "main", parent.GetOrPanic())
		must.EqOp(t, configdomain.BranchTypeObservedBranch, runtime.Config.BranchType("observed-1"))
		branches, _, err := runtime.LocalBranchesMainFirst("main")
		must.NoError(t, err)
		must.True(t, branches.Contains("observed-1"))
	})

	t.Run("CreateChildBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateChildBranch("child-1", "main")
		currentBranch := asserts.NoError1(runtime.Git.CurrentBranch(runtime)).GetOrPanic()
		must.EqOp(t, "child-1", currentBranch)
		runtime.Config.Reload(runtime.TestRunner)
		parent := runtime.Config.NormalConfig.Lineage.Parent("child-1")
		must.EqOp(t, "main", parent.GetOrPanic())
		must.EqOp(t, configdomain.BranchTypeFeatureBranch, runtime.Config.BranchType("child-1"))
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

	t.Run("CreateFolder", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFolder("myfolder")
		info, err := os.Stat(filepath.Join(runtime.WorkingDir, "myfolder"))
		must.NoError(t, err)
		must.True(t, info.IsDir())
	})

	t.Run("CreateLocalBranchUsingGitTown", func(t *testing.T) {
		t.Parallel()
		fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
		devRepo := fixture.DevRepo.GetOrPanic()
		branchSetup := datatable.BranchSetup{
			Name:       "feature-1",
			Parent:     Some(gitdomain.NewLocalBranchName("main")),
			BranchType: Some(configdomain.BranchTypeFeatureBranch),
			Locations:  []testgit.Location{testgit.LocationLocal, testgit.LocationOrigin},
		}
		devRepo.CreateLocalBranchUsingGitTown(branchSetup)
		branches, _, err := devRepo.LocalBranchesMainFirst("main")
		must.NoError(t, err)
		must.True(t, branches.Contains("feature-1"))
		devRepo.Config.Reload(devRepo.TestRunner)
		parent := devRepo.Config.NormalConfig.Lineage.Parent("feature-1")
		must.EqOp(t, "main", parent.GetOrPanic())
		must.EqOp(t, configdomain.BranchTypeFeatureBranch, devRepo.Config.BranchType("feature-1"))
	})

	t.Run("CreateStandaloneTag", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateStandaloneTag("v1.0.0")
		tags := runtime.Tags()
		must.Eq(t, []string{"v1.0.0"}, tags)
	})

	t.Run("CreateTag", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateTag("v1.0.0")
		tags := runtime.Tags()
		must.Eq(t, []string{"v1.0.0"}, tags)
	})

	t.Run("ExistingParent", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("feature-1", "main")
		runtime.CreateFeatureBranch("feature-2", "feature-1")
		runtime.Config.Reload(runtime.TestRunner)
		lineage := runtime.Config.NormalConfig.Lineage
		parent := runtime.ExistingParent("feature-2", lineage)
		must.EqOp(t, "feature-1", parent.GetOrPanic())
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.Create(t)
		origin := testruntime.Create(t)
		repo.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		repo.Fetch()
	})

	t.Run("FileContent", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("file1.txt", "test content")
		content := runtime.FileContent("file1.txt")
		must.EqOp(t, "test content", content)
	})

	t.Run("FileContentErr", func(t *testing.T) {
		t.Run("file exists", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			runtime.CreateFile("file1.txt", "test content")
			content, err := runtime.FileContentErr("file1.txt")
			must.NoError(t, err)
			must.EqOp(t, "test content", content)
		})
		t.Run("file does not exist", func(t *testing.T) {
			t.Parallel()
			runtime := testruntime.Create(t)
			_, err := runtime.FileContentErr("nonexistent.txt")
			must.Error(t, err)
		})
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

	t.Run("FilesInBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content1",
			FileName:    "file1.txt",
			Message:     "first commit",
		})
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content2",
			FileName:    "file2.txt",
			Message:     "second commit",
		})
		files := runtime.FilesInBranch("initial")
		must.Eq(t, []string{"file1.txt", "file2.txt"}, files)
	})

	t.Run("FilesInBranches", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "commit",
		})
		runtime.CreateAndCheckoutBranch("branch-1", "initial")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch-1",
			FileContent: "content2",
			FileName:    "file2.txt",
			Message:     "branch commit",
		})
		table := runtime.FilesInBranches("initial")
		want := `
| BRANCH   | NAME      | CONTENT  |
| initial  | file1.txt | content  |
| branch-1 | file1.txt | content  |
|          | file2.txt | content2 |
`[1:]
		must.Eq(t, want, table.String())
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

	t.Run("FilesInWorkspace", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("file1.txt", "content")
		runtime.CreateFile("file2.txt", "content")
		files := runtime.FilesInWorkspace()
		must.Eq(t, []string{"file1.txt", "file2.txt"}, files)
	})

	t.Run("HasBranchesOutOfSync", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateBranch("branch-1", "main")
			devRepo.CheckoutBranch("branch-1")
			devRepo.CreateFile("file1", "content")
			devRepo.StageFiles("file1")
			devRepo.CommitStagedChanges("stuff")
			devRepo.PushBranchToRemote("branch-1", gitdomain.RemoteOrigin)
			have, _ := devRepo.HasBranchesOutOfSync()
			must.False(t, have)
		})

		t.Run("branch is ahead", func(t *testing.T) {
			t.Parallel()
			fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
			devRepo := fixture.DevRepo.GetOrPanic()
			devRepo.CreateAndCheckoutBranch("branch-1", "main")
			devRepo.PushBranchToRemote("branch-1", gitdomain.RemoteOrigin)
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
			devRepo.CreateBranch("branch-1", "main")
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

	t.Run("Lineage", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("feature-1", "main")
		lineage := runtime.Lineage()
		parent := lineage.Parent("feature-1").GetOrPanic()
		must.EqOp(t, "main", parent)
	})

	t.Run("LineageText", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("feature-1", "main")
		lineage := runtime.Lineage()
		text := runtime.LineageText(lineage)
		want := `
main
  feature-1`[1:]
		must.Eq(t, want, text)
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
		devRepo := fixture.DevRepo.GetOrPanic()
		originRepo := fixture.OriginRepo.GetOrPanic()
		devRepo.CreateBranch("b1", "main")
		devRepo.CreateBranch("b2", "main")
		originRepo.CreateBranch("remote-branch", "main")
		devRepo.Fetch()
		have, _, err := devRepo.LocalBranches()
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("b1", "b2", "main")
		must.Eq(t, want, have)
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

	t.Run("LocalBranchesOrderedHierarchically", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("feature-1", "main")
		runtime.CreateFeatureBranch("feature-2", "feature-1")
		lineage := runtime.Config.NormalConfig.Lineage
		have, _, err := runtime.LocalBranchesOrderedHierarchically(lineage, configdomain.OrderAsc)
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchNames("main", "feature-1", "feature-2", "initial")
		must.Eq(t, want, have)
	})

	t.Run("MergeBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateBranch("branch-1", "initial")
		runtime.CheckoutBranch("branch-1")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch-1",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "commit on branch-1",
		})
		runtime.CheckoutBranch("initial")
		err := runtime.MergeBranch("branch-1")
		must.NoError(t, err)
		files := runtime.FilesInBranch("initial")
		must.Eq(t, []string{"file1.txt"}, files)
	})

	t.Run("PushBranch", func(t *testing.T) {
		t.Parallel()
		fixture := fixture.NewMemoized(t.TempDir()).AsFixture()
		dev := fixture.DevRepo.GetOrPanic()
		dev.CreateAndCheckoutBranch("b1", "main")
		dev.PushBranchToRemote("b1", gitdomain.RemoteOrigin)
		dev.CreateCommit(testgit.Commit{
			Branch:      "b1",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "commit",
		})
		dev.PushBranch()
		devCommits := dev.CommitsInBranch("b1", None[gitdomain.BranchName](), []string{})
		must.Len(t, 1, devCommits)
		must.EqOp(t, "commit", devCommits[0].Message)
		origin := fixture.OriginRepo.GetOrPanic()
		originCommits := origin.CommitsInBranch("b1", None[gitdomain.BranchName](), []string{})
		must.Len(t, 1, originCommits)
		must.EqOp(t, devCommits[0].Message, originCommits[0].Message)
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

	t.Run("RebaseAgainstBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content1",
			FileName:    "file1.txt",
			Message:     "base commit",
		})
		runtime.CreateBranch("branch-1", "initial")
		runtime.CheckoutBranch("branch-1")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "branch-1",
			FileContent: "content2",
			FileName:    "file2.txt",
			Message:     "branch commit",
		})
		runtime.CheckoutBranch("initial")
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content3",
			FileName:    "file3.txt",
			Message:     "another base commit",
		})
		runtime.CheckoutBranch("branch-1")
		err := runtime.RebaseAgainstBranch("initial")
		must.NoError(t, err)
		commits := runtime.CommitsInBranch("branch-1", None[gitdomain.BranchName](), []string{"FILE NAME"})
		must.Len(t, 3, commits)
		must.EqOp(t, "base commit", commits[0].Message)
		must.EqOp(t, "another base commit", commits[1].Message)
		must.EqOp(t, "branch commit", commits[2].Message)
	})

	t.Run("Reload", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.CreateGitTown(t)
		runtime.CreateFeatureBranch("feature1", "main")
		runtime.Reload()
		parent := runtime.Config.NormalConfig.Lineage.Parent("feature1")
		must.EqOp(t, "main", parent.GetOrPanic())
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

	t.Run("RemovePerennialBranchConfiguration", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		must.NoError(t, gitconfig.SetPerennialBranches(runtime.TestRunner, gitdomain.NewLocalBranchNames("qa"), configdomain.ConfigScopeLocal))
		must.NoError(t, runtime.RemovePerennialBranchConfiguration())
		runtime.Config.Reload(runtime.TestRunner)
		perennials := runtime.Config.NormalConfig.PerennialBranches
		must.Len(t, 0, perennials)
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

	t.Run("RemoveUnnecessaryFiles", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		_, err := os.Stat(filepath.Join(runtime.WorkingDir, ".git", "hooks"))
		must.NoError(t, err)
		runtime.RemoveUnnecessaryFiles()
		_, err = os.Stat(filepath.Join(runtime.WorkingDir, ".git", "hooks"))
		must.Error(t, err)
	})

	t.Run("RenameRemote", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		origin := testruntime.Create(t)
		runtime.AddRemote(gitdomain.RemoteOrigin, origin.WorkingDir)
		runtime.RenameRemote("origin", "upstream")
		remotes, err := runtime.Git.Remotes(runtime)
		must.NoError(t, err)
		must.Eq(t, gitdomain.Remotes{gitdomain.RemoteUpstream}, remotes)
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

	t.Run("SHAforBranch", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "test commit",
		})
		branchSHA := runtime.SHAforBranch("initial")
		commitSHA := runtime.SHAsForCommit("test commit").First()
		must.EqOp(t, 40, len(branchSHA.String()))
		must.Eq(t, commitSHA, branchSHA)
	})

	t.Run("StageFiles", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("file1.txt", "content")
		runtime.CreateFile("file2.txt", "content")
		runtime.StageFiles("file1.txt", "file2.txt")
		output := runtime.MustQuery("git", "status", "--porcelain")
		want := `
A  file1.txt
A  file2.txt`[1:]
		must.EqOp(t, want, output)
	})

	t.Run("Tags", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateCommit(testgit.Commit{
			Branch:      "initial",
			FileContent: "content",
			FileName:    "file1.txt",
			Message:     "test commit",
		})
		runtime.CreateTag("v1.0.0")
		runtime.CreateTag("v2.0.0")
		tags := runtime.Tags()
		must.Eq(t, []string{"v1.0.0", "v2.0.0"}, tags)
	})

	t.Run("UncommittedFiles", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("f1.txt", "one")
		runtime.CreateFile("f2.txt", "two")
		files := runtime.UncommittedFiles()
		must.Eq(t, []string{"f1.txt", "f2.txt"}, files)
	})

	t.Run("UnstashOpenFiles", func(t *testing.T) {
		t.Parallel()
		runtime := testruntime.Create(t)
		runtime.CreateFile("file1.txt", "content")
		runtime.StashOpenFiles()
		uncommitted := runtime.UncommittedFiles()
		must.Len(t, 0, uncommitted)
		err := runtime.UnstashOpenFiles()
		must.NoError(t, err)
		files := runtime.UncommittedFiles()
		must.Eq(t, []string{"file1.txt"}, files)
	})
}
