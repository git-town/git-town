package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const parkLong = "park"

// Park provides type-safe access to the CLI arguments of type Park.
func Park() (AddFunc, ReadParkFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(parkLong, false, "also mark the branch as parked")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Park], error) {
		return readBoolOptFlag[configdomain.Park](cmd.Flags(), parkLong)
	}
	return addFlag, readFlag
}

// ReadParkFlagFunc is the type signature for the function that reads the "verbose" flag from the args to the given Cobra command.
type ReadParkFlagFunc func(*cobra.Command) (Option[configdomain.Park], error)
