package output_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/output"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandsInGitTownOutput(t *testing.T) {
	t.Parallel()
	tests := map[string][]output.ExecutedGitCommand{
		// simple
		"\x1b[1m[mybranch] foo bar": {
			{Command: "foo bar", Branch: "mybranch"},
		},
		// multiline
		"\x1b[1m[branch1] command one\n\n\x1b[1m[branch2] command two\n\n": {
			{Command: "command one", Branch: "branch1"},
			{Command: "command two", Branch: "branch2"},
		},
		// no branch
		"\x1b[1mcommand one": {
			{Command: "command one", Branch: ""},
		},
	}
	for input, expected := range tests {
		assert.Equal(t, expected, output.GitCommandsInGitTownOutput(input), input)
	}
}
