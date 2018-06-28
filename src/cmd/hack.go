package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

var promptForParent bool

var hackCmd = &cobra.Command{
	Use:   "hack <branch>",
	Short: "Creates a new feature branch off the main development branch",
	Long: `Creates a new feature branch off the main development branch

Syncs the main branch and forks a new feature branch with the given name off it.

If (and only if) new-branch-push-flag is true,
pushes the new feature branch to the remote repository.

Finally, brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "new-branch-push-flag" configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getHackConfig(args)
		stepList := getAppendStepList(config)
		runState := steps.NewRunState("hack", stepList)
		steps.Run(runState)
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func getParentBranch(targetBranch string) string {
	if promptForParent {
		parentBranch := prompt.AskForBranchParent(targetBranch, git.GetMainBranch())
		prompt.EnsureKnowsParentBranches([]string{parentBranch})
		return parentBranch
	}
	return git.GetMainBranch()
}

func getHackConfig(args []string) (result appendConfig) {
	result.TargetBranch = args[0]
	result.ParentBranch = getParentBranch(result.TargetBranch)
	if git.HasRemote("origin") && !git.IsOffline() {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	return
}

func init() {
	hackCmd.Flags().BoolVarP(&promptForParent, "prompt", "p", false, "Prompt for the parent branch")
	RootCmd.AddCommand(hackCmd)
}
