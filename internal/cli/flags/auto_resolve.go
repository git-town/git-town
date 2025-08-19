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
		// Defining a string flag here, even though this is technically a bool flag,
		// so that we can parse it using our expanded bool syntax.
		flags := cmd.Flags()
		flags.Bool(autoResolveLong, true, "whether to auto-resolve phantom merge conflicts")
		noText := "no-" + autoResolveLong
		flags.Bool(noText, false, "")
		if err := flags.MarkHidden(noText); err != nil {
			panic(err)
		}
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoResolve], error) {
		if !cmd.Flags().Changed(autoResolveLong) {
			return None[configdomain.AutoResolve](), nil
		}
		value, err := cmd.Flags().GetBool(autoResolveLong)
		if err != nil {
			return None[configdomain.AutoResolve](), err
		}

		return Some(configdomain.AutoResolve(value)), nil
	}
	return addFlag, readFlag
}

// ReadAutoResolveFlagFunc is the type signature for the function that reads the "auto-resolve" flag from the args to the given Cobra command.
type ReadAutoResolveFlagFunc func(*cobra.Command) (Option[configdomain.AutoResolve], error)
