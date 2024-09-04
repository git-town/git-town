package flags

import (
	"github.com/spf13/cobra"
)

const versionLong = "version"

func Version() (AddFunc, ReadBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(versionLong, "V", false, "display the version number")
	}
	readFlag := func(cmd *cobra.Command) bool {
		value, err := cmd.Flags().GetBool(versionLong)
		if err != nil {
			panic(err)
		}
		return value
	}
	return addFlag, readFlag
}

type ReadBoolFlagFunc func(*cobra.Command) bool
