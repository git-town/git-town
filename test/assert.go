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

func assertFileExists(t *testing.T, dir, filename string) {
	filePath := path.Join(dir, filename)
	info, err := os.Stat(filePath)
	assert.Nilf(t, err, "file %q does not exist", filePath)
	assert.Falsef(t, info.IsDir(), "%q is a directory", filePath)
}

func assertFileExistsWithContent(t *testing.T, dir, filename, expectedContent string) {
	fileContent, err := ioutil.ReadFile(path.Join(dir, filename))
	assert.Nil(t, err)
	assert.Equal(t, expectedContent, string(fileContent))
}

func assertFolderExists(t *testing.T, dir string) {
	_, err := os.Stat(dir)
	assert.Falsef(t, os.IsNotExist(err), "directory %q not found", dir)
}

func assertHasGitBranches(t *testing.T, dir, expectedBranches string) {
	runner := ShellRunner{}
	err := os.Chdir(dir)
	assert.Nil(t, err)
	runResult := runner.Run("git", "branch")
	assert.Nilf(t, runResult.Err, "cannot run 'git status' in %q", dir)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(strings.TrimSpace(expectedBranches), strings.TrimSpace(runResult.Output), false)
	if len(diffs) > 1 {
		fmt.Println(dmp.DiffPrettyText(diffs))
		log.Fatalf("folder %q has the wrong Git branches", dir)
	}
}

func assertIsBareGitRepo(t *testing.T, dir string) {
	assertFolderExists(t, dir)
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	testData := []string{"HEAD", "branches", "config", "description", "hooks", "info", "objects", "refs"}
	for i, expected := range testData {
		assert.Equal(t, expected, entries[i].Name())
	}
}
