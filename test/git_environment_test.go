//nolint:testpackage
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/stretchr/testify/assert"
)

func TestGitEnvironment_CloneGitEnvironment(t *testing.T) {
	t.Parallel()
	dir := CreateTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
	assert.NoError(t, err)
	cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
	assert.NoError(t, err)
	assertIsNormalGitRepo(t, filepath.Join(dir, "cloned", "origin"))
	assertIsNormalGitRepo(t, filepath.Join(dir, "cloned", "developer"))
	assertHasGitBranch(t, filepath.Join(dir, "cloned", "developer"), "main")
	// check pushing
	err = cloned.DevRepo.PushBranchSetUpstream("main")
	assert.NoError(t, err)
}

func TestGitEnvironment_NewStandardGitEnvironment(t *testing.T) {
	t.Parallel()
	gitEnvRootDir := CreateTempDir(t)
	result, err := NewStandardGitEnvironment(gitEnvRootDir)
	assert.NoError(t, err)
	// verify the origin repo
	assertIsNormalGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
	branch, err := result.OriginRepo.CurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "master", branch, "the origin should be at the master branch so that we can push to it")
	// verify the developer repo
	assertIsNormalGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
	assertHasGlobalGitConfiguration(t, gitEnvRootDir)
	branch, err = result.DevRepo.CurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "main", branch)
}

func TestGitEnvironment_Branches_Different(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
	gitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, ""))
	assert.NoError(t, err)
	// create the branches
	err = gitEnv.DevRepo.CreateBranch("d1", "main")
	assert.NoError(t, err)
	err = gitEnv.DevRepo.CreateBranch("d2", "main")
	assert.NoError(t, err)
	err = gitEnv.OriginRepo.CreateBranch("o1", "master")
	assert.NoError(t, err)
	err = gitEnv.OriginRepo.CreateBranch("o2", "master")
	assert.NoError(t, err)
	// get branches
	table, err := gitEnv.Branches()
	assert.NoError(t, err)
	// verify
	expected := "| REPOSITORY | BRANCHES             |\n| local      | main, d1, d2         |\n| remote     | main, master, o1, o2 |\n"
	assert.Equal(t, expected, table.String())
}

func TestGitEnvironment_Branches_Same(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
	assert.NoError(t, err)
	cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
	assert.NoError(t, err)
	// create the branches
	err = cloned.DevRepo.CreateBranch("b1", "main")
	assert.NoError(t, err)
	err = cloned.DevRepo.CreateBranch("b2", "main")
	assert.NoError(t, err)
	err = cloned.OriginRepo.CreateBranch("b1", "master")
	assert.NoError(t, err)
	err = cloned.OriginRepo.CreateBranch("b2", "master")
	assert.NoError(t, err)
	// get branches
	table, err := cloned.Branches()
	assert.NoError(t, err)
	// verify
	expected := "| REPOSITORY | BRANCHES     |\n| local, remote      | main, b1, b2 |\n"
	assert.Equal(t, expected, table.String())
}

func TestGitEnvironment_CreateCommits(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
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
			FileName:    "remote-file",
			FileContent: "rc",
			Locations:   []string{"remote"},
			Message:     "remote commit",
		},
		{
			Branch:      "main",
			FileName:    "loc-rem-file",
			FileContent: "lrc",
			Locations:   []string{"local", "remote"},
			Message:     "local and remote commit",
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
	assert.Equal(t, "local and remote commit", commits[1].Message)
	assert.Equal(t, "loc-rem-file", commits[1].FileName)
	assert.Equal(t, "lrc", commits[1].FileContent)
	// verify remote commits
	commits, err = cloned.OriginRepo.Commits([]string{"FILE NAME", "FILE CONTENT"})
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "remote commit", commits[0].Message)
	assert.Equal(t, "remote-file", commits[0].FileName)
	assert.Equal(t, "rc", commits[0].FileContent)
	assert.Equal(t, "local and remote commit", commits[1].Message)
	assert.Equal(t, "loc-rem-file", commits[1].FileName)
	assert.Equal(t, "lrc", commits[1].FileContent)
	// verify origin is at master
	branch, err := cloned.OriginRepo.CurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "master", branch)
}

func TestGitEnvironment_CreateRemoteBranch(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
	assert.NoError(t, err)
	cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
	assert.NoError(t, err)
	// create the remote mranch
	err = cloned.CreateRemoteBranch("b1", "main")
	assert.NoError(t, err)
	// verify it is in the remote branches
	branches, err := cloned.OriginRepo.LocalBranchesMainFirst()
	assert.NoError(t, err)
	assert.Contains(t, branches, "b1")
	// verify it isn't in the local branches
	branches, err = cloned.DevRepo.LocalBranchesMainFirst()
	assert.NoError(t, err)
	assert.NotContains(t, branches, "b1")
}

func TestGitEnvironment_CommitTable(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
	assert.NoError(t, err)
	cloned, err := CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))
	assert.NoError(t, err)
	// create a few commits
	err = cloned.DevRepo.CreateCommit(git.Commit{
		Branch:      "main",
		FileName:    "local-remote.md",
		FileContent: "one",
		Message:     "local-remote",
	})
	assert.NoError(t, err)
	err = cloned.DevRepo.PushBranchSetUpstream("main")
	assert.NoError(t, err)
	err = cloned.OriginRepo.CreateCommit(git.Commit{
		Branch:      "main",
		FileName:    "remote.md",
		FileContent: "two",
		Message:     "2",
	})
	assert.NoError(t, err)
	// get the CommitTable
	table, err := cloned.CommitTable([]string{"LOCATION", "FILE NAME", "FILE CONTENT"})
	assert.NoError(t, err)
	assert.Len(t, table.Cells, 3)
	assert.Equal(t, table.Cells[1][0], "local, remote")
	assert.Equal(t, table.Cells[1][1], "local-remote.md")
	assert.Equal(t, table.Cells[1][2], "one")
	assert.Equal(t, table.Cells[2][0], "remote")
	assert.Equal(t, table.Cells[2][1], "remote.md")
	assert.Equal(t, table.Cells[2][2], "two")
}

func TestGitEnvironment_CommitTable_Upstream(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
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
}

func TestGitEnvironment_Remove(t *testing.T) {
	t.Parallel()
	// create GitEnvironment instance
	dir := CreateTempDir(t)
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
}
