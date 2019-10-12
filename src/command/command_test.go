package command_test

import (
	"testing"

	"github.com/Originate/git-town/src/command"
	"github.com/stretchr/testify/assert"
)

func TestCommandRun(t *testing.T) {
	res := command.Run("echo", "foo")
	assert.Equal(t, "foo", res.Output())
}

func TestCommandOutputContainsText(t *testing.T) {
	res := command.Run("echo", "hello world how are you?")
	assert.True(t, res.OutputContainsText("world"), "should contain 'world'")
	assert.False(t, res.OutputContainsText("zonk"), "should not contain 'zonk'")
}

func TestCommandOutputContainsLine(t *testing.T) {
	res := command.Run("echo", "hello world")
	assert.True(t, res.OutputContainsLine("hello world"), `should contain "hello world"`)
	assert.False(t, res.OutputContainsLine("hello"), `partial match should return false`)
	assert.False(t, res.OutputContainsLine("zonk"), `should not contain "zonk"`)
}

func TestCommandErrUnknownExecutable(t *testing.T) {
	res := command.Run("zonk")
	assert.Error(t, res.Err())
}

func TestCommandErrExitCode(t *testing.T) {
	res := command.Run("bash", "-c", "exit 2")
	assert.Error(t, res.Err())
}
