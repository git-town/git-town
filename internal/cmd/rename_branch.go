package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/spf13/cobra"
)

func renameBranchCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addForceFlag, readForceFlag := flags.Force("force rename of perennial branch")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:    "rename-branch",
		Hidden: true,
		Args:   cobra.RangeArgs(1, 2),
		Short:  renameDesc,
		Long:   cmdhelpers.Long(renameDesc, renameHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			printRenameBranchDeprecationNotice()
			result := executeRename(args, readDryRunFlag(cmd), readForceFlag(cmd), readVerboseFlag(cmd))
			printRenameBranchDeprecationNotice()
			return result
		},
	}
	addDryRunFlag(&cmd)
	addForceFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func printRenameBranchDeprecationNotice() {
	fmt.Println(messages.RenameBranchDeprecation)
	time.Sleep(2000 * time.Millisecond)
}
