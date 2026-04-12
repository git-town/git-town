package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const interactiveLong = "interactive"

// type-safe access to the CLI arguments of type configdomain.Interactive
func Interactive() (AddFunc, ReadInteractiveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(interactiveLong, false, "enable or disable interactive dialogs")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Interactive], error) {
		valueOpt, err := readBoolOptFlag[bool](cmd.Flags(), interactiveLong)
		if err != nil {
			return None[configdomain.Interactive](), err
		}
		value, hasValue := valueOpt.Get()
		if !hasValue {
			return None[configdomain.Interactive](), err
		}
		if value {
			return Some(configdomain.InteractiveEnabled), nil
		}
		return Some(configdomain.Interactive("interactivity disabled via CLI")), nil
	}
	return addFlag, readFlag
}

// ReadInteractiveFlagFunc is the type signature for the function that reads the "interactive" flag from the args to the given Cobra command.
type ReadInteractiveFlagFunc func(*cobra.Command) (Option[configdomain.Interactive], error)
