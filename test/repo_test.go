package test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/stretchr/testify/assert"
)

func TestGitRepository_FileContentInCommit(t *testing.T) {
	repo := CreateRepo(t)
	err := repo.CreateCommit(git.Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "commit",
	})
	assert.Nil(t, err)
	commits, err := repo.CommitsInBranch("master", []string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	content, err := repo.FileContentInCommit(commits[0].SHA, "hello.txt")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", content)
}

func TestNewRepo(t *testing.T) {
	dir := CreateTempDir(t)
	workingDir := filepath.Join(dir, "working")
	homeDir := filepath.Join(dir, "home")
	binDir := filepath.Join(dir, "bin")
	repo := NewRepo(workingDir, homeDir, binDir)
	assert.Equal(t, workingDir, repo.shell.workingDir)
	assert.Equal(t, homeDir, repo.shell.homeDir)
	assert.Equal(t, binDir, repo.shell.binDir)
}

func TestRepo_Clone(t *testing.T) {
	origin := CreateRepo(t)
	clonedPath := filepath.Join(origin.shell.workingDir, "cloned")
	cloned, err := origin.Clone(clonedPath)
	assert.Nil(t, err)
	assert.Equal(t, clonedPath, cloned.shell.workingDir)
	assertIsNormalGitRepo(t, clonedPath)
}
