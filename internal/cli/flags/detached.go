package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const detachedLong = "detached"

// type-safe access to the CLI arguments of type configdomain.Detached
func Detached() (AddFunc, ReadDetachedFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(detachedLong, "d", false, "don't update the perennial root branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Detached, error) {
		value, err := cmd.Flags().GetBool(detachedLong)
		return configdomain.Detached(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the detached flag from the args to the given Cobra command
type ReadDetachedFlagFunc func(*cobra.Command) (configdomain.Detached, error)
