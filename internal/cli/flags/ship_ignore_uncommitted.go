package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const shipIgnoreUncommittedLong = "ignore-uncommitted"

// type-safe access to the CLI arguments of type configdomain.ShipIgnoreUncommitted
func ShipIgnoreUncommitted() (AddFunc, ReadShipIgnoreUncommittedFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(shipIgnoreUncommittedLong, false, "ignore uncommitted changes")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.ShipIgnoreUncommitted], error) {
		return readBoolOptFlag[configdomain.ShipIgnoreUncommitted](cmd.Flags(), shipIgnoreUncommittedLong)
	}
	return addFlag, readFlag
}

// ReadShipIgnoreUncommittedFlagFunc is the type signature for the function that reads the "ignore-uncommitted" flag from the args to the given Cobra command.
type ReadShipIgnoreUncommittedFlagFunc func(*cobra.Command) (Option[configdomain.ShipIgnoreUncommitted], error)
