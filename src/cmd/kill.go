package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"
	"github.com/spf13/cobra"
)

type killConfig struct {
	InitialBranch       string
	IsTargetBranchLocal bool
	TargetBranch        string
}

var killCommand = &cobra.Command{
	Use:   "kill [<branch>]",
	Short: "Removes an obsolete feature branch",
	Long: `Removes an obsolete feature branch

Deletes the current or provided branch from the local and remote repositories.
Does not delete perennial branches nor the main branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getKillConfig(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stepList := getKillStepList(config)
		runState := steps.NewRunState("kill", stepList)
		err = steps.Run(runState)
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

func getKillConfig(args []string) (result killConfig, err error) {
	result.InitialBranch = git.GetCurrentBranchName()

	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
	} else {
		result.TargetBranch = args[0]
	}

	git.Config().EnsureIsFeatureBranch(result.TargetBranch, "You can only kill feature branches.")

	result.IsTargetBranchLocal = git.HasLocalBranch(result.TargetBranch)
	if result.IsTargetBranchLocal {
		prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	}

	if git.HasRemote("origin") && !git.Config().IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}

	if result.InitialBranch != result.TargetBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}

	return
}

func getKillStepList(config killConfig) (result steps.StepList) {
	switch {
	case config.IsTargetBranchLocal:
		targetBranchParent := git.Config().GetParentBranch(config.TargetBranch)
		if git.HasTrackingBranch(config.TargetBranch) && !git.Config().IsOffline() {
			result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: true})
		}
		if config.InitialBranch == config.TargetBranch {
			if git.HasOpenChanges() {
				result.Append(&steps.CommitOpenChangesStep{})
			}
			result.Append(&steps.CheckoutBranchStep{BranchName: targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{BranchName: config.TargetBranch, Force: true})
		for _, child := range git.Config().GetChildBranches(config.TargetBranch) {
			result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: targetBranchParent})
		}
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.TargetBranch})
	case !git.Config().IsOffline():
		result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: false})
	default:
		fmt.Printf("Cannot delete remote branch %q in offline mode", config.TargetBranch)
		os.Exit(1)
	}
	result.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.InitialBranch != config.TargetBranch && config.TargetBranch == git.GetPreviouslyCheckedOutBranch(),
	})
	return
}

func init() {
	RootCmd.AddCommand(killCommand)
}
