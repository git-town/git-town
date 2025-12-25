package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const stackLong = "stack"

// Stack provides type-safe access to the CLI arguments of type configdomain.FullStack.
func Stack(description string) (AddFunc, ReadStackFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(stackLong, "s", false, description)
	}
	readFlag := func(cmd *cobra.Command) (configdomain.FullStack, error) {
		return readBoolFlag[configdomain.FullStack](cmd.Flags(), stackLong)
	}
	return addFlag, readFlag
}

type ReadStackFlagFunc func(*cobra.Command) (configdomain.FullStack, error)
