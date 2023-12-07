package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/spf13/cobra"
)

const rootDesc = "Generic, high-level Git workflow support"

const rootHelp = `
Git Town makes software development teams who use Git even more productive and happy.

It adds Git commands that support GitHub Flow, Git Flow, the Nvie model, GitLab Flow, and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.`

// The current Git Town version (set at compile time).
var version string

// The time this Git Town binary was compiled (set at compile time).
var buildDate string //nolint:gochecknoglobals

func rootCmd() cobra.Command {
	addVersionFlag, readVersionFlag := flags.Version()
	rootCmd := cobra.Command{
		Use:           "git-town",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         rootDesc,
		Long:          long(rootDesc, rootHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRoot(cmd, readVersionFlag(cmd))
		},
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
	addVersionFlag(&rootCmd)
	return rootCmd
}

func executeRoot(cmd *cobra.Command, showVersion bool) error {
	if showVersion {
		fmt.Printf("Git Town %s (%s)\n", version, buildDate)
		return nil
	}
	return cmd.Help()
}
