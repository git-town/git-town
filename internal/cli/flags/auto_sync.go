package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const autoSyncLong = "auto-sync"

// type-safe access to the CLI arguments of type configdomain.AutoSync
func AutoSync() (AddFunc, ReadAutoSyncFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(autoSyncLong, false, "keep your branches in sync")
		defineNegatedFlag(cmd.Flags(), autoSyncLong)
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoSync], error) {
		return readNegatableFlag[configdomain.AutoSync](cmd.Flags(), autoSyncLong)
	}
	return addFlag, readFlag
}

// ReadAutoSyncFlagFunc is the type signature for the function that reads the "auto-sync" flag from the args to the given Cobra command.
type ReadAutoSyncFlagFunc func(*cobra.Command) (Option[configdomain.AutoSync], error)
