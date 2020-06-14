package run_test

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/src/run"
	"github.com/stretchr/testify/assert"
)

func TestRun_Exec(t *testing.T) {
	res, err := run.Exec("echo", "foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", res.Output())
}

func TestRun_Run_UnknownExecutable(t *testing.T) {
	_, err := run.Exec("zonk")
	assert.Error(t, err)
	var execError *exec.Error
	assert.True(t, errors.As(err, &execError))
}

func TestRun_Run_ExitCode(t *testing.T) {
	_, err := run.Exec("bash", "-c", "exit 2")
	var execError *exec.ExitError
	assert.True(t, errors.As(err, &execError))
	assert.Equal(t, 2, execError.ExitCode())
}

func TestRun_RunInDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	dirPath := filepath.Join(dir, "mydir")
	err = os.Mkdir(dirPath, 0700)
	assert.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dirPath, "one"), []byte{}, 0500)
	assert.NoError(t, err)
	res, err := run.InDir(dirPath, "ls", "-1")
	assert.NoError(t, err)
	assert.Equal(t, "one", res.OutputSanitized())
}

func TestRun_Result_OutputContainsText(t *testing.T) {
	res, err := run.Exec("echo", "hello world how are you?")
	assert.NoError(t, err)
	assert.True(t, res.OutputContainsText("world"), "should contain 'world'")
	assert.False(t, res.OutputContainsText("zonk"), "should not contain 'zonk'")
}

func TestRun_Result_OutputContainsLine(t *testing.T) {
	res, err := run.Exec("echo", "hello world")
	assert.NoError(t, err)
	assert.True(t, res.OutputContainsLine("hello world"), `should contain "hello world"`)
	assert.False(t, res.OutputContainsLine("hello"), `partial match should return false`)
	assert.False(t, res.OutputContainsLine("zonk"), `should not contain "zonk"`)
}
