package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/flags"
	"github.com/shoenig/test"
	"github.com/spf13/cobra"
)

func TestBool(t *testing.T) {
	t.Parallel()

	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Bool("myflag", "m", "desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--myflag"})
		test.NoError(t, err)
		test.EqOp(t, true, readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Bool("myflag", "m", "desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-m"})
		test.NoError(t, err)
		test.EqOp(t, true, readFlag(&cmd))
	})
}
