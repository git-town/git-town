package test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitRepository(t *testing.T) {
	_ = NewGitRepository(createTempDir(t))
}

func TestInitBareGitRepository(t *testing.T) {
	result, err := InitGitRepository(createTempDir(t), true)
	assert.Nil(t, err, "cannot initialize bare GitRepository")
	assertIsBareGitRepo(t, result.Dir)
}

func TestInitNormalGitRepository(t *testing.T) {
	result, err := InitGitRepository(createTempDir(t), false)
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, result.Dir)
}

func TestCloneGitRepository(t *testing.T) {
	rootDir := createTempDir(t)
	originPath := path.Join(rootDir, "origin")
	_, err := InitGitRepository(originPath, true)
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := path.Join(rootDir, "cloned")

	_, err = CloneGitRepository(originPath, clonedPath)

	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestGitRepositoryCreateFile(t *testing.T) {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, false)
	assert.Nil(t, err, "cannot initialize Git repo")

	err = repo.createFile("filename", "content")

	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(path.Join(dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}
