package flags

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

const forceLong = "force"

// type-safe access to the CLI arguments of type configdomain.Force
func Force(desc string) (AddFunc, ReadForceFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(forceLong, "f", false, desc)
	}
	readFlag := func(cmd *cobra.Command) configdomain.Force {
		value, err := cmd.Flags().GetBool(forceLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), forceLong))
		}
		return configdomain.Force(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the force flag from the args to the given Cobra command
type ReadForceFlagFunc func(*cobra.Command) configdomain.Force
