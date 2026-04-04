package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const proposeHeadlessLong = "headless"

// type-safe access to the CLI arguments of type configdomain.ProposeHeadless
func ProposeHeadless() (AddFunc, ReadProposeHeadlessFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(proposeHeadlessLong, false, "create the proposal without opening a browser")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.ProposeHeadless], error) {
		return readBoolOptFlag[configdomain.ProposeHeadless](cmd.Flags(), proposeHeadlessLong)
	}
	return addFlag, readFlag
}

// ReadProposeHeadlessFlagFunc is the type signature for the function that reads the "headless" flag from the args to the given Cobra command.
type ReadProposeHeadlessFlagFunc func(*cobra.Command) (Option[configdomain.ProposeHeadless], error)
