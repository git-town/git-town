package subshell_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v9/test/ostools"
	"github.com/git-town/git-town/v9/test/subshell"
	"github.com/stretchr/testify/assert"
)

func TestMockingRunner(t *testing.T) {
	t.Parallel()
	t.Run(".MockCommand()", func(t *testing.T) {
		t.Parallel()
		workDir := t.TempDir()
		devDir := filepath.Join(workDir, "dev")
		err := os.Mkdir(devDir, 0o744)
		assert.NoError(t, err)
		runner := subshell.TestRunner{
			WorkingDir: devDir,
			HomeDir:    workDir,
			BinDir:     filepath.Join(workDir, "bin"),
		}
		runner.MockCommand("foo")
		// run a program that calls the mocked command
		res, err := runner.Query("bash", "-c", "foo bar")
		assert.NoError(t, err)
		// verify that it called our overridden "foo" command
		assert.Equal(t, "foo called with: bar", res)
	})

	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		runner := subshell.TestRunner{
			WorkingDir: t.TempDir(),
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		res, err := runner.Query("echo", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", res)
	})

	t.Run(".QueryString()", func(t *testing.T) {
		t.Parallel()
		workDir := t.TempDir()
		runner := subshell.TestRunner{
			WorkingDir: workDir,
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		_, err := runner.QueryString("touch first")
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(workDir, "first"))
		assert.False(t, os.IsNotExist(err))
	})

	t.Run(".QueryWith", func(t *testing.T) {
		t.Run("without input", func(t *testing.T) {
			t.Parallel()
			dir1 := t.TempDir()
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			assert.NoError(t, err)
			r := subshell.TestRunner{
				WorkingDir: dir1,
				HomeDir:    t.TempDir(),
				BinDir:     "",
			}
			toolPath := filepath.Join(dir2, "list-dir")
			err = ostools.CreateLsTool(toolPath)
			assert.NoError(t, err)
			res, err := r.QueryWith(&subshell.Options{Dir: "subdir"}, toolPath)
			assert.NoError(t, err)
			assert.Equal(t, ostools.ScriptName("list-dir"), res)
		})

		t.Run("with input", func(t *testing.T) {
			t.Parallel()
			dir1 := t.TempDir()
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			assert.NoError(t, err)
			r := subshell.TestRunner{
				WorkingDir: dir1,
				HomeDir:    t.TempDir(),
				BinDir:     "",
			}
			toolPath := filepath.Join(dir2, "list-dir")
			err = ostools.CreateInputTool(toolPath)
			assert.NoError(t, err)
			cmd, args := ostools.CallScriptArgs(toolPath)
			res, err := r.QueryWith(&subshell.Options{Input: []string{"one\n", "two\n"}}, cmd, args...)
			assert.NoError(t, err)
			assert.Contains(t, res, "You entered one and two")
		})
	})

	t.Run(".QueryWithCode", func(t *testing.T) {
		t.Parallel()
		t.Run("exit code 0", func(t *testing.T) {
			r := subshell.TestRunner{
				BinDir:     "",
				Debug:      false,
				HomeDir:    "",
				WorkingDir: "",
			}
			output, exitCode, err := r.QueryWithCode(&subshell.Options{}, "echo", "hello")
			assert.Equal(t, "hello", output)
			assert.Equal(t, 0, exitCode)
			assert.NoError(t, err)
		})
		t.Run("exit code 1", func(t *testing.T) {
			r := subshell.TestRunner{
				BinDir:     "",
				Debug:      false,
				HomeDir:    "",
				WorkingDir: "",
			}
			output, exitCode, err := r.QueryWithCode(&subshell.Options{}, "bash", "-c", "echo hello && exit 1")
			assert.Equal(t, "hello", output)
			assert.Equal(t, 1, exitCode)
			assert.NoError(t, err)
		})
	})
}
