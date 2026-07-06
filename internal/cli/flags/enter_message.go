package flags

import (
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/spf13/cobra"
)

const enterMessageLong = "enter-message"

// type-safe access to the CLI arguments of type configdomain.ShipEnterMessage
func EnterMessage() (AddFunc, ReadEnterMessageFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(enterMessageLong, false, "manually enter the commit message")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.ShipEnterMessage], error) {
		return readBoolOptFlag[configdomain.ShipEnterMessage](cmd.Flags(), enterMessageLong)
	}
	return addFlag, readFlag
}

// ReadEnterMessageFlagFunc is the type signature for the function that reads the "enter-message" flag from the args to the given Cobra command.
type ReadEnterMessageFlagFunc func(*cobra.Command) (Option[configdomain.ShipEnterMessage], error)
