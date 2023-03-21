package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const pushHookDesc = "Configures whether Git Town should run Git's pre-push hook."

const pushHookHelp = `
Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`

func pushHookCommand(repo *git.ProdRepo) *cobra.Command {
	var globalFlag bool
	pushHookCmd := cobra.Command{
		Use:     "push-hook [--global] [(yes | no)]",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: ensure(repo, isRepository),
		Short:   pushHookDesc,
		Long:    long(pullBranchDesc, pushHookHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return setPushHook(args[0], globalFlag, repo)
			}
			return printPushHook(globalFlag, repo)
		},
	}
	pushHookCmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global push hook flag")
	return &pushHookCmd
}

func printPushHook(globalFlag bool, repo *git.ProdRepo) error {
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

func setPushHook(text string, global bool, repo *git.ProdRepo) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	if global {
		return repo.Config.SetPushHookGlobally(value)
	}
	return repo.Config.SetPushHookLocally(value)
}
