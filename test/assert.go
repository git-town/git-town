// This file defines general assertions for testing.
// Assertions take a *testing.T object and fail the test using it.

package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func assertFileExists(t *testing.T, dir, filename string) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	info, err := os.Stat(filePath)
	assert.Nilf(t, err, "file %q does not exist", filePath)
	assert.Falsef(t, info.IsDir(), "%q is a directory", filePath)
}

func assertFileExistsWithContent(t *testing.T, dir, filename, expectedContent string) {
	t.Helper()
	fileContent, err := ioutil.ReadFile(filepath.Join(dir, filename))
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, string(fileContent))
}

func assertFolderExists(t *testing.T, dir string) {
	t.Helper()
	_, err := os.Stat(dir)
	assert.Falsef(t, os.IsNotExist(err), "directory %q not found", dir)
}

func assertHasGitBranch(t *testing.T, dir, expectedBranch string) {
	t.Helper()
	outcome, err := run.InDir(dir, "git", "branch")
	assert.Nilf(t, err, "cannot run 'git status' in %q", dir)
	assert.Contains(t, outcome.OutputSanitized(), expectedBranch, "doesn't have Git branch")
}

func assertHasGlobalGitConfiguration(t *testing.T, dir string) {
	t.Helper()
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	for e := range entries {
		if entries[e].Name() == ".gitconfig" {
			return
		}
	}
	t.Fatalf(".gitconfig not found in %q", dir)
}

func assertIsNormalGitRepo(t *testing.T, dir string) {
	t.Helper()
	assertFolderExists(t, dir)
	entries, err := ioutil.ReadDir(dir)
	assert.Nilf(t, err, "cannot list directory %q", dir)
	assert.Equal(t, ".git", entries[0].Name())
}
