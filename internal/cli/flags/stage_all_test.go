package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/cli/flags"
	"github.com/shoenig/test/must"
	"github.com/spf13/cobra"
)

func TestStageAll(t *testing.T) {
	t.Parallel()

	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.StageAll()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--all"})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.True(t, have.Enabled())
	})

	t.Run("nothing given", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.StageAll()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{""})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.False(t, have.Enabled())
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.StageAll()
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-a"})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		must.True(t, have.Enabled())
	})
}
