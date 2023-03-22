package cmd

import "github.com/spf13/cobra"

const rootDesc = "Generic, high-level Git workflow support"

const rootHelp = `
Git Town makes software development teams who use Git even more productive and happy.

It adds Git commands that support GitHub Flow, Git Flow, the Nvie model, GitLab Flow, and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.`

func rootCmd() cobra.Command {
	rootCmd := cobra.Command{
		Use:           "git-town",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         rootDesc,
		Long:          long(rootDesc, rootHelp),
	}
	rootCmd.AddGroup(&cobra.Group{
		ID:    "basic",
		Title: "Basic commands:",
	}, &cobra.Group{
		ID:    "errors",
		Title: "Commands to deal with errors:",
	}, &cobra.Group{
		ID:    "lineage",
		Title: "Commands for nested feature branches:",
	}, &cobra.Group{
		ID:    "setup",
		Title: "Commands to set up Git Town on your computer:",
	})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd
}
