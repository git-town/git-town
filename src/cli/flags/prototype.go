package flags

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

const prototypeLong = "prototype"

// type-safe access to the CLI arguments of type gitdomain.Prototype
func Prototype() (AddFunc, ReadPrototypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(prototypeLong, "p", false, "Print but do not run the Git commands")
	}
	readFlag := func(cmd *cobra.Command) configdomain.Prototype {
		value, err := cmd.Flags().GetBool(prototypeLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), prototypeLong))
		}
		return configdomain.Prototype(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadPrototypeFlagFunc func(*cobra.Command) configdomain.Prototype
