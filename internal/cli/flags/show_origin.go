package flags

import (
	"github.com/spf13/cobra"
)

const showOriginLong = "show-origin"

func ShowOrigin() (AddFunc, ReadBoolFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(showOriginLong, false, "display where the configuration information is taken from")
	}
	readFlag := func(cmd *cobra.Command) (bool, error) {
		return readBoolFlag[bool](cmd.Flags(), showOriginLong)
	}
	return addFlag, readFlag
}
