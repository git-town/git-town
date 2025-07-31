package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const autoResolveLong = "no-auto-resolve"

// type-safe access to the CLI arguments of type configdomain.AutoResolve
func AutoResolve() (AddFunc, ReadNoAutoResolveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(autoResolveLong, "", true, "whether to auto-resolve phantom merge conflicts")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.AutoResolve, error) {
		value, err := cmd.Flags().GetBool(autoResolveLong)
		return configdomain.AutoResolve(value), err
	}
	return addFlag, readFlag
}

// ReadNoAutoResolveFlagFunc is the type signature for the function that reads the "no-auto-resolve" flag from the args to the given Cobra command.
type ReadNoAutoResolveFlagFunc func(*cobra.Command) (configdomain.AutoResolve, error)
