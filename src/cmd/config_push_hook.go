package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const configPushHookSummary = "Configures whether Git Town should run Git's pre-push hook."

const configPushHookDesc = `Configures whether Git Town should run Git's pre-push hook.

Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`

func pushHookCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "push-hook [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: configPushHookSummary,
		Long:  long(configPushHookSummary, configPushHookDesc),
		RunE:  runConfigPushHook,
	}
	addDebugFlag(&cmd)
	cmd.Flags().Bool(globalFlagName, false, "If set, reads or updates the push hook flag for all repos on this machine")
	return &cmd
}

func runConfigPushHook(cmd *cobra.Command, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 readDebugFlag(cmd),
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	globalFlag := readBoolFlag(cmd, globalFlagName)
	if len(args) > 0 {
		return setPushHook(args[0], globalFlag, &repo)
	}
	return printPushHook(globalFlag, &repo)
}

func printPushHook(globalFlag bool, repo *git.PublicRepo) error {
	var setting bool
	var err error
	if globalFlag {
		setting, err = repo.Config.PushHookGlobal()
	} else {
		setting, err = repo.Config.PushHook()
	}
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(setting))
	return nil
}

func setPushHook(text string, global bool, repo *git.PublicRepo) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	if global {
		return repo.Config.SetPushHookGlobally(value)
	}
	return repo.Config.SetPushHookLocally(value)
}
