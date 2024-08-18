package flags

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	commitMessageLong  = "message" // long form of the "commit message" CLI flag
	commitMessageShort = "m"       // short form of the "commit message" CLI flag
)

// type-safe access to the CLI arguments of type gitdomain.CommitMessage
func CommitMessage(desc string) (AddFunc, ReadCommitMessageFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(commitMessageLong, commitMessageShort, "", desc)
	}
	readFlag := func(cmd *cobra.Command) Option[gitdomain.CommitMessage] {
		value, err := cmd.Flags().GetString(commitMessageLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), commitMessageLong))
		}
		if value == "" {
			return None[gitdomain.CommitMessage]()
		}
		return Some(gitdomain.CommitMessage(value))
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadCommitMessageFlagFunc func(*cobra.Command) Option[gitdomain.CommitMessage]
