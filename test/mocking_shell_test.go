package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/src/command"
	"github.com/stretchr/testify/assert"
)

func TestMockingShell_MockCommand(t *testing.T) {
	workDir := CreateTempDir(t)
	devDir := filepath.Join(workDir, "dev")
	err := os.Mkdir(devDir, 0744)
	assert.Nil(t, err)
	shell := NewMockingShell(devDir, workDir, filepath.Join(workDir, "bin"))
	err = shell.MockCommand("foo")
	assert.Nil(t, err)
	// run a program that calls the mocked command
	res, err := shell.Run("bash", "-c", "foo bar")
	assert.Nil(t, err)
	// verify that it called our overridden "foo" command
	assert.Equal(t, "foo called with: bar", res.OutputSanitized())
}

func TestShellRunner_Run(t *testing.T) {
	runner := NewMockingShell(CreateTempDir(t), CreateTempDir(t), "")
	res, err := runner.Run("echo", "hello", "world")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestShellRunner_RunMany(t *testing.T) {
	workDir := CreateTempDir(t)
	runner := NewMockingShell(workDir, CreateTempDir(t), "")
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
	workDir := CreateTempDir(t)
	runner := NewMockingShell(workDir, CreateTempDir(t), "")
	_, err := runner.RunString("touch first")
	assert.Nil(t, err)
	_, err = os.Stat(filepath.Join(workDir, "first"))
	assert.False(t, os.IsNotExist(err))
}

func TestShellRunner_RunStringWith_Dir(t *testing.T) {
	dir1 := CreateTempDir(t)
	dir2 := filepath.Join(dir1, "subdir")
	err := os.Mkdir(dir2, 0744)
	assert.Nil(t, err)
	runner := NewMockingShell(dir1, CreateTempDir(t), "")
	toolPath := filepath.Join(dir2, "list-dir")
	err = CreateLsTool(toolPath)
	assert.Nil(t, err)
	res, err := runner.RunWith(command.Options{Dir: "subdir"}, toolPath)
	assert.Nil(t, err)
	assert.Equal(t, ScriptName("list-dir"), res.OutputSanitized())
}

func TestShellRunner_RunStringWith_Input(t *testing.T) {
	dir1 := CreateTempDir(t)
	dir2 := filepath.Join(dir1, "subdir")
	err := os.Mkdir(dir2, 0744)
	assert.Nil(t, err)
	runner := NewMockingShell(dir1, CreateTempDir(t), "")
	toolPath := filepath.Join(dir2, "list-dir")
	err = CreateInputTool(toolPath)
	assert.Nil(t, err)
	cmd, args := CallScriptArgs(toolPath)
	res, err := runner.RunWith(command.Options{Input: []string{"one\n", "two\n"}}, cmd, args...)
	assert.Nil(t, err)
	assert.Contains(t, res.OutputSanitized(), "You entered one and two")
}
