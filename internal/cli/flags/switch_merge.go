package flags

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/messages"
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
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), mergeLong))
		}
		return configdomain.SwitchUsingMerge(value)
	}
	return addFlag, readFlag
}

type ReadMergeFlagFunc func(*cobra.Command) configdomain.SwitchUsingMerge
