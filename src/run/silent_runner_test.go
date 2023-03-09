package run_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func TestSilentRunner(t *testing.T) {
	t.Parallel()
	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		debug := false
		shell := run.SilentRunner{Debug: &debug}
		res, err := shell.Run("echo", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", res.OutputSanitized())
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		debug := false
		shell := run.SilentRunner{Debug: &debug}
		err := shell.RunMany([][]string{
			{"mkdir", "tmp"},
			{"touch", "tmp/first"},
			{"touch", "tmp/second"},
		})
		defer os.RemoveAll("tmp")
		assert.NoError(t, err)
		entries, err := os.ReadDir("tmp")
		assert.NoError(t, err)
		assert.Equal(t, "first", entries[0].Name())
		assert.Equal(t, "second", entries[1].Name())
	})

	t.Run(".RunString()", func(t *testing.T) {
		t.Parallel()
		debug := false
		shell := run.SilentRunner{Debug: &debug}
		_, err := shell.RunString("touch first")
		defer os.Remove("first")
		assert.NoError(t, err)
		_, err = os.Stat("first")
		assert.False(t, os.IsNotExist(err))
	})
}
