package flags

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	bodyLong = "body" // long form of the "body" CLI flag
)

// type-safe access to the CLI arguments of type gitdomain.ProposalBody
func ProposalBody(short string) (AddFunc, ReadProposalBodyFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(bodyLong, short, "", "provide a body for the proposal")
	}
	readFlag := func(cmd *cobra.Command) (Option[gitdomain.ProposalBody], error) {
		value, err := cmd.Flags().GetString(bodyLong)
		return NewOption(gitdomain.ProposalBody(value)), err
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadProposalBodyFlagFunc func(*cobra.Command) (Option[gitdomain.ProposalBody], error)
