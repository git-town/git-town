package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Bool provides mistake-safe access to boolean Cobra command-line flags.
func Bool(name, short, desc string) (AddFunc, ReadBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(name, short, false, desc)
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
