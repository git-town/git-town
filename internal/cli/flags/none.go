package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const noParentLong = "none"

// type-safe access to the CLI arguments of type configdomain.None
func NoParent() (AddFunc, ReadNoneFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(noParentLong, false, "set no parent (make perennial)")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.NoParent, error) {
		return readBoolFlag[configdomain.NoParent](cmd.Flags(), noParentLong)
	}
	return addFlag, readFlag
}

// ReadNoneFlagFunc is the type signature for the function that reads the "none" flag from the args to the given Cobra command.
type ReadNoneFlagFunc func(*cobra.Command) (configdomain.NoParent, error)
