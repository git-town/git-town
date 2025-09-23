package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const pruneLong = "prune"

// type-safe access to the CLI arguments of type configdomain.Prune
func Prune() (AddFunc, ReadPruneFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(pruneLong, "p", false, "prune empty branches")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Prune, error) {
		return readBoolFlag[configdomain.Prune](cmd.Flags(), pruneLong)
	}
	return addFlag, readFlag
}

// ReadPruneFlagFunc is the type signature for the function that reads the "prune" flag from the args to the given Cobra command.
type ReadPruneFlagFunc func(*cobra.Command) (configdomain.Prune, error)
