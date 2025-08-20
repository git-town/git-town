package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const autoResolveLong = "auto-resolve"

// type-safe access to the CLI arguments of type configdomain.AutoResolve
func AutoResolve() (AddFunc, ReadAutoResolveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		defineNegatableBoolFlag(cmd.Flags(), autoResolveLong)
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoResolve], error) {
		return readNegatableBoolFlag[configdomain.AutoResolve](cmd.Flags(), autoResolveLong)
	}
	return addFlag, readFlag
}

// ReadAutoResolveFlagFunc is the type signature for the function that reads the "auto-resolve" flag from the args to the given Cobra command.
type ReadAutoResolveFlagFunc func(*cobra.Command) (Option[configdomain.AutoResolve], error)
