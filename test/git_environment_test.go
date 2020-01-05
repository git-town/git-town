package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitEnvironment(t *testing.T) {
	dir := createTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(filepath.Join(dir, "memoized"))
	assert.Nil(t, err, "cannot create memoized GitEnvironment")

	_, err = CloneGitEnvironment(memoizedGitEnv, filepath.Join(dir, "cloned"))

	assert.Nil(t, err, "cannot clone GitEnvironment")
	assertIsNormalGitRepo(t, filepath.Join(dir, "cloned", "origin"))
	assertIsNormalGitRepo(t, filepath.Join(dir, "cloned", "developer"))
	assertHasGitBranch(t, filepath.Join(dir, "cloned", "developer"), "main")
}

func TestNewStandardGitEnvironment(t *testing.T) {
	gitEnvRootDir := createTempDir(t)

	result, err := NewStandardGitEnvironment(gitEnvRootDir)

	assert.Nil(t, err)

	// verify the origin repo
	assertIsNormalGitRepo(t, filepath.Join(gitEnvRootDir, "origin"))
	branch, err := result.OriginRepo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch, "the origin should be at the master branch so that we can push to it")

	// verify the developer repo
	assertIsNormalGitRepo(t, filepath.Join(gitEnvRootDir, "developer"))
	assertHasGlobalGitConfiguration(t, gitEnvRootDir)
	branch, err = result.DeveloperRepo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "main", branch)
}
