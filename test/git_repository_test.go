package test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitRepository(t *testing.T) {
	dir := createTempDir(t)
	_ = NewGitRepository(dir, dir)
}

func TestInitNormalGitRepository(t *testing.T) {
	dir := createTempDir(t)
	result, err := InitGitRepository(dir, dir)
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, result.Dir)
}

func TestCloneGitRepository(t *testing.T) {
	rootDir := createTempDir(t)
	originPath := path.Join(rootDir, "origin")
	_, err := InitGitRepository(originPath, rootDir)
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := path.Join(rootDir, "cloned")

	_, err = CloneGitRepository(originPath, clonedPath, rootDir)

	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestGitRepositoryBranches(t *testing.T) {
	repo := createTestGitTownRepo(t)
	assert.Nil(t, repo.CreateFeatureBranch("branch3"), "cannot create branch3")
	assert.Nil(t, repo.CreateFeatureBranch("branch2"), "cannot create branch2")
	assert.Nil(t, repo.CreateFeatureBranch("branch1"), "cannot create branch1")

	branches, err := repo.Branches()
	assert.Nil(t, err)
	assert.Equal(t, []string{"branch1", "branch2", "branch3", "master"}, branches)
}

func TestGitRepositoryCreateFile(t *testing.T) {
	repo := createTestRepo(t)

	err := repo.CreateFile("filename", "content")

	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(path.Join(repo.Dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}

// HELPERS

// createTestGitRepo creates a fully initialized Git repo including a master branch.
func createTestRepo(t *testing.T) GitRepository {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir)
	assert.Nil(t, err, "cannot initialize Git repow")
	output, err := repo.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	assert.Nilf(t, err, "cannot create initial commit: %s", output)
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
