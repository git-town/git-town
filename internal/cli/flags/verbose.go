package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const (
	verboseLong  = "verbose"
	verboseShort = "v"
)

// type-safe access to the CLI arguments of type configdomain.Verbose
func Verbose() (AddFunc, ReadVerboseFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(verboseLong, verboseShort, false, "display all Git commands run under the hood")
	}
	readFlag := func(cmd *cobra.Command) configdomain.Verbose {
		value, err := cmd.Flags().GetBool(verboseLong)
		if err != nil {
			panic(err)
		}
		return configdomain.Verbose(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the verbose flag from the args to the given Cobra command
type ReadVerboseFlagFunc func(*cobra.Command) configdomain.Verbose
