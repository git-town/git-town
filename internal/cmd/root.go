package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/spf13/cobra"
)

const rootDesc = "Branching and workflow support for Git"

const rootHelp = `
Git Town helps create, sync, and ship changes efficiently and with minimal merge conflicts.`

func rootCmd() cobra.Command {
	addVersionFlag, readVersionFlag := flags.Version()
	rootCmd := cobra.Command{
		Use:           "git-town",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         rootDesc,
		Long:          cmdhelpers.Long(rootDesc, rootHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
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
		Title: "Commands for stacked changes:",
	}, &cobra.Group{
		ID:    "types",
		Title: "Commands to limit branch syncing:",
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
		fmt.Println("Git Town 15.1.0")
		return nil
	}
	return cmd.Help()
}
