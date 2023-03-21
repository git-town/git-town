package cmd

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const setParentDesc = "Prompts to set the parent branch for the current branch"

func setParentCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "set-parent",
		GroupID: "lineage",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository, isConfigured),
		Short:   setParentDesc,
		Long:    long(setParentDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setParent(repo)
		},
	}
}

func setParent(repo *git.ProdRepo) error {
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
	return validate.KnowsBranchAncestry(currentBranch, defaultParentBranch, repo)
}
