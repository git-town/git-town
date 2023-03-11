package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/validate"
	. "github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func diffParentCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "diff-parent [<branch>]",
		GroupID: "lineage",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: Ensure(repo, HasGitVersion, IsRepository, IsConfigured),
		Short:   "Shows the changes committed to a feature branch",
		Long: `Shows the changes committed to a feature branch

Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := determineDiffParentConfig(args, repo)
			if err != nil {
				return err
			}
			return repo.Logging.DiffParent(config.branch, config.parentBranch)
		},
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
	err = validate.KnowsBranchesAncestry([]string{branch}, repo)
	if err != nil {
		return nil, err
	}
	return &diffParentConfig{
		branch:       branch,
		parentBranch: repo.Config.ParentBranch(branch),
	}, nil
}
