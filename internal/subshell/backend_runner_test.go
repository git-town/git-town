package subshell_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/subshell"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBackendRunner(t *testing.T) {
	t.Parallel()

	t.Run("Query", func(t *testing.T) {
		t.Parallel()
		t.Run("happy path", func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
			output, err := runner.Query("echo", "hello", "world  ")
			must.NoError(t, err)
			must.EqOp(t, "hello world  \n", output)
		})

		t.Run("unknown executable", func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
			err := runner.Run("zonk")
			must.Error(t, err)
			var execError *exec.Error
			must.True(t, errors.As(err, &execError))
		})

		t.Run("non-zero exit code", func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
			err := runner.Run("bash", "-c", "echo hi && exit 2")
			expectedError := `
----------------------------------------
Diagnostic information of failed command

COMMAND: bash -c echo hi && exit 2
ERROR: exit status 2
OUTPUT START
hi

OUTPUT END
----------------------------------------`
			must.EqOp(t, expectedError, err.Error())
		})
	})

	t.Run("QueryTrim", func(t *testing.T) {
		t.Parallel()
		t.Run("trims whitespace", func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
			output, err := runner.QueryTrim("echo", "hello", "world  ")
			must.NoError(t, err)
			must.EqOp(t, "hello world", output)
		})
	})

	t.Run("ReplaceZeroWithNewlines", func(t *testing.T) {
		t.Parallel()

		t.Run("empty input", func(t *testing.T) {
			t.Parallel()
			give := []byte{}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{}
			must.SliceEqOp(t, want, have)
		})

		t.Run("no null bytes", func(t *testing.T) {
			t.Parallel()
			give := []byte("hello world")
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte("hello world")
			must.SliceEqOp(t, want, have)
		})

		t.Run("single null byte", func(t *testing.T) {
			t.Parallel()
			give := []byte{'h', 'e', 'l', 'l', 'o', 0x00, 'w', 'o', 'r', 'l', 'd'}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{'h', 'e', 'l', 'l', 'o', '\n', '\n', 'w', 'o', 'r', 'l', 'd'}
			must.SliceEqOp(t, want, have)
		})

		t.Run("multiple null bytes", func(t *testing.T) {
			t.Parallel()
			give := []byte{'a', 0x00, 'b', 0x00, 'c'}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{'a', '\n', '\n', 'b', '\n', '\n', 'c'}
			must.SliceEqOp(t, want, have)
		})

		t.Run("null byte at beginning", func(t *testing.T) {
			t.Parallel()
			give := []byte{0x00, 'h', 'e', 'l', 'l', 'o'}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{'\n', '\n', 'h', 'e', 'l', 'l', 'o'}
			must.SliceEqOp(t, want, have)
		})

		t.Run("null byte at end", func(t *testing.T) {
			t.Parallel()
			give := []byte{'h', 'e', 'l', 'l', 'o', 0x00}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{'h', 'e', 'l', 'l', 'o', '\n', '\n'}
			must.SliceEqOp(t, want, have)
		})

		t.Run("only null bytes", func(t *testing.T) {
			t.Parallel()
			give := []byte{0x00, 0x00, 0x00}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{'\n', '\n', '\n', '\n', '\n', '\n'}
			must.SliceEqOp(t, want, have)
		})

		t.Run("consecutive null bytes", func(t *testing.T) {
			t.Parallel()
			give := []byte{'a', 0x00, 0x00, 'b'}
			have := subshell.ReplaceZeroWithNewlines(give)
			want := []byte{'a', '\n', '\n', '\n', '\n', 'b'}
			must.SliceEqOp(t, want, have)
		})
	})
}
