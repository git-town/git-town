package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const prototypeLong = "prototype"

// type-safe access to the CLI arguments of type configdomain.Prototype
func Prototype() (AddFunc, ReadPrototypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(prototypeLong, "p", false, "create a prototype branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Prototype, error) {
		return readBoolFlag[configdomain.Prototype](cmd.Flags(), prototypeLong)
	}
	return addFlag, readFlag
}

// ReadPrototypeFlagFunc is the type signature for the function that reads the "prototype" flag from the args to the given Cobra command.
type ReadPrototypeFlagFunc func(*cobra.Command) (configdomain.Prototype, error)
