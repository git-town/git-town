package run_test

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func TestSilentRunner(t *testing.T) {
	t.Parallel()
	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		t.Run("happy path", func(t *testing.T) {
			debug := false
			runner := run.SilentRunner{Debug: &debug}
			res, err := runner.Run("echo", "hello", "world")
			assert.NoError(t, err)
			assert.Equal(t, "hello world", res.OutputSanitized())
		})
		t.Run("unknown executable", func(t *testing.T) {
			t.Parallel()
			debug := false
			runner := run.SilentRunner{Debug: &debug}
			_, err := runner.Run("zonk")
			assert.Error(t, err)
			var execError *exec.Error
			assert.True(t, errors.As(err, &execError))
		})
		t.Run("non-zero exit code", func(t *testing.T) {
			t.Parallel()
			debug := false
			runner := run.SilentRunner{Debug: &debug}
			result, err := runner.Run("bash", "-c", "echo hi && exit 2")
			assert.Equal(t, 2, result.ExitCode)
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
		t.Run("result contains text", func(t *testing.T) {
			t.Parallel()
			debug := false
			runner := run.SilentRunner{Debug: &debug}
			res, err := runner.Run("echo", "hello world how are you?")
			assert.NoError(t, err)
			assert.True(t, res.OutputContainsText("world"), "should contain 'world'")
			assert.False(t, res.OutputContainsText("zonk"), "should not contain 'zonk'")
		})
		t.Run(".OutputContainsLine()", func(t *testing.T) {
			t.Parallel()
			debug := false
			runner := run.SilentRunner{Debug: &debug}
			res, err := runner.Run("echo", "hello world")
			assert.NoError(t, err)
			assert.True(t, res.OutputContainsLine("hello world"), `should contain "hello world"`)
			assert.False(t, res.OutputContainsLine("hello"), `partial match should return false`)
			assert.False(t, res.OutputContainsLine("zonk"), `should not contain "zonk"`)
		})
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		debug := false
		runner := run.SilentRunner{Debug: &debug}
		err := runner.RunMany([][]string{
			{"mkdir", "tmp"},
			{"touch", "tmp/first"},
			{"touch", "tmp/second"},
		})
		defer os.RemoveAll("tmp")
		assert.NoError(t, err)
		entries, err := os.ReadDir("tmp")
		assert.NoError(t, err)
		assert.Equal(t, "first", entries[0].Name())
		assert.Equal(t, "second", entries[1].Name())
	})

	t.Run(".RunString()", func(t *testing.T) {
		t.Parallel()
		debug := false
		runner := run.SilentRunner{Debug: &debug}
		_, err := runner.RunString("touch first")
		defer os.Remove("first")
		assert.NoError(t, err)
		_, err = os.Stat("first")
		assert.False(t, os.IsNotExist(err))
	})
}
