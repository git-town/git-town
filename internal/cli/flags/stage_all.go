package flags

import (
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/spf13/cobra"
)

// StageAll provides type-safe access to the CLI arguments of type configdomain.StageAll.
func StageAll() (AddFunc, ReadStageAllFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(allLong, "a", false, "stage all open changes before committing")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.StageAll, error) {
		return readBoolFlag[configdomain.StageAll](cmd.Flags(), allLong)
	}
	return addFlag, readFlag
}

// ReadStageAllFlagFunc is the type signature for the function that reads the "all" flag from the args to the given Cobra command.
type ReadStageAllFlagFunc func(*cobra.Command) (configdomain.StageAll, error)
