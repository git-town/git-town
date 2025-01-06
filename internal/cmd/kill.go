package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/spf13/cobra"
)

func killCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:    "kill [<branch>]",
		Hidden: true,
		Args:   cobra.MaximumNArgs(1),
		Short:  deleteDesc,
		Long:   cmdhelpers.Long(deleteDesc, deleteHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			printKillDeprecationNotice()
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			result := executeDelete(args, dryRun, verbose)
			printKillDeprecationNotice()
			return result
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func printKillDeprecationNotice() {
	fmt.Println(messages.KillDeprecation)
	time.Sleep(2000 * time.Millisecond)
}
