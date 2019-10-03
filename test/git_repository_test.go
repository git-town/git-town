package test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	_ = NewGitRepository(dir)
}

func TestInitBareGitRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	result, err := InitGitRepository(dir, true)
	assert.Nil(t, err, "cannot initialize bare GitRepository")
	assertIsBareGitRepo(t, result.dir)
}

func TestInitNormalGitRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	result, err := InitGitRepository(dir, false)
	assert.Nil(t, err, "cannot initialize normal GitRepository")
	assertIsNormalGitRepo(t, result.dir)
}

func TestCloneGitRepository(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	originPath := path.Join(rootDir, "origin")
	_, err = InitGitRepository(originPath, true)
	assert.Nil(t, err, "cannot initialze origin Git repository")
	clonedPath := path.Join(rootDir, "cloned")

	_, err = CloneGitRepository(originPath, clonedPath)

	assert.Nil(t, err, "cannot clone repo")
	assertIsNormalGitRepo(t, clonedPath)
}

func TestGitRepositoryCreateFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	repo, err := InitGitRepository(dir, false)
	assert.Nil(t, err, "cannot initialize Git repo")

	err = repo.createFile("filename", "content")

	assert.Nil(t, err, "cannot create file in repo")
	content, err := ioutil.ReadFile(path.Join(dir, "filename"))
	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, "content", string(content))
}
