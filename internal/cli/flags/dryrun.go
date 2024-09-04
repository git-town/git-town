package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const dryRunLong = "dry-run"

// type-safe access to the CLI arguments of type configdomain.DryRun
func DryRun() (AddFunc, ReadDryRunFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(dryRunLong, "", false, "print but do not run the Git commands")
	}
	readFlag := func(cmd *cobra.Command) configdomain.DryRun {
		value, err := cmd.Flags().GetBool(dryRunLong)
		if err != nil {
			panic(err)
		}
		return configdomain.DryRun(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadDryRunFlagFunc func(*cobra.Command) configdomain.DryRun
