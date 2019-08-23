package infra

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestGitEnvironmentCreateScenarioSetup(t *testing.T) {
	gitEnvRootDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("cannot create TempDir", err)
	}
	gitEnv, err := NewGitEnvironment(gitEnvRootDir)
	if err != nil {
		t.Error("cannot create new GitEnvironment", err)
	}
	err = gitEnv.Populate()
	if err != nil {
		t.Error("cannot populate GitEnvironment", err)
	}
	verifyIsBareGitRepo(path.Join(gitEnvRootDir, "origin"))

	// verify the new GitEnvironment has a "developer" folder
	devDir := path.Join(gitEnvRootDir, "developer")
	verifyFolderExists(devDir)

	// verify the "developer" folder contains a Git repo with a main branch
	verifyFolderExists(path.Join(devDir, ".git"))
	runner := ShellRunner{}
	err = os.Chdir(devDir)
	if err != nil {
		log.Fatal("cannot enter developer dir of GitEnvironment", err)
	}
	runResult := runner.Run("git", "branch")
	if runResult.Err != nil {
		log.Fatalf("cannot run 'git branch' in '%s': %s", devDir, runResult.Err)
	}
	dmp := diffmatchpatch.New()
	expected := "* main"
	diffs := dmp.DiffMain(strings.TrimSpace(expected), strings.TrimSpace(runResult.Output), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		log.Fatalf("folder '%s' has the wrong Git branches", gitEnvRootDir)
	}
}

func TestGitEnvironmentCloneEnvironment(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("cannot create temp dir", err)
	}
	memoizedGitEnv, err := NewGitEnvironment(path.Join(dir, "memoized"))
	if err != nil {
		t.Error("cannot create memoized GitEnvironment", err)
	}
	err = memoizedGitEnv.Populate()
	if err != nil {
		t.Error("cannot populate memoized GitEnvironment", err)
	}
	_, err = CloneGitEnvironment(memoizedGitEnv, path.Join(dir, "cloned"))
	if err != nil {
		log.Fatalf("cannot clone GitEnvironment: %s", err)
	}

	// verify that the GitEnvironment was properly cloned
	verifyIsBareGitRepo(path.Join(dir, "cloned", "origin"))
	devDir := path.Join(dir, "cloned", "developer")
	verifyFolderExists(devDir)
	verifyFolderExists(path.Join(dir, "cloned", "developer", ".git"))
	verifyHasGitBranches(devDir, "* main")
}

func verifyFolderExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("directory (%s) not found", dir)
	}
}

func verifyHasGitBranches(dir, expectedBranches string) {
	runner := ShellRunner{}
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	runResult := runner.Run("git", "branch")
	if runResult.Err != nil {
		log.Fatalf("cannot run 'git status' in '%s': %s", dir, runResult.Err)
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(strings.TrimSpace(expectedBranches), strings.TrimSpace(runResult.Output), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		log.Fatalf("folder '%s' has the wrong Git branches", dir)
	}
}

func verifyIsBareGitRepo(dir string) {
	verifyFolderExists(dir)
	runner := ShellRunner{}
	runResult := runner.Run("/bin/ls", "-1", dir)
	if runResult.Err != nil {
		log.Fatalf("command failed: %s", runResult.Err)
	}
	dmp := diffmatchpatch.New()
	expected := `
branches
config
description
HEAD
hooks
info
objects
refs`
	diffs := dmp.DiffMain(strings.TrimSpace(expected), strings.TrimSpace(runResult.Output), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		log.Fatalf("folder '%s' is not a bare Git repo", dir)
	}
}
