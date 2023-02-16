//nolint:testpackage
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

func TestGitEnvironment(t *testing.T) {
	t.Parallel()
	t.Run(".CloneGitEnvironment()", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
		assert.NoError(t, err)
		cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
		assert.NoError(t, err)
		assertIsNormalGitRepo(t, filepath.Join(dir, "cloned", "origin"))
		assertIsNormalGitRepo(t, filepath.Join(dir, "cloned", "developer"))
		assertHasGitBranch(t, filepath.Join(dir, "cloned", "developer"), "main")
		// check pushing
		err = cloned.DevRepo.PushBranch(git.PushArgs{BranchName: "main", Remote: config.OriginRemote})
		assert.NoError(t, err)
	})

	t.Run(".NewStandardGitEnvironment()", func(t *testing.T) {
		t.Parallel()
		gitEnvRootDir := t.TempDir()
		result, err := NewStandardGitEnvironment(gitEnvRootDir)
		assert.NoError(t, err)
		// verify the origin repo
		assertIsNormalGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
		branch, err := result.OriginRepo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch, "the origin should be at the initial branch so that we can push to it")
		// verify the developer repo
		assertIsNormalGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
		assertHasGlobalGitConfiguration(t, gitEnvRootDir)
		branch, err = result.DevRepo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "main", branch)
	})

	t.Run(".Branches()", func(t *testing.T) {
		t.Run("different branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create GitEnvironment instance
			dir := t.TempDir()
			gitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, ""))
			assert.NoError(t, err)
			// create the branches
			err = gitEnv.DevRepo.CreateBranch("d1", "main")
			assert.NoError(t, err)
			err = gitEnv.DevRepo.CreateBranch("d2", "main")
			assert.NoError(t, err)
			err = gitEnv.OriginRepo.CreateBranch("o1", "initial")
			assert.NoError(t, err)
			err = gitEnv.OriginRepo.CreateBranch("o2", "initial")
			assert.NoError(t, err)
			// get branches
			table, err := gitEnv.Branches()
			assert.NoError(t, err)
			// verify
			expected := "| REPOSITORY | BRANCHES     |\n| local      | main, d1, d2 |\n| origin     | main, o1, o2 |\n"
			assert.Equal(t, expected, table.String())
		})

		t.Run("same branches in dev and origin repo", func(t *testing.T) {
			t.Parallel()
			// create GitEnvironment instance
			dir := t.TempDir()
			gitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, ""))
			assert.NoError(t, err)
			// create the branches
			err = gitEnv.DevRepo.CreateBranch("b1", "main")
			assert.NoError(t, err)
			err = gitEnv.DevRepo.CreateBranch("b2", "main")
			assert.NoError(t, err)
			err = gitEnv.OriginRepo.CreateBranch("b1", "main")
			assert.NoError(t, err)
			err = gitEnv.OriginRepo.CreateBranch("b2", "main")
			assert.NoError(t, err)
			// get branches
			table, err := gitEnv.Branches()
			assert.NoError(t, err)
			// verify
			expected := "| REPOSITORY    | BRANCHES     |\n| local, origin | main, b1, b2 |\n"
			assert.Equal(t, expected, table.String())
		})
	})

	t.Run(".CreateCommits()", func(t *testing.T) {
		t.Parallel()
		// create GitEnvironment instance
		dir := t.TempDir()
		memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
		assert.NoError(t, err)
		cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
		assert.NoError(t, err)
		// create the commits
		err = cloned.CreateCommits([]git.Commit{
			{
				Branch:      "main",
				FileName:    "local-file",
				FileContent: "lc",
				Locations:   []string{"local"},
				Message:     "local commit",
			},
			{
				Branch:      "main",
				FileName:    "origin-file",
				FileContent: "rc",
				Locations:   []string{"origin"},
				Message:     "origin commit",
			},
			{
				Branch:      "main",
				FileName:    "loc-rem-file",
				FileContent: "lrc",
				Locations:   []string{"local", "origin"},
				Message:     "local and origin commit",
			},
		})
		assert.NoError(t, err)
		// verify local commits
		commits, err := cloned.DevRepo.Commits([]string{"FILE NAME", "FILE CONTENT"})
		assert.NoError(t, err)
		assert.Len(t, commits, 2)
		assert.Equal(t, "local commit", commits[0].Message)
		assert.Equal(t, "local-file", commits[0].FileName)
		assert.Equal(t, "lc", commits[0].FileContent)
		assert.Equal(t, "local and origin commit", commits[1].Message)
		assert.Equal(t, "loc-rem-file", commits[1].FileName)
		assert.Equal(t, "lrc", commits[1].FileContent)
		// verify origin commits
		commits, err = cloned.OriginRepo.Commits([]string{"FILE NAME", "FILE CONTENT"})
		assert.NoError(t, err)
		assert.Len(t, commits, 2)
		assert.Equal(t, "origin commit", commits[0].Message)
		assert.Equal(t, "origin-file", commits[0].FileName)
		assert.Equal(t, "rc", commits[0].FileContent)
		assert.Equal(t, "local and origin commit", commits[1].Message)
		assert.Equal(t, "loc-rem-file", commits[1].FileName)
		assert.Equal(t, "lrc", commits[1].FileContent)
		// verify origin is at "initial" branch
		branch, err := cloned.OriginRepo.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", branch)
	})

	t.Run(".CreateOriginBranch()", func(t *testing.T) {
		t.Parallel()
		// create GitEnvironment instance
		dir := t.TempDir()
		memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
		assert.NoError(t, err)
		cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
		assert.NoError(t, err)
		// create the origin branch
		err = cloned.CreateOriginBranch("b1", "main")
		assert.NoError(t, err)
		// verify it is in the origin branches
		branches, err := cloned.OriginRepo.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.Contains(t, branches, "b1")
		// verify it isn't in the local branches
		branches, err = cloned.DevRepo.LocalBranchesMainFirst()
		assert.NoError(t, err)
		assert.NotContains(t, branches, "b1")
	})

	t.Run(".CommitTable()", func(t *testing.T) {
		t.Run("without upstream repo", func(t *testing.T) {
			t.Parallel()
			// create GitEnvironment instance
			dir := t.TempDir()
			memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
			assert.NoError(t, err)
			cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
			assert.NoError(t, err)
			// create a few commits
			err = cloned.DevRepo.CreateCommit(git.Commit{
				Branch:      "main",
				FileName:    "local-origin.md",
				FileContent: "one",
				Message:     "local-origin",
			})
			assert.NoError(t, err)
			err = cloned.DevRepo.PushBranch(git.PushArgs{BranchName: "main", Remote: config.OriginRemote})
			assert.NoError(t, err)
			err = cloned.OriginRepo.CreateCommit(git.Commit{
				Branch:      "main",
				FileName:    "origin.md",
				FileContent: "two",
				Message:     "2",
			})
			assert.NoError(t, err)
			// get the CommitTable
			table, err := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
			assert.NoError(t, err)
			assert.Len(t, table.Cells, 3)
			assert.Equal(t, table.Cells[1][0], "local, origin")
			assert.Equal(t, table.Cells[1][1], "local-origin.md")
			assert.Equal(t, table.Cells[1][2], "one")
			assert.Equal(t, table.Cells[2][0], "origin")
			assert.Equal(t, table.Cells[2][1], "origin.md")
			assert.Equal(t, table.Cells[2][2], "two")
		})

		t.Run("with upstream repo", func(t *testing.T) {
			t.Parallel()
			// create GitEnvironment instance
			dir := t.TempDir()
			memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
			assert.NoError(t, err)
			cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
			assert.NoError(t, err)
			err = cloned.AddUpstream()
			assert.NoError(t, err)
			// create a few commits
			err = cloned.DevRepo.CreateCommit(git.Commit{
				Branch:      "main",
				FileName:    "local.md",
				FileContent: "one",
				Message:     "local",
			})
			assert.NoError(t, err)
			err = cloned.UpstreamRepo.CreateCommit(git.Commit{
				Branch:      "main",
				FileName:    "upstream.md",
				FileContent: "two",
				Message:     "2",
			})
			assert.NoError(t, err)
			// get the CommitTable
			table, err := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
			assert.NoError(t, err)
			assert.Len(t, table.Cells, 3)
			assert.Equal(t, table.Cells[1][0], "local")
			assert.Equal(t, table.Cells[1][1], "local.md")
			assert.Equal(t, table.Cells[1][2], "one")
			assert.Equal(t, table.Cells[2][0], "upstream")
			assert.Equal(t, table.Cells[2][1], "upstream.md")
			assert.Equal(t, table.Cells[2][2], "two")
		})
	})

	t.Run(".Remove()", func(t *testing.T) {
		t.Parallel()
		// create GitEnvironment instance
		dir := t.TempDir()
		memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
		assert.NoError(t, err)
		cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
		assert.NoError(t, err)
		// remove it
		err = cloned.Remove()
		assert.NoError(t, err)
		// verify
		_, err = os.Stat(cloned.Dir)
		assert.True(t, os.IsNotExist(err))
	})
}
