package output_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/output"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandsInGitTownOutput(t *testing.T) {
	t.Parallel()
	t.Run("single line", func(t *testing.T) {
		give := "\x1b[1m[mybranch] foo bar"
		want := []output.ExecutedGitCommand{
			{Command: "foo bar", Branch: "mybranch"},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("multiple lines", func(t *testing.T) {
		give := "\x1b[1m[branch1] command one\n\n\x1b[1m[branch2] command two\n\n"
		want := []output.ExecutedGitCommand{
			{Command: "command one", Branch: "branch1"},
			{Command: "command two", Branch: "branch2"},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("no branch", func(t *testing.T) {
		give := "\x1b[1mcommand one"
		want := []output.ExecutedGitCommand{
			{Command: "command one", Branch: ""},
		}
		have := output.GitCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
}
