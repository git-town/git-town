package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		pushNewBranches, err := prodRepo.Config.ShouldNewBranchPush()
		if err != nil {
			cli.Exit(err)
		}
		pushHook, err := prodRepo.Config.PushHook()
		if err != nil {
			cli.Exit(err)
		}
		isOffline, err := prodRepo.Config.IsOffline()
		if err != nil {
			cli.Exit(err)
		}
		deleteOrigin, err := prodRepo.Config.ShouldShipDeleteOriginBranch()
		if err != nil {
			cli.Exit(err)
		}
		fmt.Println()
		cli.PrintHeader("Branches")
		cli.PrintEntry("main branch", cli.StringSetting(prodRepo.Config.MainBranch()))
		cli.PrintEntry("perennial branches", cli.StringSetting(strings.Join(prodRepo.Config.PerennialBranches(), ", ")))
		fmt.Println()
		cli.PrintHeader("Configuration")
		cli.PrintEntry("offline", cli.BoolSetting(isOffline))
		cli.PrintEntry("pull branch strategy", prodRepo.Config.PullBranchStrategy())
		cli.PrintEntry("run pre-push hook", cli.BoolSetting(pushHook))
		cli.PrintEntry("push new branches", cli.BoolSetting(pushNewBranches))
		cli.PrintEntry("ship removes the remote branch", cli.BoolSetting(deleteOrigin))
		cli.PrintEntry("sync strategy", prodRepo.Config.SyncStrategy())
		cli.PrintEntry("sync with upstream", cli.BoolSetting(prodRepo.Config.ShouldSyncUpstream()))
		fmt.Println()
		cli.PrintHeader("Hosting")
		cli.PrintEntry("hosting service override", cli.StringSetting(prodRepo.Config.HostingService()))
		cli.PrintEntry("GitHub token", cli.StringSetting(prodRepo.Config.GitHubToken()))
		cli.PrintEntry("GitLab token", cli.StringSetting(prodRepo.Config.GitLabToken()))
		cli.PrintEntry("Gitea token", cli.StringSetting(prodRepo.Config.GiteaToken()))
		fmt.Println()
		if prodRepo.Config.MainBranch() != "" {
			cli.PrintLabelAndValue("Branch Ancestry", cli.PrintableBranchAncestry(&prodRepo.Config))
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

// MAIN BRANCH SUBCOMMAND

var mainBranchConfigCommand = &cobra.Command{
	Use:   "main-branch [<branch>]",
	Short: "Displays or sets your main development branch",
	Long: `Displays or sets your main development branch

The main branch is the Git branch from which new feature branches are cut.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printMainBranch()
		} else {
			err := setMainBranch(args[0], prodRepo)
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func printMainBranch() {
	cli.Println(cli.StringSetting(prodRepo.Config.MainBranch()))
}

func setMainBranch(branchName string, repo *git.ProdRepo) error {
	hasBranch, err := repo.Silent.HasLocalBranch(branchName)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branchName)
	}
	return repo.Config.SetMainBranch(branchName)
}

// OFFLINE SUBCOMMAND

var offlineCommand = &cobra.Command{
	Use:   "offline [(yes | no)]",
	Short: "Displays or sets offline mode",
	Long: `Displays or sets offline mode

Git Town avoids network operations in offline mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			isOffline, err := prodRepo.Config.IsOffline()
			if err != nil {
				cli.Exit(err)
			}
			cli.Println(cli.FormatBool(isOffline))
		} else {
			value, err := cli.ParseBool(args[0])
			if err != nil {
				cli.Exit(fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no".\n`, args[0]))
			}
			err = prodRepo.Config.SetOffline(value)
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

// PERENNIAL-BRANCHES SUBCOMMAND

var perennialBranchesCommand = &cobra.Command{
	Use:   "perennial-branches",
	Short: "Displays your perennial branches",
	Long: `Displays your perennial branches

Perennial branches are long-lived branches.
They cannot be shipped.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Println(cli.StringSetting(strings.Join(prodRepo.Config.PerennialBranches(), "\n")))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

var updatePrennialBranchesCommand = &cobra.Command{
	Use:   "update",
	Short: "Prompts to update your perennial branches",
	Long:  `Prompts to update your perennial branches`,
	Run: func(cmd *cobra.Command, args []string) {
		err := dialog.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

// PULL-BRANCH-STRATEGY SUBCOMMAND

var pullBranchStrategyCommand = &cobra.Command{
	Use:   "pull-branch-strategy [(rebase | merge)]",
	Short: "Displays or sets your pull branch strategy",
	Long: `Displays or sets your pull branch strategy

The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Println(prodRepo.Config.PullBranchStrategy())
		} else {
			err := prodRepo.Config.SetPullBranchStrategy(args[0])
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] != "rebase" && args[0] != "merge" {
			return fmt.Errorf("invalid value: %q", args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

// PUSH-NEW-BRANCHES SUBCOMMAND

var pushNewBranchesCommand = &cobra.Command{
	Use:   "push-new-branches [--global] [(yes | no)]",
	Short: "Displays or changes whether new branches get pushed to origin",
	Long: `Displays or changes whether new branches get pushed to origin.

If "push-new-branches" is true, the Git Town commands hack, append, and prepend
push the new branch to the origin remote.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := printPushNewBranches(prodRepo)
			if err != nil {
				cli.Exit(err)
			}
		} else {
			err := setPushNewBranches(args[0], prodRepo)
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func printPushNewBranches(repo *git.ProdRepo) error {
	var setting bool
	var err error
	if globalFlag {
		setting, err = repo.Config.ShouldNewBranchPushGlobal()
	} else {
		setting, err = repo.Config.ShouldNewBranchPush()
	}
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(setting))
	return nil
}

func setPushNewBranches(text string, repo *git.ProdRepo) error {
	value, err := cli.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	return repo.Config.SetNewBranchPush(value, globalFlag)
}

// PUSH-HOOK SUBCOMMAND

var pushHookCommand = &cobra.Command{
	Use:   "push-hook [--global] [(yes | no)]",
	Short: "Configures whether Git Town should run Git's pre-push hook.",
	Long: `Configures whether Git Town should run Git's pre-push hook.

Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) == 0 {
			err = printPushHook(prodRepo)
		} else {
			err = setPushHook(args[0], prodRepo)
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

func printPushHook(repo *git.ProdRepo) error {
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

func setPushHook(text string, repo *git.ProdRepo) error {
	value, err := cli.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	if globalFlag {
		return repo.Config.SetPushHookGlobally(value)
	}
	return repo.Config.SetPushHookLocally(value)
}

// RESET SUBCOMMAND

var resetConfigCommand = &cobra.Command{
	Use:   "reset",
	Short: "Resets your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := prodRepo.Config.RemoveLocalGitConfiguration()
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

// SETUP SUBCOMMAND

var setupConfigCommand = &cobra.Command{
	Use:   "setup",
	Short: "Prompts to setup your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := dialog.ConfigureMainBranch(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		err = dialog.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

// SYNC-STRATEGY SUBCOMMAND

var syncStrategyCommand = &cobra.Command{
	Use:   "sync-strategy [(merge | rebase)]",
	Short: "Displays or sets your sync strategy",
	Long: `Displays or sets your sync strategy

The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printSyncStrategy(prodRepo)
		} else {
			err := setSyncStrategy(prodRepo, args[0])
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] != "merge" && args[0] != "rebase" {
			return fmt.Errorf("invalid value: %q", args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func printSyncStrategy(repo *git.ProdRepo) {
	var strategy string
	if globalFlag {
		strategy = repo.Config.SyncStrategyGlobal()
	} else {
		strategy = repo.Config.SyncStrategy()
	}
	cli.Println(strategy)
}

func setSyncStrategy(repo *git.ProdRepo, value string) error {
	if globalFlag {
		return repo.Config.SetSyncStrategyGlobal(value)
	}
	return repo.Config.SetSyncStrategy(value)
}

func init() {
	configCommand.AddCommand(mainBranchConfigCommand)
	pushNewBranchesCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	configCommand.AddCommand(pushNewBranchesCommand)
	configCommand.AddCommand(offlineCommand)
	perennialBranchesCommand.AddCommand(updatePrennialBranchesCommand)
	configCommand.AddCommand(perennialBranchesCommand)
	configCommand.AddCommand(pullBranchStrategyCommand)
	configCommand.AddCommand(pushHookCommand)
	pushHookCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global push hook flag")
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	configCommand.AddCommand(syncStrategyCommand)
	syncStrategyCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global sync strategy")
	RootCmd.AddCommand(configCommand)
}
