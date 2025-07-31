package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const autoResolveLong = "no-auto-resolve"

// type-safe access to the CLI arguments of type configdomain.AutoResolve
func AutoResolve() (AddFunc, ReadAutoResolveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(autoResolveLong, "", true, "whether to auto-resolve phantom merge conflicts")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoResolve], error) {
		if !cmd.Flags().Changed(autoResolveLong) {
			return None[configdomain.AutoResolve](), nil
		}
		value, err := cmd.Flags().GetBool(autoResolveLong)
		return Some(configdomain.AutoResolve(value)), err
	}
	return addFlag, readFlag
}

// ReadAutoResolveFlagFunc is the type signature for the function that reads the "no-auto-resolve" flag from the args to the given Cobra command.
type ReadAutoResolveFlagFunc func(*cobra.Command) (Option[configdomain.AutoResolve], error)
