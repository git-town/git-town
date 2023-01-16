package cmd

import (
	"errors"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra object.
func RootCmd() *cobra.Command {
	debugFlag := false
	rootCmd := cobra.Command{
		Use:   "git-town",
		Short: "Generic, high-level Git workflow support",
		Long: `Git Town makes software development teams who use Git even more productive and happy.

It adds Git commands that support GitHub Flow, Git Flow, the Nvie model, GitLab Flow, and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli.SetDebug(debugFlag)
		},
	}
	rootCmd.AddCommand(abortCmd())
	rootCmd.AddCommand(appendCmd())
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(continueCmd())
	rootCmd.AddCommand(diffParentCommand())
	rootCmd.AddCommand(discardCmd())
	rootCmd.AddCommand(hackCmd())
	rootCmd.AddCommand(installCommand(&rootCmd))
	rootCmd.AddCommand(killCommand())
	rootCmd.AddCommand(newPullRequestCommand())
	rootCmd.AddCommand(prependCommand())
	rootCmd.AddCommand(pruneBranchesCommand())
	rootCmd.AddCommand(renameBranchCommand())
	rootCmd.AddCommand(repoCommand())
	rootCmd.AddCommand(setParentBranchCommand())
	rootCmd.AddCommand(shipCmd())
	rootCmd.AddCommand(skipCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(undoCmd())
	rootCmd.AddCommand(versionCmd())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Developer tool to print git commands run under the hood")
	return &rootCmd
}

// Execute runs the Cobra stack.
func Execute() {
	majorVersion, minorVersion, err := prodRepo.Silent.Version()
	if err != nil {
		cli.Exit(err)
	}
	if !IsAcceptableGitVersion(majorVersion, minorVersion) {
		cli.Exit(errors.New("this app requires Git 2.7.0 or higher"))
	}
	color.NoColor = false // Prevent color from auto disable
	if err := RootCmd().Execute(); err != nil {
		cli.Exit(err)
	}
}

// IsAcceptableGitVersion indicates whether the given Git version works for Git Town.
func IsAcceptableGitVersion(major, minor int) bool {
	return major > 2 || (major == 2 && minor >= 7)
}
