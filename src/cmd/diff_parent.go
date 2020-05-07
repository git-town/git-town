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
	Short: "Show differences between current branch and parent branch",
	Long: `Show the difference between a feature branch and its parent

Works on either the current branch or the branch name provided. If the branch has a parent, then
the diff will be output directly. If the branch does not have a parent, one will be asked to
identify the parent branch.

Does not output anything for perennial branches nor the main branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getDiffParentConfig(args)
		runDiffParent(config)
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

func runDiffParent(config diffParentConfig) {
	script.RunCommandSafe("git", "diff", config.parentBranch+".."+config.branch)
}

func init() {
	RootCmd.AddCommand(diffParentCommand)
}
