package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const commitLong = "commit"

// type-safe access to the CLI arguments of type configdomain.Commit
func Commit() (AddFunc, ReadCommitFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(commitLong, "c", false, "commit the stashed changes into the new branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Commit, error) {
		return readBoolFlag[configdomain.Commit](cmd.Flags(), commitLong)
	}
	return addFlag, readFlag
}

// ReadCommitFlagFunc is the type signature for the function that reads the "commit" flag from the args to the given Cobra command.
type ReadCommitFlagFunc func(*cobra.Command) (configdomain.Commit, error)
