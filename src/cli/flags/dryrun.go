package flags

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

const dryRunLong = "dry-run"

// type-safe access to the CLI arguments of type configdomain.DryRun
func DryRun() (AddFunc, ReadDryRunFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(dryRunLong, "", false, "print but do not run the Git commands")
	}
	readFlag := func(cmd *cobra.Command) configdomain.DryRun {
		value, err := cmd.Flags().GetBool(dryRunLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), dryRunLong))
		}
		return configdomain.DryRun(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadDryRunFlagFunc func(*cobra.Command) configdomain.DryRun
