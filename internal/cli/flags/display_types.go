package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const displayTypesLong = "display-types"

// type-safe access to the CLI arguments of type configdomain.Displaytypes
func Displaytypes() (AddFunc, ReadDisplayTypesFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(displayTypesLong, "d", "", "display the branch types")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.DisplayTypes], error) {
		text, err := cmd.Flags().GetString(displayTypesLong)
		if err != nil {
			return None[configdomain.DisplayTypes](), err
		}
		return configdomain.ParseDisplayTypes(text, "CLI flag "+displayTypesLong)
	}
	return addFlag, readFlag
}

// ReadDisplayTypesFlagFunc is the type signature for the function that reads the "display-types" flag from the args to the given Cobra command.
type ReadDisplayTypesFlagFunc func(*cobra.Command) (Option[configdomain.DisplayTypes], error)
