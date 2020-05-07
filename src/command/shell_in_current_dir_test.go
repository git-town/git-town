package command_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/git-town/git-town/src/command"
	"github.com/stretchr/testify/assert"
)

func TestShellInCurrentDir_MustRun(t *testing.T) {
	runner := command.ShellInCurrentDir{}
	res := runner.MustRun("echo", "hello", "world")
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestShellInCurrentDir_Run_arguments(t *testing.T) {
	runner := command.ShellInCurrentDir{}
	res, err := runner.Run("echo", "hello", "world")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestShellInCurrentDir_RunMany(t *testing.T) {
	runner := command.ShellInCurrentDir{}
	err := runner.RunMany([][]string{
		{"mkdir", "tmp"},
		{"touch", "tmp/first"},
		{"touch", "tmp/second"},
	})
	defer os.RemoveAll("tmp")
	assert.Nil(t, err)
	infos, err := ioutil.ReadDir("tmp")
	assert.Nil(t, err)
	assert.Equal(t, "first", infos[0].Name())
	assert.Equal(t, "second", infos[1].Name())
}

func TestShellInCurrentDir_RunString(t *testing.T) {
	runner := command.ShellInCurrentDir{}
	_, err := runner.RunString("touch first")
	defer os.Remove("first")
	assert.Nil(t, err)
	_, err = os.Stat("first")
	assert.False(t, os.IsNotExist(err))
}

func TestShellInCurrentDir_RunStringWith(t *testing.T) {
	runner := command.ShellInCurrentDir{}
	res, err := runner.RunStringWith("ls -1", command.Options{Dir: ".."})
	assert.Nil(t, err)
	assert.Contains(t, res.OutputSanitized(), "cmd")
}
