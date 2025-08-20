package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	verboseLong  = "verbose"
	verboseShort = "v"
)

// type-safe access to the CLI arguments of type configdomain.Verbose
func Verbose() (AddFunc, ReadVerboseFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(verboseLong, verboseShort, false, "display all Git commands run under the hood")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Verbose], error) {
		return readBoolOptFlag[configdomain.Verbose](cmd.Flags(), verboseLong)
	}
	return addFlag, readFlag
}

// ReadVerboseFlagFunc is the type signature for the function that reads the "verbose" flag from the args to the given Cobra command.
type ReadVerboseFlagFunc func(*cobra.Command) (Option[configdomain.Verbose], error)
