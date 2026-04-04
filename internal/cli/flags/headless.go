package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const headlessLong = "headless"

// type-safe access to the CLI arguments of type configdomain.Headless
func Headless() (AddFunc, ReadHeadlessFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(headlessLong, false, "disable all interactive features")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Headless], error) {
		return readBoolOptFlag[configdomain.Headless](cmd.Flags(), headlessLong)
	}
	return addFlag, readFlag
}

// ReadHeadlessFlagFunc is the type signature for the function that reads the "headless" flag from the args to the given Cobra command.
type ReadHeadlessFlagFunc func(*cobra.Command) (Option[configdomain.Headless], error)
