package cmd

import (
	"cmp"
	"fmt"
	"time"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/messages"
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
			dryRun, err1 := readDryRunFlag(cmd)
			verbose, err2 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2); err != nil {
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
