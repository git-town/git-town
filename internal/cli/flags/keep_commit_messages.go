package flags

import (
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const keepCommitMessagesLong = "keep-commit-messages"

// type-safe access to the CLI arguments of type configdomain.KeepCommitMessages
func KeepCommitMessages() (AddFunc, ReadKeepCommitMessagesFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(keepCommitMessagesLong, false, "preserve the commit messages of squashed commits")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.KeepCommitMessages, error) {
		return readBoolFlag[configdomain.KeepCommitMessages](cmd.Flags(), keepCommitMessagesLong)
	}
	return addFlag, readFlag
}

// ReadKeepCommitMessagesFlagFunc is the type signature for the function that reads the "keep-commit-messages" flag from the args to the given Cobra command.
type ReadKeepCommitMessagesFlagFunc func(*cobra.Command) (configdomain.KeepCommitMessages, error)
