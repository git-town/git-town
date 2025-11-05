package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const stashLong = "stash"

func Stash() (AddFunc, ReadStashFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(stashLong, false, "stash uncommitted changes when creating branches")
		defineNegatedFlag(cmd.Flags(), stashLong)
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.Stash], error) {
		return readNegatableFlag[configdomain.Stash](cmd.Flags(), stashLong)
	}
	return addFlag, readFlag
}

type ReadStashFlagFunc func(*cobra.Command) (Option[configdomain.Stash], error)
