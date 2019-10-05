// This file defines general assertions for testing.
// Assertions take a *testing.T object and fail the test using it.

package test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

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

func assertHasGitBranch(t *testing.T, dir, expectedBranch string) {
	runner := NewShellRunner(dir)
	output, err := runner.Run("git", "branch")
	assert.Nilf(t, err, "cannot run 'git status' in %q", dir)
	assert.Contains(t, output, expectedBranch, "doesn't have Git branch")
}

func assertIsBareGitRepo(t *testing.T, dir string) {
	assertFolderExists(t, dir)
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	entryNames := make([]string, len(entries))
	for i := range entries {
		entryNames[i] = entries[i].Name()
	}
	assert.Contains(t, entryNames, "HEAD")
	assert.Contains(t, entryNames, "config")
	assert.Contains(t, entryNames, "description")
	assert.Contains(t, entryNames, "hooks")
	assert.Contains(t, entryNames, "objects")
	assert.Contains(t, entryNames, "refs")
}

func assertIsNormalGitRepo(t *testing.T, dir string) {
	assertFolderExists(t, dir)
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	testData := []string{".git"}
	for i, expected := range testData {
		assert.Equal(t, expected, entries[i].Name())
	}
}
