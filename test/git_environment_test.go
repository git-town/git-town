package test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitEnvironmentPopulate(t *testing.T) {
	gitEnvRootDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	gitEnv, err := NewGitEnvironment(gitEnvRootDir)
	assert.Nil(t, err, "cannot create new GitEnvironment")
	err = gitEnv.Populate()
	assert.Nil(t, err, "cannot populate GitEnvironment")
	assertIsBareGitRepo(t, path.Join(gitEnvRootDir, "origin"))

	// verify the new GitEnvironment has a "developer" folder
	devDir := path.Join(gitEnvRootDir, "developer")
	assertFolderExists(t, devDir)

	// verify the "developer" folder contains a Git repo with a main branch
	assertFolderExists(t, path.Join(devDir, ".git"))
	runner := ShellRunner{}
	err = os.Chdir(devDir)
	assert.Nil(t, err, "cannot enter developer dir of GitEnvironment")
	output, err := runner.Run("git", "branch")
	assert.Nilf(t, err, "cannot run 'git branch' in %q", devDir)
	assert.Contains(t, output, "* main")
}

func TestGitEnvironmentCloneEnvironment(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create temp dir")
	memoizedGitEnv, err := NewGitEnvironment(path.Join(dir, "memoized"))
	assert.Nil(t, err, "cannot create memoized GitEnvironment")
	err = memoizedGitEnv.Populate()
	assert.Nil(t, err, "cannot populate memoized GitEnvironment")
	_, err = CloneGitEnvironment(memoizedGitEnv, path.Join(dir, "cloned"))
	assert.Nil(t, err, "cannot clone GitEnvironment")

	// verify that the GitEnvironment was properly cloned
	assertIsBareGitRepo(t, path.Join(dir, "cloned", "origin"))
	devDir := path.Join(dir, "cloned", "developer")
	assertFolderExists(t, devDir)
	assertFolderExists(t, path.Join(dir, "cloned", "developer", ".git"))
	assertHasGitBranch(t, devDir, "* main")
}
