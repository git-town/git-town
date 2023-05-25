package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Parallel()
	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.String("myflag", "m", "default", "desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--myflag", "my-value"})
		assert.NoError(t, err)
		assert.Equal(t, "my-value", readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.String("myflag", "m", "default", "desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-m", "my-value"})
		assert.NoError(t, err)
		assert.Equal(t, "my-value", readFlag(&cmd))
	})
}
