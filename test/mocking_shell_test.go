//nolint:testpackage
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func TestMockingShell(t *testing.T) {
	t.Parallel()
	t.Run(".MockCommand()", func(t *testing.T) {
		t.Parallel()
		workDir := CreateTempDir(t)
		devDir := filepath.Join(workDir, "dev")
		err := os.Mkdir(devDir, 0o744)
		assert.NoError(t, err)
		shell := NewMockingShell(devDir, workDir, filepath.Join(workDir, "bin"))
		err = shell.MockCommand("foo")
		assert.NoError(t, err)
		// run a program that calls the mocked command
		res, err := shell.Run("bash", "-c", "foo bar")
		assert.NoError(t, err)
		// verify that it called our overridden "foo" command
		assert.Equal(t, "foo called with: bar", res.OutputSanitized())
	})

	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		runner := NewMockingShell(CreateTempDir(t), CreateTempDir(t), "")
		res, err := runner.Run("echo", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", res.OutputSanitized())
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		workDir := CreateTempDir(t)
		runner := NewMockingShell(workDir, CreateTempDir(t), "")
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
		workDir := CreateTempDir(t)
		runner := NewMockingShell(workDir, CreateTempDir(t), "")
		_, err := runner.RunString("touch first")
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(workDir, "first"))
		assert.False(t, os.IsNotExist(err))
	})

	t.Run(".RunStringWith()", func(t *testing.T) {
		t.Run("without input", func(t *testing.T) {
			t.Parallel()
			dir1 := CreateTempDir(t)
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			assert.NoError(t, err)
			runner := NewMockingShell(dir1, CreateTempDir(t), "")
			toolPath := filepath.Join(dir2, "list-dir")
			err = CreateLsTool(toolPath)
			assert.NoError(t, err)
			res, err := runner.RunWith(&run.Options{Dir: "subdir"}, toolPath)
			assert.NoError(t, err)
			assert.Equal(t, ScriptName("list-dir"), res.OutputSanitized())
		})

		t.Run("with input", func(t *testing.T) {
			t.Parallel()
			dir1 := CreateTempDir(t)
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			assert.NoError(t, err)
			runner := NewMockingShell(dir1, CreateTempDir(t), "")
			toolPath := filepath.Join(dir2, "list-dir")
			err = CreateInputTool(toolPath)
			assert.NoError(t, err)
			cmd, args := CallScriptArgs(toolPath)
			res, err := runner.RunWith(&run.Options{Input: []string{"one\n", "two\n"}}, cmd, args...)
			assert.NoError(t, err)
			assert.Contains(t, res.OutputSanitized(), "You entered one and two")
		})
	})
}
