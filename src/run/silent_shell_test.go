package run_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func TestSilentShell(t *testing.T) {
	t.Parallel()
	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		shell := run.SilentShell{}
		res, err := shell.Run("echo", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", res.OutputSanitized())
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		shell := run.SilentShell{}
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
	})

	t.Run(".RunString()", func(t *testing.T) {
		t.Parallel()
		shell := run.SilentShell{}
		_, err := shell.RunString("touch first")
		defer os.Remove("first")
		assert.NoError(t, err)
		_, err = os.Stat("first")
		assert.False(t, os.IsNotExist(err))
	})

	t.Run(".RunStringWith()", func(t *testing.T) {
		t.Parallel()
		shell := run.SilentShell{}
		res, err := shell.RunStringWith("ls -1", run.Options{Dir: ".."})
		assert.NoError(t, err)
		assert.Contains(t, res.OutputSanitized(), "cmd")
	})
}
