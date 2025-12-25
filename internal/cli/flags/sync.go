package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const syncLong = "sync"

// Sync provides type-safe access to the CLI arguments of type configdomain.AutoSync.
func Sync() (AddFunc, ReadSyncFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(syncLong, true, "sync branches")
		defineNegatedFlag(cmd.Flags(), syncLong, "don't sync branches")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoSync], error) {
		return readNegatableFlag[configdomain.AutoSync](cmd.Flags(), syncLong)
	}
	return addFlag, readFlag
}

type ReadSyncFlagFunc func(*cobra.Command) (Option[configdomain.AutoSync], error)
