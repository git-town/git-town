package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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
		must.Eq(t, Some(gitdomain.CommitMessage("my-value")), readFlag(&cmd))
	})

	t.Run("short version", func(t *testing.T) {
		t.Parallel()
		cmd := cobra.Command{}
		addFlag, readFlag := flags.CommitMessage("desc")
		addFlag(&cmd)
		err := cmd.ParseFlags([]string{"-m", "my-value"})
		must.NoError(t, err)
		must.EqOp(t, Some(gitdomain.CommitMessage("my-value")), readFlag(&cmd))
	})
}
