package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/shoenig/test/must"
	"github.com/spf13/cobra"
)

func TestVerbose(t *testing.T) {
	t.Parallel()

	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Verbose()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--verbose"})
		must.NoError(t, err)
		must.EqOp(t, true, readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Verbose()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-v"})
		must.NoError(t, err)
		must.EqOp(t, true, readFlag(&cmd))
	})
}
