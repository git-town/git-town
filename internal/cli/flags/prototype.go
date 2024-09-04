package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const prototypeLong = "prototype"

// type-safe access to the CLI arguments of type gitdomain.Prototype
func Prototype() (AddFunc, ReadPrototypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(prototypeLong, "p", false, "create a prototype branch")
	}
	readFlag := func(cmd *cobra.Command) configdomain.Prototype {
		value, err := cmd.Flags().GetBool(prototypeLong)
		if err != nil {
			panic(err)
		}
		return configdomain.Prototype(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadPrototypeFlagFunc func(*cobra.Command) configdomain.Prototype
