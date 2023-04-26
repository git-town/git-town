package subshell_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/test/ostools"
	"github.com/git-town/git-town/v8/test/subshell"
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
		runner := subshell.Mocking{
			WorkingDir: devDir,
			HomeDir:    workDir,
			BinDir:     filepath.Join(workDir, "bin"),
		}
		err = runner.MockCommand("foo")
		assert.NoError(t, err)
		// run a program that calls the mocked command
		res, err := runner.Run("bash", "-c", "foo bar")
		assert.NoError(t, err)
		// verify that it called our overridden "foo" command
		assert.Equal(t, "foo called with: bar", res)
	})

	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		runner := subshell.Mocking{
			WorkingDir: t.TempDir(),
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		res, err := runner.Run("echo", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", res)
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		workDir := t.TempDir()
		runner := subshell.Mocking{
			WorkingDir: workDir,
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		err := runner.RunMany([][]string{
			{"touch", "first"},
			{"touch", "second"},
		})
		assert.NoError(t, err)
		entries, err := os.ReadDir(workDir)
		assert.NoError(t, err)
		assert.Len(t, entries, 2)
		assert.Equal(t, "first", entries[0].Name())
		assert.Equal(t, "second", entries[1].Name())
	})

	t.Run(".RunString()", func(t *testing.T) {
		t.Parallel()
		workDir := t.TempDir()
		runner := subshell.Mocking{
			WorkingDir: workDir,
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		_, err := runner.RunString("touch first")
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(workDir, "first"))
		assert.False(t, os.IsNotExist(err))
	})

	t.Run(".RunStringWith", func(t *testing.T) {
		t.Run("without input", func(t *testing.T) {
			t.Parallel()
			dir1 := t.TempDir()
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			assert.NoError(t, err)
			r := subshell.Mocking{
				WorkingDir: dir1,
				HomeDir:    t.TempDir(),
				BinDir:     "",
			}
			toolPath := filepath.Join(dir2, "list-dir")
			err = ostools.CreateLsTool(toolPath)
			assert.NoError(t, err)
			res, err := r.RunWith(&subshell.Options{Dir: "subdir"}, toolPath)
			assert.NoError(t, err)
			assert.Equal(t, ostools.ScriptName("list-dir"), res)
		})

		t.Run("with input", func(t *testing.T) {
			t.Parallel()
			dir1 := t.TempDir()
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			assert.NoError(t, err)
			r := subshell.Mocking{
				WorkingDir: dir1,
				HomeDir:    t.TempDir(),
				BinDir:     "",
			}
			toolPath := filepath.Join(dir2, "list-dir")
			err = ostools.CreateInputTool(toolPath)
			assert.NoError(t, err)
			cmd, args := ostools.CallScriptArgs(toolPath)
			res, err := r.RunWith(&subshell.Options{Input: []string{"one\n", "two\n"}}, cmd, args...)
			assert.NoError(t, err)
			assert.Contains(t, res, "You entered one and two")
		})
	})
}
