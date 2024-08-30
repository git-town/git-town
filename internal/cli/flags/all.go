package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const allLong = "all"

// type-safe access to the CLI arguments of type configdomain.SyncAllBranches
func All() (AddFunc, ReadAllFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(allLong, "a", false, "sync all local branches")
	}
	readFlag := func(cmd *cobra.Command) configdomain.SyncAllBranches {
		value, err := cmd.Flags().GetBool(allLong)
		if err != nil {
			panic(err)
		}
		return configdomain.SyncAllBranches(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadAllFlagFunc func(*cobra.Command) configdomain.SyncAllBranches
