package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const noAutoResolveLong = "no-auto-resolve"

// type-safe access to the CLI arguments of type configdomain.NoAutoResolve
func NoAutoResolve() (AddFunc, ReadNoAutoResolveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(noAutoResolveLong, "", false, "don't auto-resolve phantom merge conflicts")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.NoAutoResolve, error) {
		value, err := cmd.Flags().GetBool(noAutoResolveLong)
		return configdomain.NoAutoResolve(value), err
	}
	return addFlag, readFlag
}

// ReadNoAutoResolveFlagFunc is the type signature for the function that reads the "no-auto-resolve" flag from the args to the given Cobra command.
type ReadNoAutoResolveFlagFunc func(*cobra.Command) (configdomain.NoAutoResolve, error)
