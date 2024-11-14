package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const allLong = "all"

// type-safe access to the CLI arguments of type configdomain.SyncAllBranches
func All(desc string) (AddFunc, ReadAllFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(allLong, "a", false, desc)
	}
	readFlag := func(cmd *cobra.Command) (configdomain.AllBranches, error) {
		value, err := cmd.Flags().GetBool(allLong)
		return configdomain.AllBranches(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadAllFlagFunc func(*cobra.Command) (configdomain.AllBranches, error)
