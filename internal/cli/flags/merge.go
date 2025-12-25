package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const mergeLong = "merge"

// type-safe access to the CLI arguments of type configdomain.ShipIntoNonPerennialParent
func Merge() (AddFunc, ReadMergeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(mergeLong, "m", false, "merge uncommitted changes into the target branch")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.SwitchUsingMerge, error) {
		return readBoolFlag[configdomain.SwitchUsingMerge](cmd.Flags(), mergeLong)
	}
	return addFlag, readFlag
}

type ReadMergeFlagFunc func(*cobra.Command) (configdomain.SwitchUsingMerge, error)
