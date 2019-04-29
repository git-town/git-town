package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandsInOutput(t *testing.T) {
	actual := CommandsInOutput("\x1b[1m[mybranch] foo bar")
	assert.Equal(t, actual, []string{"foo bar"})
}

func TestCommandsInOutputMultiline(t *testing.T) {
	actual := CommandsInOutput("\x1b[1m[branch1] command one\n\n\x1b[1m[branch2] command two\n\n")
	assert.Equal(t, actual, []string{"command one", "command two"})
}

func TestCommandsInOutputNoBranch(t *testing.T) {
	actual := CommandsInOutput("\x1b[1mcommand one")
	assert.Equal(t, actual, []string{"command one"})
}
