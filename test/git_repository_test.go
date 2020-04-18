package test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitRepository(t *testing.T) {
	rootDir := createTempDir(t)
	originPath := filepath.Join(rootDir, "origin")
	_, err := InitGitRepository(originPath, rootDir)
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := filepath.Join(rootDir, "cloned")
	_, err = CloneGitRepository(originPath, clonedPath, rootDir)
	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestInitGitRepository(t *testing.T) {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir)
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, repo.Dir)
	// ensure the Git repo works, i.e. we can commit into it
	repo.CreateFile("test.txt", "hello")
	repo.StageFile("test.txt")
	repo.CommitStagedChanges()
}

func TestNewGitRepository(t *testing.T) {
	dir := createTempDir(t)
	_ = NewGitRepository(dir, dir)
}

func TestGitRepository_Branches(t *testing.T) {
	repo := createTestGitTownRepo(t)
	assert.Nil(t, repo.CreateFeatureBranch("branch3", false), "cannot create branch3")
	assert.Nil(t, repo.CreateFeatureBranch("branch2", false), "cannot create branch2")
	assert.Nil(t, repo.CreateFeatureBranch("branch1", false), "cannot create branch1")
	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "branch2", "branch3", "master"}, branches)
}

func TestGitRepository_CreateFile(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

func TestGitRepository_CreateFile_InSubFolder(t *testing.T) {
	repo := createTestRepo(t)
	err := repo.CreateFile("folder/filename", "content")
	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(filepath.Join(repo.Dir, "folder/filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

// HELPERS

// createTestGitRepo creates a fully initialized Git repo including a master branch.
func createTestRepo(t *testing.T) GitRepository {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir)
	assert.Nil(t, err, "cannot initialize Git repow")
	err = repo.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
	})
	assert.Nil(t, err, "cannot create initial commit: %s")
	return repo
}

func createTestGitTownRepo(t *testing.T) GitRepository {
	repo := createTestRepo(t)
	err := repo.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "master"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.Nil(t, err)
	return repo
}
