package subshell_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v21/internal/test/ostools"
	"github.com/git-town/git-town/v21/internal/test/subshell"
	"github.com/shoenig/test/must"
)

func TestMockingRunner(t *testing.T) {
	t.Parallel()

	t.Run("MockCommand", func(t *testing.T) {
		t.Parallel()
		workDir := t.TempDir()
		devDir := filepath.Join(workDir, "dev")
		err := os.Mkdir(devDir, 0o744)
		must.NoError(t, err)
		runner := subshell.TestRunner{
			WorkingDir: devDir,
			HomeDir:    workDir,
			BinDir:     filepath.Join(workDir, "bin"),
		}
		runner.MockCommand("foo")
		// run a program that calls the mocked command
		res, err := runner.Query("bash", "-c", "foo bar")
		must.NoError(t, err)
		// verify that it called our overridden "foo" command
		must.EqOp(t, "foo called with: bar", res)
	})

	t.Run("MockCommitMessage", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		runner := subshell.TestRunner{
			WorkingDir: dir,
			HomeDir:    dir,
			BinDir:     filepath.Join(dir, "bin"),
		}
		runner.MockCommitMessage("test commit message")
		// Simulate Git calling the mock editor configured by MockCommitMessage and
		// verify its effect.
		// MockCommitMessage creates a custom editor that the runner makes available
		// via the GIT_EDITOR environment variable. We verify the following:
		// - GIT_EDITOR is available and executable.
		// - The contents of the file provided in the first argument to $GIT_EDITOR
		//   our expected commit message after the command has finished.
		_ = runner.MustQuery("bash", "-c", `"$GIT_EDITOR" output`)
		data, err := os.ReadFile(filepath.Join(dir, "output"))
		must.NoError(t, err)
		must.Eq(t, "test commit message\n", string(data))
	})

	t.Run("QueryString", func(t *testing.T) {
		t.Parallel()
		workDir := t.TempDir()
		runner := subshell.TestRunner{
			WorkingDir: workDir,
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		_, err := runner.QueryString("touch first")
		must.NoError(t, err)
		_, err = os.Stat(filepath.Join(workDir, "first"))
		must.False(t, os.IsNotExist(err))
	})

	t.Run("QueryWith", func(t *testing.T) {
		t.Run("without input", func(t *testing.T) {
			t.Parallel()
			dir1 := t.TempDir()
			dir2 := filepath.Join(dir1, "subdir")
			err := os.Mkdir(dir2, 0o744)
			must.NoError(t, err)
			r := subshell.TestRunner{
				WorkingDir: dir1,
				HomeDir:    t.TempDir(),
				BinDir:     "",
			}
			toolPath := filepath.Join(dir2, "list-dir")
			ostools.CreateLsTool(toolPath)
			res, err := r.QueryWith(&subshell.Options{Dir: "subdir"}, toolPath)
			must.NoError(t, err)
			must.EqOp(t, ostools.ScriptName("list-dir"), res)
		})
	})

	t.Run("QueryWithCode", func(t *testing.T) {
		t.Parallel()
		t.Run("exit code 0", func(t *testing.T) {
			t.Parallel()
			r := subshell.TestRunner{
				BinDir:     "",
				Verbose:    false,
				HomeDir:    "",
				WorkingDir: "",
			}
			output, exitCode, err := r.QueryWithCode(&subshell.Options{}, "echo", "hello")
			must.EqOp(t, "hello", output)
			must.EqOp(t, 0, exitCode)
			must.NoError(t, err)
		})
		t.Run("exit code 1", func(t *testing.T) {
			t.Parallel()
			r := subshell.TestRunner{
				BinDir:     "",
				Verbose:    false,
				HomeDir:    "",
				WorkingDir: "",
			}
			output, exitCode, err := r.QueryWithCode(&subshell.Options{}, "bash", "-c", "echo hello && exit 1")
			must.EqOp(t, "hello", output)
			must.EqOp(t, 1, exitCode)
			must.NoError(t, err)
		})
	})

	t.Run("Run", func(t *testing.T) {
		t.Parallel()
		runner := subshell.TestRunner{
			WorkingDir: t.TempDir(),
			HomeDir:    t.TempDir(),
			BinDir:     "",
		}
		res, err := runner.Query("echo", "hello", "world")
		must.NoError(t, err)
		must.EqOp(t, "hello world", res)
	})
}
