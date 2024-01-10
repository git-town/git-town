package cmd

import (
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/spf13/cobra"
)

func debugCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "debug <thing>",
		GroupID: "lineage",
		Args:    cobra.MaximumNArgs(1),
		Short:   diffParentDesc,
		Long:    cmdhelpers.Long(diffParentDesc, diffParentHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDiffParent(args, readVerboseFlag(cmd))
		},
	}
}
