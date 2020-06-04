package command_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/git-town/git-town/src/command"
	"github.com/stretchr/testify/assert"
)

func TestSilentShell_MustRun(t *testing.T) {
	shell := command.SilentShell{}
	res := shell.MustRun("echo", "hello", "world")
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestSilentShell_Run_arguments(t *testing.T) {
	shell := command.SilentShell{}
	res, err := shell.Run("echo", "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestSilentShell_RunMany(t *testing.T) {
	shell := command.SilentShell{}
	err := shell.RunMany([][]string{
		{"mkdir", "tmp"},
		{"touch", "tmp/first"},
		{"touch", "tmp/second"},
	})
	defer os.RemoveAll("tmp")
	assert.NoError(t, err)
	infos, err := ioutil.ReadDir("tmp")
	assert.NoError(t, err)
	assert.Equal(t, "first", infos[0].Name())
	assert.Equal(t, "second", infos[1].Name())
}

func TestSilentShell_RunString(t *testing.T) {
	shell := command.SilentShell{}
	_, err := shell.RunString("touch first")
	defer os.Remove("first")
	assert.NoError(t, err)
	_, err = os.Stat("first")
	assert.False(t, os.IsNotExist(err))
}

func TestSilentShell_RunStringWith(t *testing.T) {
	shell := command.SilentShell{}
	res, err := shell.RunStringWith("ls -1", command.Options{Dir: ".."})
	assert.NoError(t, err)
	assert.Contains(t, res.OutputSanitized(), "cmd")
}
