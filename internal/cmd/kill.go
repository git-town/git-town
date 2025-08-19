package cmd

import (
	"cmp"
	"fmt"
	"time"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
			dryRun, errDryRun := readDryRunFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDryRun, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				Detached:    None[configdomain.Detached](),
				DryRun:      dryRun,
				Verbose:     verbose,
			})
			result := executeDelete(args, cliConfig)
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
