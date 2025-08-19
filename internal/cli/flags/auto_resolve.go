package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const autoResolveLong = "auto-resolve"

// type-safe access to the CLI arguments of type configdomain.AutoResolve
func AutoResolve() (AddFunc, ReadAutoResolveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		// Defining a string flag here, even though this is technically a bool flag,
		// so that we can parse it using our expanded bool syntax.
		cmd.Flags().StringP(autoResolveLong, "", "yes", "whether to auto-resolve phantom merge conflicts")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.AutoResolve], error) {
		if !cmd.Flags().Changed(autoResolveLong) {
			return None[configdomain.AutoResolve](), nil
		}
		text, err := cmd.Flags().GetString(autoResolveLong)
		if err != nil {
			return None[configdomain.AutoResolve](), err
		}
		parsed, err := gohacks.ParseBool[configdomain.AutoResolve](text, autoResolveLong)
		return Some(configdomain.AutoResolve(parsed)), err
	}
	return addFlag, readFlag
}

// ReadAutoResolveFlagFunc is the type signature for the function that reads the "auto-resolve" flag from the args to the given Cobra command.
type ReadAutoResolveFlagFunc func(*cobra.Command) (Option[configdomain.AutoResolve], error)
