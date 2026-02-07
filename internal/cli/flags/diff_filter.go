package flags

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const diffFilterLong = "diff-filter"

// type-safe access to the CLI arguments of type gitdomain.DiffFilter
func DiffFilter() (AddFunc, ReadDiffFilterFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().String(diffFilterLong, "", "set Git's --diff-filter flag")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.DiffFilter], error) {
		if !cmd.Flags().Changed(diffFilterLong) {
			return None[configdomain.DiffFilter](), nil
		}
		value, err := cmd.Flags().GetString(diffFilterLong)
		if err != nil {
			return None[configdomain.DiffFilter](), err
		}
		return Some(configdomain.DiffFilter(value)), nil
	}
	return addFlag, readFlag
}

// ReadCommitMessageFlagFunc defines the type signature for helper functions that provide the value a string CLI flag associated with a Cobra command.
type ReadDiffFilterFlagFunc func(*cobra.Command) (Option[configdomain.DiffFilter], error)
