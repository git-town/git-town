package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCommandsInGitTownOutput(t *testing.T) {
	testData := map[string][]string{
		"\x1b[1m[mybranch] foo bar":                                        {"foo bar"},                    // simple
		"\x1b[1m[branch1] command one\n\n\x1b[1m[branch2] command two\n\n": {"command one", "command two"}, // multiline
		"\x1b[1mcommand one":                                               {"command one"},                // no branch
	}
	for input, expected := range testData {
		actual := GitCommandsInGitTownOutput(input)
		assert.Equal(t, expected, actual)
	}
}
