package output_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/output"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandsInGitTownOutput(t *testing.T) {
	t.Parallel()
	t.Run("single frontend line", func(t *testing.T) {
		t.Parallel()
		give := "\x1b[1m[mybranch] foo bar"
		want := []output.ExecutedGitCommand{
			{Command: "foo bar", Branch: "mybranch", CommandType: output.CommandTypeFrontend},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("multiple frontend lines", func(t *testing.T) {
		t.Parallel()
		give := "\x1b[1m[branch1] command one\n\n\x1b[1m[branch2] command two\n\n"
		want := []output.ExecutedGitCommand{
			{Command: "command one", Branch: "branch1", CommandType: output.CommandTypeFrontend},
			{Command: "command two", Branch: "branch2", CommandType: output.CommandTypeFrontend},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("frontend line without branch", func(t *testing.T) {
		t.Parallel()
		give := "\x1b[1mcommand one"
		want := []output.ExecutedGitCommand{
			{Command: "command one", Branch: "", CommandType: output.CommandTypeFrontend},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("single debug line", func(t *testing.T) {
		t.Parallel()
		give := "(debug) foo bar"
		want := []output.ExecutedGitCommand{
			{Command: "foo bar", CommandType: output.CommandTypeBackend, Branch: ""},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("multiple debug lines", func(t *testing.T) {
		t.Parallel()
		give := "(debug) command one\n\n(debug) command two\n\n"
		want := []output.ExecutedGitCommand{
			{Command: "command one", CommandType: output.CommandTypeBackend, Branch: ""},
			{Command: "command two", CommandType: output.CommandTypeBackend, Branch: ""},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("line withouth a command", func(t *testing.T) {
		t.Parallel()
		give := "hello world"
		want := []output.ExecutedGitCommand{}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
}
