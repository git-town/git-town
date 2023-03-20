package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pushHookCommand() *cobra.Command {
	globalFlag := false
	debug := false
	cmd := cobra.Command{
		Use:   "push-hook [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Configures whether Git Town should run Git's pre-push hook.",
		Long: `Configures whether Git Town should run Git's pre-push hook.

Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigurePushHook(debug, globalFlag, args)
		},
	}
	debugFlag(&cmd, &debug)
	cmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global push hook flag")
	return &cmd
}

func runConfigurePushHook(debug, global bool, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setPushHook(args[0], global, &repo)
	}
	return printPushHook(global, &repo)
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
