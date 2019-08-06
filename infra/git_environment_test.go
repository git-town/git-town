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
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	ge, err := NewGitEnvironment(dir)
	if err != nil {
		t.Error(err)
	}

	err = ge.Populate()
	if err != nil {
		t.Error(err)
	}

	verifyIsBareGitRepo(path.Join(dir, "origin"))

	// verify the "developer" folder exists
	devDir := path.Join(dir, "developer")
	verifyFolderExists(devDir)

	// verify the "developer" folder contains a Git repo with a main branch
	verifyFolderExists(path.Join(dir, "origin", "developer", ".git"))
	runner := ShellRunner{}
	err = os.Chdir(devDir)
	if err != nil {
		log.Fatal(err)
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
		log.Fatalf("folder '%s' has the wrong Git branches", dir)
	}
}

func TestGitEnvironmentCloneEnvironment(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	ge, err := NewGitEnvironment(path.Join(dir, "memoized"))
	if err != nil {
		t.Error(err)
	}
	err = ge.Populate()
	if err != nil {
		t.Error(err)
	}

	ce, err := CloneGitEnvironment(ge, path.Join(dir, "cloned"))
	if err != nil {
		log.Fatalf("cannot clone environment: %s", err)
	}
	if ce == nil {
		log.Fatal("returned nil for cloned environment")
	}

	verifyIsBareGitRepo(path.Join(dir, "cloned", "origin"))

	// verify the "developer" folder exists
	devDir := path.Join(dir, "cloned", "developer")
	verifyFolderExists(devDir)

	// verify the "developer" folder contains a Git repo
	verifyFolderExists(path.Join(dir, "cloned", "developer", ".git"))
	runner := ShellRunner{}
	err = os.Chdir(devDir)
	if err != nil {
		log.Fatal(err)
	}
	runResult := runner.Run("git", "branch")
	if runResult.Err != nil {
		log.Fatalf("cannot run 'git status' in '%s': %s", devDir, runResult.Err)
	}
	dmp := diffmatchpatch.New()
	expected := "* main"
	diffs := dmp.DiffMain(strings.TrimSpace(expected), strings.TrimSpace(runResult.Output), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		log.Fatalf("folder '%s' has the wrong Git branches", dir)
	}
}

func verifyFolderExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("directory (%s) not found", dir)
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
