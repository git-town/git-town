package flags

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const prototypeLong = "prototype"

// type-safe access to the CLI arguments of type gitdomain.Prototype
func Prototype() (AddFunc, ReadPrototypeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(prototypeLong, "p", false, "create a prototype branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Prototype, error) {
		value, err := cmd.Flags().GetBool(prototypeLong)
		return configdomain.Prototype(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadPrototypeFlagFunc func(*cobra.Command) (configdomain.Prototype, error)
