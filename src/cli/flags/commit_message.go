package flags

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const (
	commitMessageLong  = "message" // long form of the "commit message" CLI flag
	commitMessageShort = "m"       // short form of the "commit message" CLI flag
)

// CommitMessage provides type-safe access to the CLI arguments of type gitdomain.CommitMessage.
func CommitMessage(desc string) (AddFunc, ReadCommitMessageFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(commitMessageLong, commitMessageShort, "", desc)
	}
	readFlag := func(cmd *cobra.Command) gitdomain.CommitMessage {
		value, err := cmd.Flags().GetString(commitMessageLong)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a string %q flag", cmd.Name(), commitMessageLong))
		}
		return gitdomain.CommitMessage(value)
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadCommitMessageFlagFunc func(*cobra.Command) gitdomain.CommitMessage
