package flags

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

const (
	verboseLong  = "verbose"
	verboseShort = "v"
)

// ProposalTitle provides type-safe access to the CLI arguments of type gitdomain.ProposalTitle.
func Verbose() (AddFunc, ReadVerboseFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(verboseLong, verboseShort, false, "Display all Git commands run under the hood")
	}
	readFlag := func(cmd *cobra.Command) configdomain.Verbose {
		value, err := cmd.Flags().GetBool(verboseLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), verboseLong))
		}
		return configdomain.Verbose(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadVerboseFlagFunc func(*cobra.Command) configdomain.Verbose
