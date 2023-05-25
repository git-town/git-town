package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDryRun(t *testing.T) {
	t.Parallel()
	cmd := cobra.Command{}
	addFlag, readFlag := flags.DryRun()
	addFlag(&cmd)
	err := cmd.ParseFlags([]string{"--dry-run"})
	assert.NoError(t, err)
	assert.Equal(t, true, readFlag(&cmd))
}
