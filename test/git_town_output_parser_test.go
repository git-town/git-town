package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCommandsInGitTownOutput(t *testing.T) {
	testData := map[string][]ExecutedGitCommand{
		// simple
		"\x1b[1m[mybranch] foo bar": []ExecutedGitCommand{
			{Command: "foo bar", Branch: "mybranch"}},

		// multiline
		"\x1b[1m[branch1] command one\n\n\x1b[1m[branch2] command two\n\n": []ExecutedGitCommand{
			{Command: "command one", Branch: "branch1"},
			{Command: "command two", Branch: "branch2"}},

		// no branch
		"\x1b[1mcommand one": []ExecutedGitCommand{
			{Command: "command one", Branch: ""}},
	}
	for input, expected := range testData {
		assert.Equal(t, expected, GitCommandsInGitTownOutput(input))
	}
}
