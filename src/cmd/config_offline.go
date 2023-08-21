package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/spf13/cobra"
)

const offlineDesc = "Displays or sets offline mode"

const offlineHelp = `
Git Town avoids network operations in offline mode.`

func offlineCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "offline [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: offlineDesc,
		Long:  long(offlineDesc, offlineHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return offline(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func offline(args []string, debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  true,
		ValidateIsOnline: false,
		ValidateGitRepo:  false,
	})
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return setOfflineStatus(args[0], &repo.Runner)
	}
	return displayOfflineStatus(&repo.Runner)
}

func displayOfflineStatus(run *git.ProdRunner) error {
	isOffline, err := run.Config.IsOffline()
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(isOffline))
	return nil
}

func setOfflineStatus(text string, run *git.ProdRunner) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.ValueInvalid, config.KeyOffline, text)
	}
	return run.Config.SetOffline(value)
}
