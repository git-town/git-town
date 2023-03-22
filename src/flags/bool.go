package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Bool provides mistake-safe access to boolean Cobra command-line flags.
func Bool(name, short, desc string) (AddFunc, readBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(name, short, false, desc)
	}
	readFlag := func(cmd *cobra.Command) bool {
		value, err := cmd.Flags().GetBool(name)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a boolean %q flag", cmd.Name(), name))
		}
		return value
	}
	return addFlag, readFlag
}

// readBoolFlagFunc defines the type signature for helper functions that provide the value a boolean CLI flag associated with a Cobra command.
type readBoolFlagFunc func(*cobra.Command) bool
