package flags

import (
	"github.com/spf13/cobra"
)

const versionLong = "version"

func Version() (AddFunc, ReadBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(versionLong, "V", false, "display the version number")
	}
	readFlag := func(cmd *cobra.Command) (bool, error) {
		value, err := cmd.Flags().GetBool(versionLong)
		return value, err
	}
	return addFlag, readFlag
}

type ReadBoolFlagFunc func(*cobra.Command) (bool, error)
