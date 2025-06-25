package flags

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	titleLong  = "title" // long form of the "title" CLI flag
	titleShort = "t"     // short form of the "title" CLI flag
)

// type-safe access to the CLI arguments of type gitdomain.ProposalTitle
func ProposalTitle() (AddFunc, ReadProposalTitleFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(titleLong, titleShort, "", "provide a title for the proposal")
	}
	readFlag := func(cmd *cobra.Command) (Option[gitdomain.ProposalTitle], error) {
		value, err := cmd.Flags().GetString(titleLong)
		return NewOption(gitdomain.ProposalTitle(value)), err
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadProposalTitleFlagFunc func(*cobra.Command) (Option[gitdomain.ProposalTitle], error)
