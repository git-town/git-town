package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	autoResolveLong = "auto-resolve"
	noAutoResolve   = "no-" + autoResolveLong
)

// type-safe access to the CLI arguments of type configdomain.AutoResolve
func AutoResolve() (AddFunc, ReadAutoResolveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		// Defining a string flag here, even though this is technically a bool flag,
		// so that we can parse it using our expanded bool syntax.
		flags := cmd.Flags()
		flags.Bool(autoResolveLong, true, "whether to auto-resolve phantom merge conflicts")
		flags.Bool(noAutoResolve, false, "")
		if err := flags.MarkHidden(noAutoResolve); err != nil {
			panic(err)
		}
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoResolve], error) {
		flags := cmd.Flags()
		if flags.Changed(autoResolveLong) {
			value, err := flags.GetBool(autoResolveLong)
			if err != nil {
				return None[configdomain.AutoResolve](), err
			}
			return Some(configdomain.AutoResolve(value)), nil
		}
		if flags.Changed(noAutoResolve) {
			value, err := flags.GetBool(noAutoResolve)
			if err != nil {
				return None[configdomain.AutoResolve](), err
			}
			return Some(configdomain.AutoResolve(!value)), nil
		}
		return None[configdomain.AutoResolve](), nil
	}
	return addFlag, readFlag
}

// ReadAutoResolveFlagFunc is the type signature for the function that reads the "auto-resolve" flag from the args to the given Cobra command.
type ReadAutoResolveFlagFunc func(*cobra.Command) (Option[configdomain.AutoResolve], error)
