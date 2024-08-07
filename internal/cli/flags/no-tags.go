package flags

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/spf13/cobra"
)

const noTagsLong = "no-tags"

// type-safe access to the CLI arguments of type gitdomain.NoTags
func NoTags() (AddFunc, ReadNoTagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(noTagsLong, "", false, "do not sync tags")
	}
	readFlag := func(cmd *cobra.Command) configdomain.SyncTags {
		value, err := cmd.Flags().GetBool(noPushLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), noTagsLong))
		}
		return configdomain.SyncTags(!value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadNoTagFunc func(*cobra.Command) configdomain.SyncTags
