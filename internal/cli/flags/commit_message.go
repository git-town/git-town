package flags

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	commitMessageLong  = "message" // long form of the "commit message" CLI flag
	commitMessageShort = "m"       // short form of the "commit message" CLI flag
)

// type-safe access to the CLI arguments of type gitdomain.CommitMessage
func CommitMessage(desc string) (AddFunc, ReadCommitMessageFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(commitMessageLong, commitMessageShort, "", desc)
	}
	readFlag := func(cmd *cobra.Command) (Option[gitdomain.CommitMessage], error) {
		value, err := cmd.Flags().GetString(commitMessageLong)
		if err != nil {
			return None[gitdomain.CommitMessage](), err
		}
		if value == "" {
			return None[gitdomain.CommitMessage](), nil
		}
		return Some(gitdomain.CommitMessage(value)), nil
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadCommitMessageFlagFunc func(*cobra.Command) (Option[gitdomain.CommitMessage], error)
