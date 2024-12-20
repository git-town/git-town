package flags

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const proposeLong = "propose"

// type-safe access to the CLI arguments of type gitdomain.Propose
func Propose() (AddFunc, ReadProposeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(proposeLong, "p", false, "propose the new branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Propose, error) {
		value, err := cmd.Flags().GetBool(proposeLong)
		return configdomain.Propose(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadProposeFlagFunc func(*cobra.Command) (configdomain.Propose, error)
