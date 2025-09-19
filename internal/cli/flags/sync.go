package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const syncLong = "sync"

// type-safe access to the CLI arguments of type configdomain.AutoSync
func Sync() (AddFunc, ReadSyncFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(syncLong, true, "sync branches")
		defineNegatedFlag(cmd.Flags(), syncLong)
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoSync], error) {
		return readNegatableFlag[configdomain.AutoSync](cmd.Flags(), syncLong)
	}
	return addFlag, readFlag
}

// ReadPrototypeFlagFunc is the type signature for the function that reads the "prototype" flag from the args to the given Cobra command.
type ReadSyncFlagFunc func(*cobra.Command) (Option[configdomain.AutoSync], error)
