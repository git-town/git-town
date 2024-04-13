package testruntime_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/git-town/git-town/v14/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestRunner(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		workingDir := filepath.Join(dir, "working")
		err := os.Mkdir(workingDir, 0o744)
		must.NoError(t, err)
		homeDir := filepath.Join(dir, "home")
		binDir := filepath.Join(dir, "bin")
		runtime := testruntime.New(workingDir, homeDir, binDir)
		must.EqOp(t, workingDir, runtime.WorkingDir)
		must.EqOp(t, homeDir, runtime.HomeDir)
		must.EqOp(t, binDir, runtime.BinDir)
	})

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		clonedPath := filepath.Join(origin.WorkingDir, "cloned")
		cloned := testruntime.Clone(origin.TestRunner, clonedPath)
		must.EqOp(t, clonedPath, cloned.WorkingDir)
		asserts.IsGitRepo(t, clonedPath)
	})
}
