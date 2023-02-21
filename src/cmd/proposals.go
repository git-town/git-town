package cmd

import (
	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

func proposalsCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "proposals",
		Short: "Analyzes the Git Town setup",
		Run: func(cmd *cobra.Command, args []string) {
			connector, err := hosting.NewConnector(&repo.Config, &repo.Silent, cli.PrintDriverAction)
			if err != nil {
				cli.Exit(err)
			}
			proposals, err := connector.ChangeRequests()
			if err != nil {
				cli.Exit(err)
			}
			for _, proposal := range proposals {
				cli.Print(" ")
				cli.PrintColor(mergeability(proposal.CanMergeWithAPI))
				cli.Printf("  %s (#%d)\n", proposal.Title, proposal.Number)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			return repo.Config.ValidateIsOnline()
		},
		Hidden: true,
	}
}

func mergeability(mergeable bool) (*color.Color, string) {
	if mergeable {
		return color.New(color.FgGreen), "✔"
	}
	return color.New(color.FgRed), "✗"
}
