package cmd

import (
	"errors"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra object.
func RootCmd(repo *git.ProdRepo) *cobra.Command {
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
	rootCmd.AddCommand(abortCmd(repo))
	rootCmd.AddCommand(appendCmd(repo))
	rootCmd.AddCommand(configCmd(repo))
	rootCmd.AddCommand(continueCmd(repo))
	rootCmd.AddCommand(diffParentCommand(repo))
	rootCmd.AddCommand(discardCmd(repo))
	rootCmd.AddCommand(hackCmd(repo))
	rootCmd.AddCommand(installCommand(repo, &rootCmd))
	rootCmd.AddCommand(killCommand(repo))
	rootCmd.AddCommand(newPullRequestCommand(repo))
	rootCmd.AddCommand(prependCommand(repo))
	rootCmd.AddCommand(pruneBranchesCommand(repo))
	rootCmd.AddCommand(renameBranchCommand(repo))
	rootCmd.AddCommand(repoCommand(repo))
	rootCmd.AddCommand(setParentBranchCommand(repo))
	rootCmd.AddCommand(shipCmd(repo))
	rootCmd.AddCommand(skipCmd(repo))
	rootCmd.AddCommand(syncCmd(repo))
	rootCmd.AddCommand(undoCmd(repo))
	rootCmd.AddCommand(versionCmd())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Developer tool to print git commands run under the hood")
	return &rootCmd
}

// Execute runs the Cobra stack.
func Execute() {
	repo := git.NewProdRepo()
	majorVersion, minorVersion, err := repo.Silent.Version()
	if err != nil {
		cli.Exit(err)
	}
	if !IsAcceptableGitVersion(majorVersion, minorVersion) {
		cli.Exit(errors.New("this app requires Git 2.7.0 or higher"))
	}
	color.NoColor = false // Prevent color from auto disable
	if err := RootCmd(&repo).Execute(); err != nil {
		cli.Exit(err)
	}
}

// IsAcceptableGitVersion indicates whether the given Git version works for Git Town.
func IsAcceptableGitVersion(major, minor int) bool {
	return major > 2 || (major == 2 && minor >= 7)
}
