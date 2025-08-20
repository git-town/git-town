package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	detachedLong = "detached"
	noDetached   = "no-" + detachedLong
)

// type-safe access to the CLI arguments of type configdomain.Detached
func Detached() (AddFunc, ReadDetachedFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(detachedLong, "d", false, "don't update the perennial root branch")
		cmd.Flags().Bool(noDetached, true, "")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Detached], error) {
		if cmd.Flags().Changed(detachedLong) {
			value, err := cmd.Flags().GetBool(detachedLong)
			return Some(configdomain.Detached(value)), err
		}
		if cmd.Flags().Changed(noDetached) {
			value, err := cmd.Flags().GetBool(noDetached)
			return Some(configdomain.Detached(!value)), err
		}
		return None[configdomain.Detached](), nil
	}
	return addFlag, readFlag
}

// ReadDetachedFlagFunc is the type signature for the function that reads the "detached" flag from the args to the given Cobra command.
type ReadDetachedFlagFunc func(*cobra.Command) (Option[configdomain.Detached], error)
