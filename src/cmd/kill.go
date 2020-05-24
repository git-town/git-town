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
		repo := git.ProdRepoInCurrentDir()
		config, err := getKillConfig(args, &repo.Silent)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		stepList, err := getKillStepList(config, &repo.Silent)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		runState := steps.NewRunState("kill", stepList)
		err = steps.Run(runState)
		if err != nil {
			fmt.Println("Error:", err)
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

func getKillConfig(args []string, runner *git.Runner) (result killConfig, err error) {
	result.InitialBranch, err = runner.CurrentBranch()
	if err != nil {
		return result, err
	}
	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
	} else {
		result.TargetBranch = args[0]
	}
	if !runner.IsFeatureBranch(result.TargetBranch) {
		return result, fmt.Errorf("you can only kill feature branches")
	}
	result.IsTargetBranchLocal, err = runner.HasLocalBranch(result.TargetBranch)
	if err != nil {
		return result, err
	}
	if result.IsTargetBranchLocal {
		prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
		runner.Configuration.Reload()
	}
	hasOrigin, err := runner.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if hasOrigin && !runner.IsOffline() {
		err := script.Fetch()
		if err != nil {
			return result, err
		}
	}
	if result.InitialBranch != result.TargetBranch {
		hasTargetBranch, err := runner.HasLocalOrRemoteBranch(result.TargetBranch)
		if err != nil {
			return result, err
		}
		if !hasTargetBranch {
			return result, fmt.Errorf("there is no branch named %q", result.TargetBranch)
		}
	}
	return result, nil
}

func getKillStepList(config killConfig, runner *git.Runner) (result steps.StepList, err error) {
	switch {
	case config.IsTargetBranchLocal:
		hasTrackingBranch, err := runner.HasTrackingBranch(config.TargetBranch)
		if err != nil {
			return result, err
		}
		if hasTrackingBranch && !runner.IsOffline() {
			result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: true})
		}
		targetBranchParent := runner.GetParentBranch(config.TargetBranch)
		if config.InitialBranch == config.TargetBranch {
			hasOpenChanges, err := runner.HasOpenChanges()
			if err != nil {
				return result, err
			}
			if hasOpenChanges {
				result.Append(&steps.CommitOpenChangesStep{})
			}
			result.Append(&steps.CheckoutBranchStep{BranchName: targetBranchParent})
		}
		result.Append(&steps.DeleteLocalBranchStep{BranchName: config.TargetBranch, Force: true})
		for _, child := range runner.GetChildBranches(config.TargetBranch) {
			result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: targetBranchParent})
		}
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.TargetBranch})
	case !runner.IsOffline():
		result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.TargetBranch, IsTracking: false})
	default:
		fmt.Printf("Cannot delete remote branch %q in offline mode", config.TargetBranch)
		os.Exit(1)
	}
	previousBranch, err := runner.PreviouslyCheckedOutBranch()
	if err != nil {
		return result, err
	}
	result.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.InitialBranch != config.TargetBranch && config.TargetBranch == previousBranch,
	})
	return result, nil
}

func init() {
	RootCmd.AddCommand(killCommand)
}
