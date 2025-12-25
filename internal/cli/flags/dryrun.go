package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const dryRunLong = "dry-run"

// type-safe access to the CLI arguments of type configdomain.DryRun
func DryRun() (AddFunc, ReadDryRunFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(dryRunLong, "", false, "print but do not run the Git commands")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.DryRun], error) {
		return readBoolOptFlag[configdomain.DryRun](cmd.Flags(), dryRunLong)
	}
	return addFlag, readFlag
}

// ReadDryRunFlagFunc is the type signature for the function that reads the "dry-run" flag from the args to the given Cobra command.
type ReadDryRunFlagFunc func(*cobra.Command) (Option[configdomain.DryRun], error)
