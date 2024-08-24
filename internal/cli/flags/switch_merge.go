package flags

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const mergeLong = "merge"

// type-safe access to the CLI arguments of type configdomain.ShipIntoNonPerennialParent
func SwitchMerge() (AddFunc, ReadMergeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(mergeLong, "m", false, "merge uncommitted changes into the target branch")
	}
	readFlag := func(cmd *cobra.Command) configdomain.SwitchUsingMerge {
		value, err := cmd.Flags().GetBool(mergeLong)
		if err != nil {
			panic(err)
		}
		return configdomain.SwitchUsingMerge(value)
	}
	return addFlag, readFlag
}

type ReadMergeFlagFunc func(*cobra.Command) configdomain.SwitchUsingMerge
