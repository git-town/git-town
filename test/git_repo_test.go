package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitRepository(t *testing.T) {
	rootDir := CreateTempDir(t)
	originPath := filepath.Join(rootDir, "origin")
	_, err := InitGitRepository(originPath, rootDir, "")
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := filepath.Join(rootDir, "cloned")
	_, err = CloneGitRepo(originPath, clonedPath, rootDir, "")
	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestGitRepository_FileContentInCommit(t *testing.T) {
	repo := CreateTestRepo(t)
	err := repo.CreateCommit(Commit{
		Branch:      "master",
		FileName:    "hello.txt",
		FileContent: "hello world",
		Message:     "commit",
	})
	assert.Nil(t, err)
	commits, err := repo.commitsInBranch("master", []string{})
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	content, err := repo.FileContentInCommit(commits[0].SHA, "hello.txt")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", content)
}

func TestInitGitRepository(t *testing.T) {
	dir := CreateTempDir(t)
	repo, err := InitGitRepository(dir, dir, "")
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, repo.Dir)
	// ensure the Git repo works, i.e. we can commit into it
	err = repo.CreateFile("test.txt", "hello")
	assert.Nil(t, err)
	err = repo.StageFiles("test.txt")
	assert.Nil(t, err)
	err = repo.CommitStagedChanges(true)
	assert.Nil(t, err)
}

func TestNewGitRepository(t *testing.T) {
	dir := CreateTempDir(t)
	_ = NewGitRepository(dir, dir, NewMockingShell(dir, dir, ""))
}
