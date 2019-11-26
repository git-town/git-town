package command_test

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
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

func TestCommand_Run_UnknownExecutable(t *testing.T) {
	_, err := command.Run("zonk")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, exec.ErrNotFound))
}

func TestCommand_Run_ExitCode(t *testing.T) {
	_, err := command.Run("bash", "-c", "exit 2")
	var execError *exec.ExitError
	assert.True(t, errors.As(err, &execError))
	assert.Equal(t, 2, execError.ExitCode())
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

func TestCommand_Result_OutputContainsText(t *testing.T) {
	res, err := command.Run("echo", "hello world how are you?")
	assert.Nil(t, err)
	assert.True(t, res.OutputContainsText("world"), "should contain 'world'")
	assert.False(t, res.OutputContainsText("zonk"), "should not contain 'zonk'")
}

func TestCommand_Result_OutputContainsLine(t *testing.T) {
	res, err := command.Run("echo", "hello world")
	assert.Nil(t, err)
	assert.True(t, res.OutputContainsLine("hello world"), `should contain "hello world"`)
	assert.False(t, res.OutputContainsLine("hello"), `partial match should return false`)
	assert.False(t, res.OutputContainsLine("zonk"), `should not contain "zonk"`)
}
