package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const pushLong = "push"

// type-safe access to the CLI arguments of type configdomain.PushBranches
func Push() (AddFunc, ReadPushFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(pushLong, "", false, "push local branches")
		defineNegatedFlag(cmd.Flags(), pushLong)
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.PushBranches], error) {
		if !cmd.Flags().Changed(pushLong) {
			return None[configdomain.PushBranches](), nil
		}
		result, err := readBoolFlag[configdomain.PushBranches](cmd.Flags(), pushLong)
		return Some(result), err
	}
	return addFlag, readFlag
}

// ReadPushFlagFunc is the type signature for the function that reads the "no-push" flag from the args to the given Cobra command.
type ReadPushFlagFunc func(*cobra.Command) (Option[configdomain.PushBranches], error)
