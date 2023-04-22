//nolint:testpackage
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/src/git"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Parallel()
	t.Run("newRunner", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		workingDir := filepath.Join(dir, "working")
		homeDir := filepath.Join(dir, "home")
		binDir := filepath.Join(dir, "bin")
		repo := newRunner(workingDir, homeDir, binDir)
		assert.Equal(t, workingDir, repo.workingDir)
		assert.Equal(t, homeDir, repo.homeDir)
		assert.Equal(t, binDir, repo.binDir)
	})

	t.Run(".AddRemote()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		remotes, err := runner.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, remotes)
		origin := CreateRunner(t)
		err = runner.AddRemote(config.OriginRemote, origin.workingDir)
		assert.NoError(t, err)
		remotes, err = runner.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})

	t.Run(".Clone()", func(t *testing.T) {
		t.Parallel()
		origin := CreateRunner(t)
		clonedPath := filepath.Join(origin.workingDir, "cloned")
		cloned, err := origin.Clone(clonedPath)
		assert.NoError(t, err)
		assert.Equal(t, clonedPath, cloned.workingDir)
		assertIsNormalGitRepo(t, clonedPath)
	})

	t.Run(".Commits()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		err := runner.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file1",
			FileContent: "hello",
			Message:     "first commit",
		})
		assert.NoError(t, err)
		err = runner.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "file2",
			FileContent: "hello again",
			Message:     "second commit",
		})
		assert.NoError(t, err)
		commits, err := runner.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
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
		origin := CreateRunner(t)
		repoDir := filepath.Join(t.TempDir(), "repo") // need a non-existing directory
		err := CopyDirectory(origin.workingDir, repoDir)
		assert.NoError(t, err)
		runner := newRunner(repoDir, repoDir, "")
		err = runner.AddRemote(config.OriginRemote, origin.workingDir)
		assert.NoError(t, err)
		err = runner.Fetch()
		assert.NoError(t, err)
		err = runner.ConnectTrackingBranch("initial")
		assert.NoError(t, err)
		err = runner.PushBranch()
		assert.NoError(t, err)
	})

	t.Run(".CreateBranch()", func(t *testing.T) {
		t.Run("simple branch name", func(t *testing.T) {
			t.Parallel()
			runner := CreateRunner(t)
			err := runner.CreateBranch("branch1", "initial")
			assert.NoError(t, err)
			currentBranch, err := runner.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, "initial", currentBranch)
			branches, err := runner.LocalBranchesMainFirst("initial")
			assert.NoError(t, err)
			assert.Equal(t, []string{"initial", "branch1"}, branches)
		})

		t.Run("branch name with slashes", func(t *testing.T) {
			t.Parallel()
			runner := CreateRunner(t)
			err := runner.CreateBranch("my/feature", "initial")
			assert.NoError(t, err)
			currentBranch, err := runner.CurrentBranch()
			assert.NoError(t, err)
			assert.Equal(t, "initial", currentBranch)
			branches, err := runner.LocalBranchesMainFirst("initial")
			assert.NoError(t, err)
			assert.Equal(t, []string{"initial", "my/feature"}, branches)
		})
	})

	t.Run(".CreateChildFeatureBranch()", func(t *testing.T) {
		t.Parallel()
		runner := CreateTestGitTownRunner(t)
		err := runner.CreateFeatureBranch("f1")
		assert.NoError(t, err)
		err = runner.CreateChildFeatureBranch("f1a", "f1")
		assert.NoError(t, err)
		output, err := runner.BackendRunner.Run("git-town", "config")
		assert.NoError(t, err)
		has := strings.Contains(output, "Branch Ancestry:\n  main\n    f1\n      f1a")
		if !has {
			fmt.Printf("unexpected output: %s", output)
		}
		assert.True(t, has)
	})

	t.Run(".CreateCommit()", func(t *testing.T) {
		t.Run("minimal arguments", func(t *testing.T) {
			t.Parallel()
			runner := CreateRunner(t)
			err := runner.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
			})
			assert.NoError(t, err)
			commits, err := runner.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
			assert.NoError(t, err)
			assert.Len(t, commits, 1)
			assert.Equal(t, "hello.txt", commits[0].FileName)
			assert.Equal(t, "hello world", commits[0].FileContent)
			assert.Equal(t, "test commit", commits[0].Message)
			assert.Equal(t, "initial", commits[0].Branch)
		})

		t.Run("set the author", func(t *testing.T) {
			t.Parallel()
			runner := CreateRunner(t)
			err := runner.CreateCommit(git.Commit{
				Branch:      "initial",
				FileName:    "hello.txt",
				FileContent: "hello world",
				Message:     "test commit",
				Author:      "developer <developer@example.com>",
			})
			assert.NoError(t, err)
			commits, err := runner.Commits([]string{"FILE NAME", "FILE CONTENT"}, "initial")
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
			runner := CreateRunner(t)
			err := runner.CreateFile("filename", "content")
			assert.Nil(t, err, "cannot create file in repo")
			content, err := os.ReadFile(filepath.Join(runner.workingDir, "filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})

		t.Run("create file in subfolder", func(t *testing.T) {
			t.Parallel()
			runner := CreateRunner(t)
			err := runner.CreateFile("folder/filename", "content")
			assert.Nil(t, err, "cannot create file in repo")
			content, err := os.ReadFile(filepath.Join(runner.workingDir, "folder/filename"))
			assert.Nil(t, err, "cannot read file")
			assert.Equal(t, "content", string(content))
		})
	})

	t.Run(".CreatePerennialBranches()", func(t *testing.T) {
		t.Parallel()
		runner := CreateTestGitTownRunner(t)
		err := runner.CreatePerennialBranches("p1", "p2")
		assert.NoError(t, err)
		branches, err := runner.LocalBranchesMainFirst("main")
		assert.NoError(t, err)
		assert.Equal(t, []string{"main", "initial", "p1", "p2"}, branches)
		runner.Config.Reload()
		assert.True(t, runner.Config.IsPerennialBranch("p1"))
		assert.True(t, runner.Config.IsPerennialBranch("p2"))
	})

	t.Run(".Fetch()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		origin := CreateRunner(t)
		err := runner.AddRemote(config.OriginRemote, origin.workingDir)
		assert.NoError(t, err)
		err = runner.Fetch()
		assert.NoError(t, err)
	})

	t.Run(".FileContentInCommit()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		err := runner.CreateCommit(git.Commit{
			Branch:      "initial",
			FileName:    "hello.txt",
			FileContent: "hello world",
			Message:     "commit",
		})
		assert.NoError(t, err)
		commits, err := runner.CommitsInBranch("initial", []string{})
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		content, err := runner.FileContentInCommit(commits[0].SHA, "hello.txt")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", content)
	})

	t.Run(".FilesInCommit()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		err := runner.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		err = runner.CreateFile("f2.txt", "two")
		assert.NoError(t, err)
		err = runner.StageFiles("f1.txt", "f2.txt")
		assert.NoError(t, err)
		err = runner.CommitStagedChanges("stuff")
		assert.NoError(t, err)
		commits, err := runner.Commits([]string{}, "initial")
		assert.NoError(t, err)
		assert.Len(t, commits, 1)
		fileNames, err := runner.FilesInCommit(commits[0].SHA)
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, fileNames)
	})

	t.Run(".HasBranchesOutOfSync()", func(t *testing.T) {
		t.Run("branches are in sync", func(t *testing.T) {
			t.Parallel()
			env, err := NewStandardFixture(t.TempDir())
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
			env, err := NewStandardFixture(t.TempDir())
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
			env, err := NewStandardFixture(t.TempDir())
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
		runner := CreateRunner(t)
		err := runner.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		has, err := runner.HasFile("f1.txt", "one")
		assert.NoError(t, err)
		assert.True(t, has)
		_, err = runner.HasFile("f1.txt", "zonk")
		assert.Error(t, err)
		_, err = runner.HasFile("zonk.txt", "one")
		assert.Error(t, err)
	})

	t.Run(".HasGitTownConfigNow()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		res := runner.HasGitTownConfigNow()
		assert.False(t, res)
		err := runner.CreateBranch("main", "initial")
		assert.NoError(t, err)
		err = runner.CreateFeatureBranch("foo")
		assert.NoError(t, err)
		res = runner.HasGitTownConfigNow()
		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run(".PushBranchToRemote()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		origin := CreateRunner(t)
		err := runner.AddRemote(config.OriginRemote, origin.workingDir)
		assert.NoError(t, err)
		err = runner.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		err = runner.PushBranchToRemote("b1", config.OriginRemote)
		assert.NoError(t, err)
		branches, err := origin.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
	})

	t.Run(".RemoveBranch()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		err := runner.CreateBranch("b1", "initial")
		assert.NoError(t, err)
		branches, err := runner.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "b1"}, branches)
		err = runner.RemoveBranch("b1")
		assert.NoError(t, err)
		branches, err = runner.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial"}, branches)
	})

	t.Run(".RemoveRemote()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		origin := CreateRunner(t)
		err := runner.AddRemote(config.OriginRemote, origin.workingDir)
		assert.NoError(t, err)
		err = runner.RemoveRemote(config.OriginRemote)
		assert.NoError(t, err)
		remotes, err := runner.Remotes()
		assert.NoError(t, err)
		assert.Len(t, remotes, 0)
	})

	t.Run(".ShaForCommit()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		err := runner.CreateCommit(git.Commit{Branch: "initial", FileName: "foo", FileContent: "bar", Message: "commit"})
		assert.NoError(t, err)
		sha, err := runner.ShaForCommit("commit")
		assert.NoError(t, err)
		assert.Len(t, sha, 40)
	})

	t.Run(".UncommittedFiles()", func(t *testing.T) {
		t.Parallel()
		runner := CreateRunner(t)
		err := runner.CreateFile("f1.txt", "one")
		assert.NoError(t, err)
		err = runner.CreateFile("f2.txt", "two")
		assert.NoError(t, err)
		files, err := runner.UncommittedFiles()
		assert.NoError(t, err)
		assert.Equal(t, []string{"f1.txt", "f2.txt"}, files)
	})
}
