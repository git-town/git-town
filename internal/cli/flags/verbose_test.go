package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.True(t, have.EqualSome(true))
	})

	t.Run("nothing given", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Verbose()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{""})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.Eq(t, None[configdomain.Verbose](), have)
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.Verbose()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-v"})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.True(t, have.EqualSome(true))
	})
}
