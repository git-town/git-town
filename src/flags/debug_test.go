package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	t.Parallel()
	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Debug()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--debug"})
		assert.NoError(t, err)
		assert.Equal(t, true, readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Debug()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-d"})
		assert.NoError(t, err)
		assert.Equal(t, true, readFlag(&cmd))
	})
}
