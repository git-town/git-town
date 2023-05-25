package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	t.Parallel()
	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Bool("myflag", "m", "desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--myflag"})
		assert.NoError(t, err)
		assert.Equal(t, true, readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Bool("myflag", "m", "desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-m"})
		assert.NoError(t, err)
		assert.Equal(t, true, readFlag(&cmd))
	})
}
