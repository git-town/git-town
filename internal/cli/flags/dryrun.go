package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const dryRunLong = "dry-run"

// type-safe access to the CLI arguments of type configdomain.DryRun
func DryRun() (AddFunc, ReadDryRunFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(dryRunLong, "", false, "print but do not run the Git commands")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.DryRun], error) {
		if !cmd.Flags().Changed(dryRunLong) {
			return None[configdomain.DryRun](), nil
		}
		value, err := cmd.Flags().GetBool(dryRunLong)
		return Some(configdomain.DryRun(value)), err
	}
	return addFlag, readFlag
}

// ReadDryRunFlagFunc is the type signature for the function that reads the "dry-run" flag from the args to the given Cobra command.
type ReadDryRunFlagFunc func(*cobra.Command) (Option[configdomain.DryRun], error)
