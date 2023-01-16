package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pushHookCommand() *cobra.Command {
	var globalFlag bool
	pushHookCmd := cobra.Command{
		Use:   "push-hook [--global] [(yes | no)]",
		Short: "Configures whether Git Town should run Git's pre-push hook.",
		Long: `Configures whether Git Town should run Git's pre-push hook.

Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 0 {
				err = printPushHook(globalFlag, prodRepo)
			} else {
				err = setPushHook(args[0], globalFlag, prodRepo)
			}
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(prodRepo)
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
