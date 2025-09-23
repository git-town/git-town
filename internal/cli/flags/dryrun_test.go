package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
	"github.com/spf13/cobra"
)

func TestDryRun(t *testing.T) {
	t.Parallel()

	t.Run("user provides flag", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.DryRun()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--dry-run"})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.True(t, have.EqualSome(true))
	})

	t.Run("user provides no flag", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.DryRun()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{""})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.Eq(t, None[configdomain.DryRun](), have)
	})
}
