package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const ignoreUncommittedLong = "ignore-uncommitted"

// type-safe access to the CLI arguments of type configdomain.IgnoreUncommitted
func IgnoreUncommitted() (AddFunc, ReadIgnoreUncommittedFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(ignoreUncommittedLong, false, "ignore uncommitted changes")
		defineNegatedFlag(cmd.Flags(), ignoreUncommittedLong, "warn about uncommitted changes")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.IgnoreUncommitted], error) {
		return readBoolOptFlag[configdomain.IgnoreUncommitted](cmd.Flags(), ignoreUncommittedLong)
	}
	return addFlag, readFlag
}

// ReadIgnoreUncommittedFlagFunc is the type signature for the function that reads the "ignore-uncommitted" flag from the args to the given Cobra command.
type ReadIgnoreUncommittedFlagFunc func(*cobra.Command) (Option[configdomain.IgnoreUncommitted], error)
