package flags

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const (
	proposalLineageLong  = "lineage" // long form of the "lineage" CLI flag
	proposalLineageShort = "l"       // short form of the "lineage" CLI flag
)

// type-safe access to the CLI arguments of type bool
func ProposalLineage(desc string, enabledAction configdomain.ProposalLineage) (AddFunc, ReadProposalLineageFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(proposalLineageLong, proposalLineageShort, false, "Include the proposal lineage in the proposal")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.ProposalLineage, error) {
		value, err := cmd.Flags().GetBool(proposalLineageLong)
		if err == nil && value {
			return enabledAction, nil
		}
		return configdomain.ProposalLineageNone, err
	}
	return addFlag, readFlag
}

// ReadProposalLineageFlagFunc reads a boolean from the CLI args.
type ReadProposalLineageFlagFunc func(*cobra.Command) (configdomain.ProposalLineage, error)
