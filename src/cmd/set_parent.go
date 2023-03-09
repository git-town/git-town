package cmd

import (
	"errors"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func setParentCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "set-parent",
		Short: "Prompts to set the parent branch for the current branch",
		Long:  `Prompts to set the parent branch for the current branch`,
		RunE: func(cmd *cobra.Command, args []string) error {
			currentBranch, err := repo.Silent.CurrentBranch()
			if err != nil {
				return err
			}
			if !repo.Config.IsFeatureBranch(currentBranch) {
				return errors.New("only feature branches can have parent branches")
			}
			defaultParentBranch := repo.Config.ParentBranch(currentBranch)
			if defaultParentBranch == "" {
				defaultParentBranch = repo.Config.MainBranch()
			}
			err = repo.Config.RemoveParentBranch(currentBranch)
			if err != nil {
				return err
			}
			parentDialog := dialog.ParentBranches{}
			return parentDialog.AskForBranchAncestry(currentBranch, defaultParentBranch, repo)
		},
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository, isConfigured),
		GroupID: "lineage",
	}
}
