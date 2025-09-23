package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const (
	pendingLong  = "pending"
	pendingShort = "p"
)

// type-safe access to the CLI arguments of type configdomain.Pending
func Pending() (AddFunc, ReadPendingFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(pendingLong, pendingShort, false, "display just the name of the pending Git Town command")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Pending, error) {
		return readBoolFlag[configdomain.Pending](cmd.Flags(), pendingLong)
	}
	return addFlag, readFlag
}

// ReadPendingFlagFunc is the type signature for the function that reads the "pending" flag from the args to the given Cobra command.
type ReadPendingFlagFunc func(*cobra.Command) (configdomain.Pending, error)
