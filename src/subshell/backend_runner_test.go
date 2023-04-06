package subshell_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/execute"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestBackendRunner(t *testing.T) {
	t.Parallel()
	t.Run("Run", func(t *testing.T) {
		t.Parallel()
		t.Run("happy path", func(t *testing.T) {
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: &tmpDir, Verbose: false, Stats: &execute.NoStatistics{}}
			output, err := runner.Run("echo", "hello", "world")
			assert.NoError(t, err)
			assert.Equal(t, "hello world", output)
		})

		t.Run("unknown executable", func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: &tmpDir, Verbose: false, Stats: &execute.NoStatistics{}}
			_, err := runner.Run("zonk")
			assert.Error(t, err)
			var execError *exec.Error
			assert.True(t, errors.As(err, &execError))
		})

		t.Run("non-zero exit code", func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: &tmpDir, Verbose: false, Stats: &execute.NoStatistics{}}
			_, err := runner.Run("bash", "-c", "echo hi && exit 2")
			expectedError := `
----------------------------------------
Diagnostic information of failed command

Command: bash -c echo hi && exit 2
Error: exit status 2
Output:
hi

----------------------------------------`
			assert.Equal(t, expectedError, err.Error())
		})
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: &tmpDir, Verbose: false, Stats: &execute.NoStatistics{}}
		err := runner.RunMany([][]string{
			{"mkdir", "tmp"},
			{"touch", "tmp/first"},
			{"touch", "tmp/second"},
		})
		assert.NoError(t, err)
		entries, err := os.ReadDir(filepath.Join(tmpDir, "tmp"))
		assert.NoError(t, err)
		assert.Equal(t, "first", entries[0].Name())
		assert.Equal(t, "second", entries[1].Name())
	})
}
