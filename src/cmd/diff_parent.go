package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func diffParentCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "diff-parent [<branch>]",
		Short: "Shows the changes committed to a feature branch",
		Long: `Shows the changes committed to a feature branch

Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := determineDiffParentConfig(args, repo)
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
		GroupID: "lineage",
	}
}

type diffParentConfig struct {
	branch       string
	parentBranch string
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentConfig(args []string, repo *git.ProdRepo) (*diffParentConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	var branch string
	if len(args) > 0 {
		branch = args[0]
	} else {
		branch = initialBranch
	}
	if initialBranch != branch {
		hasBranch, err := repo.Silent.HasLocalBranch(branch)
		if err != nil {
			return nil, err
		}
		if !hasBranch {
			return nil, fmt.Errorf("there is no local branch named %q", branch)
		}
	}
	if !repo.Config.IsFeatureBranch(branch) {
		return nil, fmt.Errorf("you can only diff-parent feature branches")
	}
	parentDialog := dialog.ParentBranches{}
	err = parentDialog.EnsureKnowsParentBranches([]string{branch}, repo)
	if err != nil {
		return nil, err
	}
	return &diffParentConfig{
		branch:       branch,
		parentBranch: repo.Config.ParentBranch(branch),
	}, nil
}
