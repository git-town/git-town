package cmd

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/util"
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
		config := getDiffParentConfig(args)
		script.RunCommandSafe("git", "diff", config.parentBranch+".."+config.branch)
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

// Does not return error because "Ensure" functions will call exit directly
func getDiffParentConfig(args []string) (config diffParentConfig) {
	initialBranch := git.GetCurrentBranchName()

	if len(args) == 0 {
		config.branch = initialBranch
	} else {
		config.branch = args[0]
	}

	if initialBranch != config.branch {
		git.EnsureHasLocalBranch(config.branch)
	}

	git.Config().EnsureIsFeatureBranch(config.branch, "You can only diff-parent feature branches.")

	prompt.EnsureKnowsParentBranches([]string{config.branch})
	config.parentBranch = git.Config().GetParentBranch(config.branch)
	return
}

func init() {
	RootCmd.AddCommand(diffParentCommand)
}
