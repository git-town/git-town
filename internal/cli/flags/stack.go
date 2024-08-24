package flags

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const stackLong = "stack"

// type-safe access to the CLI arguments of type configdomain.FullStack
func Stack(description string) (AddFunc, ReadStackFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(stackLong, "s", false, description)
	}
	readFlag := func(cmd *cobra.Command) configdomain.FullStack {
		value, err := cmd.Flags().GetBool(stackLong)
		if err != nil {
			panic(err)
		}
		return configdomain.FullStack(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadStackFlagFunc func(*cobra.Command) configdomain.FullStack
