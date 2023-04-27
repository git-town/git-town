package commands_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/test/asserts"
	"github.com/git-town/git-town/v8/test/commands"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Parallel()
	t.Run("newRunner", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		workingDir := filepath.Join(dir, "working")
		homeDir := filepath.Join(dir, "home")
		binDir := filepath.Join(dir, "bin")
		runtime := commands.New(workingDir, homeDir, binDir)
		assert.Equal(t, workingDir, runtime.WorkingDir)
		assert.Equal(t, homeDir, runtime.HomeDir)
		assert.Equal(t, binDir, runtime.BinDir)
	})

	t.Run(".Clone()", func(t *testing.T) {
		t.Parallel()
		origin := commands.Create(t)
		clonedPath := filepath.Join(origin.Dir(), "cloned")
		cloned, err := commands.Clone(&origin, clonedPath)
		assert.NoError(t, err)
		assert.Equal(t, clonedPath, cloned.WorkingDir)
		asserts.IsGitRepo(t, clonedPath)
	})
}
