package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	stashLong    = "stash"
	stashDefault = true
)

// type-safe access to the CLI arguments of type configdomain.Stash
func Stash() (AddFunc, ReadStashFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(stashLong, stashDefault, "stash uncommitted changes")
		defineNegatedFlag(cmd.Flags(), stashLong, stashDefault)
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Stash], error) {
		return readNegatableFlag[configdomain.Stash](cmd.Flags(), stashLong)
	}
	return addFlag, readFlag
}

// ReadPrototypeFlagFunc is the type signature for the function that reads the "prototype" flag from the args to the given Cobra command.
type ReadStashFlagFunc func(*cobra.Command) (Option[configdomain.Stash], error)
