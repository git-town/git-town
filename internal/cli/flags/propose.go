package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const proposeLong = "propose"

// type-safe access to the CLI arguments of type configdomain.Propose
func Propose() (AddFunc, ReadProposeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(proposeLong, "", false, "propose the new branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Propose, error) {
		return readBoolFlag[configdomain.Propose](cmd.Flags(), proposeLong)
	}
	return addFlag, readFlag
}

// ReadProposeFlagFunc is the type signature for the function that reads the "propose" flag from the args to the given Cobra command.
type ReadProposeFlagFunc func(*cobra.Command) (configdomain.Propose, error)
