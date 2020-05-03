package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/test"
	"github.com/stretchr/testify/assert"
)

func TestShellInDir_Run_workingDir(t *testing.T) {
	workDir := test.CreateTempDir(t)
	runner := command.ShellInDir{workDir}
	res, err := runner.Run("pwd")
	assert.Nil(t, err)
	assert.Equal(t, workDir, res.OutputSanitized())
}

func TestShellInDir_Run_arguments(t *testing.T) {
	runner := command.ShellInDir{test.CreateTempDir(t)}
	res, err := runner.Run("echo", "hello", "world")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestShellRunner_RunMany(t *testing.T) {
	workDir := test.CreateTempDir(t)
	runner := command.ShellInDir{workDir}
	err := runner.RunMany([][]string{
		{"touch", "first"},
		{"touch", "second"},
	})
	assert.Nil(t, err)
	infos, err := ioutil.ReadDir(workDir)
	assert.Nil(t, err)
	assert.Len(t, infos, 2)
	assert.Equal(t, "first", infos[0].Name())
	assert.Equal(t, "second", infos[1].Name())
}

func TestShellRunner_RunString(t *testing.T) {
	workDir := test.CreateTempDir(t)
	runner := command.ShellInDir{workDir}
	_, err := runner.RunString("touch first")
	assert.Nil(t, err)
	_, err = os.Stat(filepath.Join(workDir, "first"))
	assert.False(t, os.IsNotExist(err))
}
