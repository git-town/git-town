package flags

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

const displayTypesLong = "display-types"

// type-safe access to the CLI arguments of type configdomain.Displaytypes
func Displaytypes() (AddFunc, ReadDisplayTypesFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(displayTypesLong, "d", false, "display the branch types")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.DisplayTypes, error) {
		value, err := cmd.Flags().GetBool(displayTypesLong)
		return configdomain.DisplayTypes(value), err
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the display-types flag from the args to the given Cobra command
type ReadDisplayTypesFlagFunc func(*cobra.Command) (configdomain.DisplayTypes, error)
