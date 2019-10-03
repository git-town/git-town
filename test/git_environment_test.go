package test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStandardGitEnvironment(t *testing.T) {
	gitEnvRootDir := createTempDir(t)
	_, err := NewStandardGitEnvironment(gitEnvRootDir)

	assert.Nil(t, err, "cannot create new GitEnvironment")
	assertIsBareGitRepo(t, path.Join(gitEnvRootDir, "origin"))

	// verify the new GitEnvironment has a "developer" folder
	devDir := path.Join(gitEnvRootDir, "developer")
	assertFolderExists(t, devDir)

	// verify the "developer" folder contains a Git repo with a main branch
	assertFolderExists(t, path.Join(devDir, ".git"))
	runner := NewShellRunner(devDir)
	output, err := runner.Run("git", "branch")
	assert.Nilf(t, err, "cannot run 'git branch' in %q", devDir)
	assert.Contains(t, output, "* main")
}

func TestGitEnvironmentCloneEnvironment(t *testing.T) {
	dir := createTempDir(t)
	memoizedGitEnv, err := NewStandardGitEnvironment(path.Join(dir, "memoized"))
	assert.Nil(t, err, "cannot create memoized GitEnvironment")

	_, err = CloneGitEnvironment(memoizedGitEnv, path.Join(dir, "cloned"))

	assert.Nil(t, err, "cannot clone GitEnvironment")
	assertIsBareGitRepo(t, path.Join(dir, "cloned", "origin"))
	devDir := path.Join(dir, "cloned", "developer")
	assertFolderExists(t, devDir)
	assertFolderExists(t, path.Join(dir, "cloned", "developer", ".git"))
	assertHasGitBranch(t, devDir, "* main")
}
