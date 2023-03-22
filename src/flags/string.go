package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stringFlag provides access to Cobra command-line flags containing strings
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the programmer didn't forget to define or read the flag.
func String(name, short, defaultValue, desc string) (AddFunc, readStringFlagFunc) {
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

// readStringFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type readStringFlagFunc func(*cobra.Command) string
