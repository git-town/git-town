package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const nameOnlyLong = "name-only"

// NameOnly provides type-safe access to the CLI arguments of type configdomain.NameOnly.
func NameOnly() (AddFunc, ReadNameOnlyFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().Bool(nameOnlyLong, false, "show only the names of changed files")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.NameOnly], error) {
		return readBoolOptFlag[configdomain.NameOnly](cmd.Flags(), nameOnlyLong)
	}
	return addFlag, readFlag
}

// ReadNameOnlyFlagFunc is the type signature for the function that reads the "name-only" flag from the args to the given Cobra command.
type ReadNameOnlyFlagFunc func(*cobra.Command) (Option[configdomain.NameOnly], error)
