package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

type diffParentConfig struct {
	branch       string
	parentBranch string
}

func diffParentCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "diff-parent [<branch>]",
		Short: "Shows the changes committed to a feature branch",
		Long: `Shows the changes committed to a feature branch

Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := createDiffParentConfig(args, repo)
			if err != nil {
				cli.Exit(err)
			}
			err = repo.Logging.DiffParent(config.branch, config.parentBranch)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := ValidateIsRepository(repo)
			if err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
}

// Does not return error because "Ensure" functions will call exit directly.
func createDiffParentConfig(args []string, repo *git.ProdRepo) (diffParentConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return diffParentConfig{}, err
	}
	config := diffParentConfig{}
	if len(args) == 0 {
		config.branch = initialBranch
	} else {
		config.branch = args[0]
	}
	if initialBranch != config.branch {
		hasBranch, err := repo.Silent.HasLocalBranch(config.branch)
		if err != nil {
			return diffParentConfig{}, err
		}
		if !hasBranch {
			return diffParentConfig{}, fmt.Errorf("there is no local branch named %q", config.branch)
		}
	}
	if !repo.Config.IsFeatureBranch(config.branch) {
		return diffParentConfig{}, fmt.Errorf("you can only diff-parent feature branches")
	}
	err = dialog.EnsureKnowsParentBranches([]string{config.branch}, repo)
	if err != nil {
		return diffParentConfig{}, err
	}
	config.parentBranch = repo.Config.ParentBranch(config.branch)
	return config, nil
}
