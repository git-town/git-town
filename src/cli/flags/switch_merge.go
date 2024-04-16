package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	mergeMessageLong  = "merge"
	mergeMessageShort = "m"
	mergeDesc         = `When switching branches, if you have local modifications to one or more files that are different between the current branch and the branch to which you are switching, the command refuses to switch branches in order to preserve your modifications in context. However,
with this option, a three-way merge between the current branch, your working tree contents, and the new branch is done, and you will be on the new branch.`
)

// SwitchMerge returns two functions: addFlag and readFlag.
// addFlag is a function that takes a *cobra.Command argument and adds a bool flag
// with the long name "merge" and short name "m" to the command's persistent flags.
// The flag's default value is false and it has a description defined by the
// constant mergeDesc.
// readFlag is a function that takes a *cobra.Command argument and returns
// the boolean value for the "merge" flag. It uses the GetBool method of the command's flags
// to retrieve the value and panics if there is an error or if the flag is not found.
// The returned boolean value indicates whether the flag is set or not.
func SwitchMerge() (AddFunc, ReadSwitchMergeFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(mergeMessageLong, mergeMessageShort, false, mergeDesc)
	}
	readFlag := func(cmd *cobra.Command) bool {
		value, err := cmd.Flags().GetBool(mergeMessageLong)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a string %q flag", cmd.Name(), mergeMessageLong))
		}
		return value
	}
	return addFlag, readFlag
}

// ReadSwitchMergeFlagFunc represents a function type that takes a *cobra.Command
// argument and returns a boolean value. The function is used to read a switch
// or merge flag from the command and determine whether it is set or not.
type ReadSwitchMergeFlagFunc func(*cobra.Command) bool
