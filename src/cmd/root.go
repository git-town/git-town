package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra object.
var RootCmd = &cobra.Command{
	Use:   "git-town",
	Short: "Generic, high-level Git workflow support",
	Long: `Git Town makes software development teams who use Git even more productive and happy.

It adds Git commands that support GitHub Flow, Git Flow, the Nvie model, GitLab Flow, and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cli.SetDebug(debugFlag)
	},
}

// Execute runs the Cobra stack.
func Execute() {
	majorVersion, minorVersion, err := git.Version()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !IsAcceptableVersion(majorVersion, minorVersion) {
		fmt.Println("Git Town requires Git 2.7.0 or higher")
		os.Exit(1)
	}
	color.NoColor = false // Prevent color from auto disable
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// IsAcceptableVersion indicates whether the given Git version works for Git Town.
func IsAcceptableVersion(major, minor int) bool {
	return major > 2 || (major == 2 && minor >= 7)
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Developer tool to print git commands run under the hood")
}
