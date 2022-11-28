package run_test

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()
	t.Run(".Exec()", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			t.Parallel()
			res, err := run.Exec("echo", "foo")
			assert.NoError(t, err)
			assert.Equal(t, "foo\n", res.Output())
		})

		t.Run("unknown executable", func(t *testing.T) {
			t.Parallel()
			_, err := run.Exec("zonk")
			assert.Error(t, err)
			var execError *exec.Error
			assert.True(t, errors.As(err, &execError))
		})

		t.Run("non-zero exit code", func(t *testing.T) {
			t.Parallel()
			result, err := run.Exec("bash", "-c", "echo hi && exit 2")
			assert.Equal(t, 2, result.ExitCode())
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

	t.Run(".InDir()", func(t *testing.T) {
		t.Parallel()
		dir, err := ioutil.TempDir("", "")
		assert.NoError(t, err)
		dirPath := filepath.Join(dir, "mydir")
		err = os.Mkdir(dirPath, 0o700)
		assert.NoError(t, err)
		err = ioutil.WriteFile(filepath.Join(dirPath, "one"), []byte{}, 0o500)
		assert.NoError(t, err)
		res, err := run.InDir(dirPath, "ls", "-1")
		assert.NoError(t, err)
		assert.Equal(t, "one", res.OutputSanitized())
	})

	t.Run("exec result", func(t *testing.T) {
		t.Run(".OutputContainsText()", func(t *testing.T) {
			t.Parallel()
			res, err := run.Exec("echo", "hello world how are you?")
			assert.NoError(t, err)
			assert.True(t, res.OutputContainsText("world"), "should contain 'world'")
			assert.False(t, res.OutputContainsText("zonk"), "should not contain 'zonk'")
		})

		t.Run(".OutputContainsLine()", func(t *testing.T) {
			t.Parallel()
			res, err := run.Exec("echo", "hello world")
			assert.NoError(t, err)
			assert.True(t, res.OutputContainsLine("hello world"), `should contain "hello world"`)
			assert.False(t, res.OutputContainsLine("hello"), `partial match should return false`)
			assert.False(t, res.OutputContainsLine("zonk"), `should not contain "zonk"`)
		})
	})
}
