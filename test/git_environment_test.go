package test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneGitEnvironment(t *testing.T) {
	dir := createTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(path.Join(dir, "memoized"))
	assert.Nil(t, err, "cannot create memoized GitEnvironment")
	_, err = CloneGitEnvironment(memoizedGitEnv, path.Join(dir, "cloned"))
	assert.Nil(t, err, "cannot clone GitEnvironment")
	assertIsNormalGitRepo(t, path.Join(dir, "cloned", "origin"))
	assertIsNormalGitRepo(t, path.Join(dir, "cloned", "developer"))
	assertHasGitBranch(t, path.Join(dir, "cloned", "developer"), "main")
}

func TestNewStandardGitEnvironment(t *testing.T) {
	gitEnvRootDir := createTempDir(t)
	result, err := NewStandardGitEnvironment(gitEnvRootDir)
	assert.Nil(t, err)
	// verify the origin repo
	assertIsNormalGitRepo(t, path.Join(gitEnvRootDir, "origin"))
	branch, err := result.OriginRepo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch, "the origin should be at the master branch so that we can push to it")
	// verify the developer repo
	assertIsNormalGitRepo(t, path.Join(gitEnvRootDir, "developer"))
	assertHasGlobalGitConfiguration(t, gitEnvRootDir)
	branch, err = result.DeveloperRepo.CurrentBranch()
	assert.Nil(t, err)
	assert.Equal(t, "main", branch)
}
