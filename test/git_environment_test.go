package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
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
	runResult := runner.Run("git", "branch")
	assert.Nilf(t, runResult.Err, "cannot run 'git branch' in %q", devDir)
	dmp := diffmatchpatch.New()
	expected := "* main"
	diffs := dmp.DiffMain(expected, strings.TrimSpace(runResult.Output), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		log.Fatalf("folder %q has the wrong Git branches", gitEnvRootDir)
	}
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
	assertHasGitBranches(t, devDir, "* main")
}
