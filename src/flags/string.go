package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Bool provides mistake-safe access to string Cobra command-line flags.
func String(name, short, defaultValue, desc string) (AddFunc, ReadStringFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(name, short, defaultValue, desc)
	}
	readFlag := func(cmd *cobra.Command) string {
		value, err := cmd.Flags().GetString(name)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a string %q flag", cmd.Name(), name))
		}
		return value
	}
	return addFlag, readFlag
}

// ReadStringFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadStringFlagFunc func(*cobra.Command) string
