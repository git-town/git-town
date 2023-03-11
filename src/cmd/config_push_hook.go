package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pushHookCommand(repo *git.ProdRepo) *cobra.Command {
	var globalFlag bool
	pushHookCmd := cobra.Command{
		Use:     "push-hook [--global] [(yes | no)]",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: Ensure(repo, IsRepository),
		Short:   "Configures whether Git Town should run Git's pre-push hook.",
		Long: `Configures whether Git Town should run Git's pre-push hook.

Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`,
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
	value, err := cli.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	if global {
		return repo.Config.SetPushHookGlobally(value)
	}
	return repo.Config.SetPushHookLocally(value)
}
