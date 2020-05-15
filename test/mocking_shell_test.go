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
	workDir := createTempDir(t)
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
	runner := NewMockingShell(createTempDir(t), createTempDir(t), "")
	res, err := runner.Run("echo", "hello", "world")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", res.OutputSanitized())
}

func TestShellRunner_RunMany(t *testing.T) {
	workDir := createTempDir(t)
	runner := NewMockingShell(workDir, createTempDir(t), "")
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
	workDir := createTempDir(t)
	runner := NewMockingShell(workDir, createTempDir(t), "")
	_, err := runner.RunString("touch first")
	assert.Nil(t, err)
	_, err = os.Stat(filepath.Join(workDir, "first"))
	assert.False(t, os.IsNotExist(err))
}

func TestShellRunner_RunStringWith_Dir(t *testing.T) {
	dir1 := createTempDir(t)
	dir2 := filepath.Join(dir1, "subdir")
	err := os.Mkdir(dir2, 0744)
	assert.Nil(t, err)
	runner := NewMockingShell(dir1, createTempDir(t), "")
	toolPath := filepath.Join(dir2, "list-dir")
	err = ioutil.WriteFile(toolPath, []byte("#!/usr/bin/env bash\n\nls\n"), 0744)
	assert.Nil(t, err)
	res, err := runner.RunStringWith(toolPath, command.Options{Dir: "subdir"})
	assert.Nil(t, err)
	assert.Equal(t, "list-dir", res.OutputSanitized())
}

func TestShellRunner_RunStringWith_Input(t *testing.T) {
	dir1 := createTempDir(t)
	dir2 := filepath.Join(dir1, "subdir")
	err := os.Mkdir(dir2, 0744)
	assert.Nil(t, err)
	runner := NewMockingShell(dir1, createTempDir(t), "")
	toolPath := filepath.Join(dir2, "list-dir")
	err = ioutil.WriteFile(toolPath, []byte(`#!/usr/bin/env bash
read i1
read i2
echo Hello $i1 and $i2
`), 0744)
	assert.Nil(t, err)
	res, err := runner.RunStringWith(toolPath, command.Options{Input: []string{"one\n", "two\n"}})
	assert.Nil(t, err)
	assert.Equal(t, "Hello one and two", res.OutputSanitized())
}
