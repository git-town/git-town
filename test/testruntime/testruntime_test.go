package testruntime_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Parallel()
	t.Run("New", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		workingDir := filepath.Join(dir, "working")
		err := os.Mkdir(workingDir, 0x744)
		assert.NoError(t, err)
		homeDir := filepath.Join(dir, "home")
		binDir := filepath.Join(dir, "bin")
		runtime := testruntime.New(workingDir, homeDir, binDir)
		assert.Equal(t, workingDir, runtime.WorkingDir)
		assert.Equal(t, homeDir, runtime.HomeDir)
		assert.Equal(t, binDir, runtime.BinDir)
	})

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		clonedPath := filepath.Join(origin.WorkingDir, "cloned")
		cloned := testruntime.Clone(origin.TestRunner, clonedPath)
		assert.Equal(t, clonedPath, cloned.WorkingDir)
		asserts.IsGitRepo(t, clonedPath)
	})
}
