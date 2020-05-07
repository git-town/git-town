package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"
	"github.com/spf13/cobra"
)

type diffParentConfig struct {
	InitialBranch string
	TargetBranch  string
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
		stepList := getDiffParentStepList(config)
		runState := steps.NewRunState("diff-parent", stepList)
		err := steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
func getDiffParentConfig(args []string) (result diffParentConfig) {
	result.InitialBranch = git.GetCurrentBranchName()

	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
	} else {
		result.TargetBranch = args[0]
	}

	if result.InitialBranch != result.TargetBranch {
		git.EnsureHasLocalBranch(result.TargetBranch)
	}

	git.Config().EnsureIsFeatureBranch(result.TargetBranch, "You can only diff-parent feature branches.")

	prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	return
}

func getDiffParentStepList(config diffParentConfig) (result steps.StepList) {
	targetBranchParent := git.Config().GetParentBranch(config.TargetBranch)
	result.Append(&steps.DiffParentBranchStep{BranchName: config.TargetBranch, ParentBranch: targetBranchParent})
	return
}

func init() {
	RootCmd.AddCommand(diffParentCommand)
}
