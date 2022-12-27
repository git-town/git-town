//nolint:testpackage
package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCommandsInGitTownOutput(t *testing.T) {
	t.Parallel()
	tests := map[string][]ExecutedGitCommand{
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
		assert.Equal(t, expected, GitCommandsInGitTownOutput(input), input)
	}
}
