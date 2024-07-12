package flags

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const (
	titleLong  = "title" // long form of the "title" CLI flag
	titleShort = "t"     // short form of the "title" CLI flag
)

// ProposalTitle provides type-safe access to the CLI arguments of type gitdomain.ProposalTitle.
func ProposalTitle() (AddFunc, ReadProposalTitleFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(titleLong, titleShort, "", "provide a title for the proposal")
	}
	readFlag := func(cmd *cobra.Command) gitdomain.ProposalTitle {
		value, err := cmd.Flags().GetString(titleLong)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a string %q flag", cmd.Name(), commitMessageLong))
		}
		return gitdomain.ProposalTitle(value)
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadProposalTitleFlagFunc func(*cobra.Command) gitdomain.ProposalTitle
