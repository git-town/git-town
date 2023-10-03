package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/flags"
	"github.com/shoenig/test"
	"github.com/spf13/cobra"
)

func TestDryRun(t *testing.T) {
	t.Parallel()
	cmd := cobra.Command{}
	addFlag, readFlag := flags.DryRun()
	addFlag(&cmd)
	err := cmd.ParseFlags([]string{"--dry-run"})
	test.NoError(t, err)
	test.EqOp(t, true, readFlag(&cmd))
}
