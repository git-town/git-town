//nolint:testpackage
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitManager_CreateMemoizedEnvironment(t *testing.T) {
	t.Parallel()
	dir := CreateTempDir(t)
	_, err := NewGitManager(dir)
	assert.Nil(t, err, "creating memoized environment failed")
	memoizedPath := filepath.Join(dir, "memoized")
	_, err = os.Stat(memoizedPath)
	assert.Falsef(t, os.IsNotExist(err), "memoized directory %q not found", memoizedPath)
}

func TestGitManager_CreateScenarioEnvironment(t *testing.T) {
	t.Parallel()
	dir := CreateTempDir(t)
	gm, err := NewGitManager(dir)
	assert.Nil(t, err, "creating memoized environment failed")
	result, err := gm.CreateScenarioEnvironment("foo")
	assert.Nil(t, err, "cannot create scenario environment")
	_, err = os.Stat(result.DevRepo.shell.workingDir)
	assert.False(t, os.IsNotExist(err), "scenario environment directory %q not found", result.DevRepo.WorkingDir)
}
