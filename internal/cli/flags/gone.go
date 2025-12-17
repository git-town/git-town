package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const goneLong = "gone"

// type-safe access to the CLI arguments of type configdomain.Gone
func Gone() (AddFunc, ReadGoneFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(goneLong, "p", false, "sync only branches whose remote is gone")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.Gone, error) {
		return readBoolFlag[configdomain.Gone](cmd.Flags(), goneLong)
	}
	return addFlag, readFlag
}

type ReadGoneFlagFunc func(*cobra.Command) (configdomain.Gone, error)
