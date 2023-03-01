package cmd

import (
	"errors"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func setParentCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "set-parent",
		Short: "Prompts to set the parent branch for the current branch",
		Long:  `Prompts to set the parent branch for the current branch`,
		Run: func(cmd *cobra.Command, args []string) {
			currentBranch, err := repo.Silent.CurrentBranch()
			if err != nil {
				cli.Exit(err)
			}
			if !repo.Config.IsFeatureBranch(currentBranch) {
				cli.Exit(errors.New("only feature branches can have parent branches"))
			}
			defaultParentBranch := repo.Config.ParentBranch(currentBranch)
			if defaultParentBranch == "" {
				defaultParentBranch = repo.Config.MainBranch()
			}
			err = repo.Config.RemoveParentBranch(currentBranch)
			if err != nil {
				cli.Exit(err)
			}
			parentDialog := dialog.ParentBranches{}
			err = parentDialog.AskForBranchAncestry(currentBranch, defaultParentBranch, repo)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
		GroupID: "lineage",
	}
}
