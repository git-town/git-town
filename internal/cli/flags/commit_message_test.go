package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
	"github.com/spf13/cobra"
)

func TestCommitMessage(t *testing.T) {
	t.Parallel()

	t.Run("long version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.CommitMessage("desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"--message", "my-value"})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		want := Some(gitdomain.CommitMessage("my-value"))
		must.Eq(t, want, have)
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.CommitMessage("desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-m", "my-value"})
		must.NoError(t, err)
		have, err := readFlag(&cmd)
		must.NoError(t, err)
		want := Some(gitdomain.CommitMessage("my-value"))
		must.Eq(t, want, have)
	})
}
