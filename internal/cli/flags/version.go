package flags

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/spf13/cobra"
)

const versionLong = "version"

func Version() (AddFunc, ReadBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(versionLong, "V", false, "display the version number")
	}
	readFlag := func(cmd *cobra.Command) bool {
		value, err := cmd.Flags().GetBool(versionLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), versionLong))
		}
		return value
	}
	return addFlag, readFlag
}

type ReadBoolFlagFunc func(*cobra.Command) bool
