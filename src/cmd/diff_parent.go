package cmd

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/spf13/cobra"
)

type diffParentConfig struct {
	branch       string
	parentBranch string
}

var diffParentCommand = &cobra.Command{
	Use:   "diff-parent [<branch>]",
	Short: "Shows the changes committed to a feature branch",
	Long: `Shows the changes committed to a feature branch

Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getDiffParentConfig(args, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		err = prodRepo.Logging.DiffParent(config.branch, config.parentBranch)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

// Does not return error because "Ensure" functions will call exit directly.
func getDiffParentConfig(args []string, repo *git.ProdRepo) (config diffParentConfig, err error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return config, err
	}
	if len(args) == 0 {
		config.branch = initialBranch
	} else {
		config.branch = args[0]
	}
	if initialBranch != config.branch {
		hasBranch, err := repo.Silent.HasLocalBranch(config.branch)
		if err != nil {
			return config, err
		}
		if !hasBranch {
			return config, fmt.Errorf("there is no local branch named %q", config.branch)
		}
	}
	if !prodRepo.Config.IsFeatureBranch(config.branch) {
		return config, fmt.Errorf("you can only diff-parent feature branches")
	}
	err = prompt.EnsureKnowsParentBranches([]string{config.branch}, repo)
	if err != nil {
		return config, err
	}
	config.parentBranch = repo.Config.GetParentBranch(config.branch)
	return config, nil
}

func init() {
	RootCmd.AddCommand(diffParentCommand)
}
