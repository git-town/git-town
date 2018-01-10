package cmd

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/src/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra object
var RootCmd = &cobra.Command{
	Use:   "git-town",
	Short: "Generic, high-level Git workflow support",
	Long: `Git Town makes software development teams who use Git even more productive and happy.

It adds Git commands that support GitHub Flow, Git Flow, the Nvie model, GitLab Flow, and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.`,
}

// Execute runs the Cobra stack
func Execute() {
	git.EnsureVersionRequirementSatisfied()
	color.NoColor = false // Prevent color from auto disable

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
