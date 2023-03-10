//nolint:testpackage
package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	t.Parallel()
	t.Run("NewRepo", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		workingDir := filepath.Join(dir, "working")
		homeDir := filepath.Join(dir, "home")
		binDir := filepath.Join(dir, "bin")
		repo := NewRepo(workingDir, homeDir, binDir)
		assert.Equal(t, workingDir, repo.runner.workingDir)
		assert.Equal(t, homeDir, repo.runner.homeDir)
		assert.Equal(t, binDir, repo.runner.binDir)
	})

	t.Run(".Clone()", func(t *testing.T) {
		t.Parallel()
		origin := CreateRepo(t)
		clonedPath := filepath.Join(origin.runner.workingDir, "cloned")
		cloned, err := origin.Clone(clonedPath)
		assert.NoError(t, err)
		assert.Equal(t, clonedPath, cloned.runner.workingDir)
		assertIsNormalGitRepo(t, clonedPath)
	})
}
