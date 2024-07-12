package flags

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

const (
	bodyLong  = "body" // long form of the "body" CLI flag
	bodyShort = "b"    // short form of the "body" CLI flag
)

// ProposalBody provides type-safe access to the CLI arguments of type gitdomain.ProposalBody.
func ProposalBody() (AddFunc, ReadProposalBodyFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(bodyLong, bodyShort, "", "provide a body for the proposal")
	}
	readFlag := func(cmd *cobra.Command) gitdomain.ProposalBody {
		value, err := cmd.Flags().GetString(bodyLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), commitMessageLong))
		}
		return gitdomain.ProposalBody(value)
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadProposalBodyFlagFunc func(*cobra.Command) gitdomain.ProposalBody
