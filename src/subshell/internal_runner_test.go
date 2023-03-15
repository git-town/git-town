package subshell_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestSilentRunner(t *testing.T) {
	t.Parallel()
	t.Run(".Run()", func(t *testing.T) {
		t.Parallel()
		runner := subshell.InternalRunner{}
		res, err := runner.Run("echo", "hello", "world")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", res.OutputSanitized())
	})

	t.Run(".RunMany()", func(t *testing.T) {
		t.Parallel()
		runner := subshell.InternalRunner{}
		err := runner.RunMany([][]string{
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
		runner := subshell.InternalRunner{}
		_, err := runner.RunString("touch first")
		defer os.Remove("first")
		assert.NoError(t, err)
		_, err = os.Stat("first")
		assert.False(t, os.IsNotExist(err))
	})
}
