package command_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Originate/git-town/src/command"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	res, err := command.Run("echo", "foo")
	assert.Nil(t, err)
	assert.Equal(t, "foo\n", res.Output())
}

func TestCommand_RunInDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	dirPath := path.Join(dir, "mydir")
	err = os.Mkdir(dirPath, 0744)
	assert.Nil(t, err)
	err = ioutil.WriteFile(path.Join(dirPath, "one"), []byte{}, 0744)
	assert.Nil(t, err)
	res, err := command.RunInDir(dirPath, "ls", "-1")
	assert.Nil(t, err)
	assert.Equal(t, "one", res.OutputSanitized())
}

func TestCommand_OutputContainsText(t *testing.T) {
	res, err := command.Run("echo", "hello world how are you?")
	assert.Nil(t, err)
	assert.True(t, res.OutputContainsText("world"), "should contain 'world'")
	assert.False(t, res.OutputContainsText("zonk"), "should not contain 'zonk'")
}

func TestCommand_OutputContainsLine(t *testing.T) {
	res, err := command.Run("echo", "hello world")
	assert.Nil(t, err)
	assert.True(t, res.OutputContainsLine("hello world"), `should contain "hello world"`)
	assert.False(t, res.OutputContainsLine("hello"), `partial match should return false`)
	assert.False(t, res.OutputContainsLine("zonk"), `should not contain "zonk"`)
}

func TestCommand_ErrUnknownExecutable(t *testing.T) {
	_, err := command.Run("zonk")
	assert.Error(t, err)
}

func TestCommand_ErrExitCode(t *testing.T) {
	res, err := command.Run("bash", "-c", "exit 2")
	assert.Nil(t, err)
	assert.Equal(t, 2, res.ExitCode())
}
