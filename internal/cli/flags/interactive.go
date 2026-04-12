package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const interactiveLong = "interactive"

// type-safe access to the CLI arguments of type configdomain.Interactive
func Interactive() (AddFunc, ReadInteractiveFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		flags := cmd.Flags()
		flags.Bool(interactiveLong, false, "enable or disable interactive dialogs")
		negateName := "non-" + interactiveLong
		flags.Bool(negateName, false, "disable interactive dialogs")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Interactive], error) {
		// read negated flag
		flags := cmd.Flags()
		valueOpt, err := readBoolOptFlag[bool](flags, "non-"+interactiveLong)
		if err != nil {
			return None[configdomain.Interactive](), err
		}
		value, hasValue := valueOpt.Get()
		if hasValue && value {
			return Some(configdomain.Interactive(messages.InteractivityDisabledViaCLI)), nil
		}

		// read normal flag
		valueOpt, err = readBoolOptFlag[bool](flags, interactiveLong)
		if err != nil {
			return None[configdomain.Interactive](), err
		}
		value, hasValue = valueOpt.Get()
		if !hasValue {
			return None[configdomain.Interactive](), err
		}
		if !value {
			return Some(configdomain.Interactive(messages.InteractivityDisabledViaCLI)), nil
		}
		return Some(configdomain.InteractiveEnabled), nil
	}
	return addFlag, readFlag
}

// ReadInteractiveFlagFunc is the type signature for the function that reads the "interactive" flag from the args to the given Cobra command.
type ReadInteractiveFlagFunc func(*cobra.Command) (Option[configdomain.Interactive], error)
