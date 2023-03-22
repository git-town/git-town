package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	cmd := cobra.Command{}
	addFlag, readFlag := flags.Debug()
	addFlag(&cmd)
	err := cmd.ParseFlags([]string{"--debug"})
	assert.NoError(t, err)
	assert.Equal(t, true, readFlag(&cmd))
}
