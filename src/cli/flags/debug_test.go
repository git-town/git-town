package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cli/flags"
	"github.com/shoenig/test/must"
	"github.com/spf13/cobra"
)

func TestDebug(t *testing.T) {
	t.Parallel()

	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Debug()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--debug"})
		must.NoError(t, err)
		must.EqOp(t, true, readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Debug()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-d"})
		must.NoError(t, err)
		must.EqOp(t, true, readFlag(&cmd))
	})
}
