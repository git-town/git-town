package flags

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const stackLong = "stack"

// type-safe access to the CLI arguments of type configdomain.FullStack
func Stack(description string) (AddFunc, ReadStackFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(stackLong, "s", false, description)
	}
	readFlag := func(cmd *cobra.Command) (configdomain.FullStack, error) {
		value, err := cmd.Flags().GetBool(stackLong)
		return configdomain.FullStack(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadStackFlagFunc func(*cobra.Command) (configdomain.FullStack, error)
