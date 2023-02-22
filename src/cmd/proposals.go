package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/stringslice"
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
			for p := range proposals { //nolint:varnamelen
				fmt.Print(".")
				proposal, err := connector.ChangeRequestDetails(proposals[p].Number)
				if err != nil {
					cli.Exit(err)
				}
				proposals[p] = *proposal
			}
			fmt.Println()
			titles := make([]string, len(proposals))
			for p := range proposals {
				titles[p] = proposals[p].Title
			}
			longest := stringslice.Longest(titles)
			format := "  %-" + strconv.Itoa(longest+1) + "s #%d\n"
			for _, proposal := range proposals {
				cli.Print(" ")
				cli.PrintColor(mergeability(proposal.CanMergeWithAPI))
				cli.Printf(format, proposal.Title, proposal.Number)
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
