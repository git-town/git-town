package output_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/output"
	"github.com/stretchr/testify/assert"
)

func TestDebugCommandsInGitTownOutput(t *testing.T) {
	t.Parallel()
	t.Run("single line", func(t *testing.T) {
		give := "(debug) foo bar"
		want := []string{"foo bar"}
		have := output.DebugCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("multiple lines", func(t *testing.T) {
		give := "(debug) command one\n\n(debug) command two\n\n"
		want := []string{"command one", "command two"}
		have := output.DebugCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
	t.Run("no debug command", func(t *testing.T) {
		give := "command one"
		want := []string{}
		have := output.DebugCommandsInGitTownOutput(give)
		assert.Equal(t, want, have)
	})
}
