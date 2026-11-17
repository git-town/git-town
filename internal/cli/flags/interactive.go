package flags

import (
	"github.com/git-town/git-town/v24/internal/cli/clidomain"
	"github.com/git-town/git-town/v24/internal/messages"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/spf13/cobra"
)

const interactiveLong = "interactive"

// type-safe access to the CLI arguments of type clidomain.Interactive
func Interactive() (AddFunc, ReadInteractiveFlagFunc) {
	negateName := "non-" + interactiveLong
	addFlag := func(cmd *cobra.Command) {
		flags := cmd.Flags()
		flags.Bool(interactiveLong, false, "enable or disable interactive dialogs")
		flags.Bool(negateName, false, "disable interactive dialogs")
	}
	readFlag := func(cmd *cobra.Command) (Option[clidomain.Interactive], error) {
		// read negated flag
		flags := cmd.Flags()
		valueOpt, err := readBoolOptFlag[bool](flags, negateName)
		if err != nil {
			return None[clidomain.Interactive](), err
		}
		value, hasValue := valueOpt.Get()
		if hasValue && value {
			return Some(clidomain.Interactive(messages.InteractivityDisabledViaCLI)), nil
		}

		// read normal flag
		valueOpt, err = readBoolOptFlag[bool](flags, interactiveLong)
		if err != nil {
			return None[clidomain.Interactive](), err
		}
		value, hasValue = valueOpt.Get()
		if !hasValue {
			return None[clidomain.Interactive](), err
		}
		if !value {
			return Some(clidomain.Interactive(messages.InteractivityDisabledViaCLI)), nil
		}
		return Some(clidomain.InteractiveEnabled), nil
	}
	return addFlag, readFlag
}

// ReadInteractiveFlagFunc is the type signature for the function that reads the "interactive" flag from the args to the given Cobra command.
type ReadInteractiveFlagFunc func(*cobra.Command) (Option[clidomain.Interactive], error)
