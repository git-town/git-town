package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const (
	proposalLong = "proposal" // long form of the "proposal" CLI flag
)

// type-safe access to the CLI arguments of type bool
func Proposal(description string) (AddFunc, ReadProposalFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(proposalLong, "", false, description)
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Proposal, error) {
		value, err := cmd.Flags().GetBool(proposalLong)
		return configdomain.Proposal(value), err
	}
	return addFlag, readFlag
}

// ReadProposalFlagFunc reads a boolean from the CLI args.
type ReadProposalFlagFunc func(*cobra.Command) (configdomain.Proposal, error)
