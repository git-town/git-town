package cmd

import (
	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/steps"
	"github.com/Originate/git-town/lib/util"
	"github.com/spf13/cobra"
)

type KillFlags struct {
	Undo bool
}

type KillConfig struct {
	InitialBranch string
	TargetBranch  string
}

var killFlags KillFlags

var killCommand = &cobra.Command{
	Use:   "kill",
	Short: "Removes an obsolete feature branch",
	Long:  "Removes an obsolete feature branch",
	Run: func(cmd *cobra.Command, args []string) {
		steps.Run(steps.RunOptions{
			CanSkip:              func() bool { return false },
			Command:              "kill",
			IsAbort:              false,
			IsContinue:           false,
			IsSkip:               false,
			IsUndo:               killFlags.Undo,
			SkipMessageGenerator: func() string { return "" },
			StepListGenerator: func() steps.StepList {
				killConfig := checkKillPreconditions(args)
				return getKillStepList(killConfig)
			},
		})
	},
}

func checkKillPreconditions(args []string) (result KillConfig) {
	result.InitialBranch = git.GetCurrentBranchName()

	if len(args) == 0 {
		result.TargetBranch = result.InitialBranch
	} else if len(args) == 1 {
		result.TargetBranch = args[0]
	} else {
		util.ExitWithErrorMessage("Too many arguments")
	}

	config.EnsureIsFeatureBranch(result.TargetBranch, "You can only kill feature branches.")

	if config.HasLocalBranch(result.TargetBranch) {
		prompt.EnsureKnowsParentBranches([]string{result.TargetBranch})
	}

	if config.HasRemote("origin") {
		steps.FetchStep{}.Run()
	}

	if result.InitialBranch != result.TargetBranch {
		git.EnsureHasBranch(result.TargetBranch)
	}

	return
}

func getKillStepList(killConfig KillConfig) (result steps.StepList) {
	if config.HasLocalBranch(killConfig.TargetBranch) {
		result.Append(DeleteLocalBranch{BranchName: killConfig.TargetBranch, Force: true})
	} else {
		result.Append(DeleteRemoteBranch{BranchName: killConfig.TargetBranch})
	}
	return
}

func init() {
	killCommand.Flags().BoolVar(&killFlags.Undo, "undo", false, "Undo a previous command")
	RootCmd.AddCommand(killCommand)
}
