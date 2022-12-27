package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/userinput"
	"github.com/spf13/cobra"
)

var promptForParent bool

var hackCmd = &cobra.Command{
	Use:   "hack <branch>",
	Short: "Creates a new feature branch off the main development branch",
	Long: `Creates a new feature branch off the main development branch

Syncs the main branch,
forks a new feature branch with the given name off the main branch,
pushes the new feature branch to origin
(if and only if "new-branch-push-flag" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := createHackConfig(args, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := createAppendStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := runstate.New("hack", stepList)
		err = runstate.Execute(runState, prodRepo, nil)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

func determineParentBranch(targetBranch string, repo *git.ProdRepo) (string, error) {
	if promptForParent {
		parentBranch, err := userinput.AskForBranchParent(targetBranch, repo.Config.MainBranch(), repo)
		if err != nil {
			return "", err
		}
		err = userinput.EnsureKnowsParentBranches([]string{parentBranch}, repo)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return repo.Config.MainBranch(), nil
}

func createHackConfig(args []string, repo *git.ProdRepo) (appendConfig, error) {
	targetBranch := args[0]
	parentBranch, err := determineParentBranch(targetBranch, repo)
	if err != nil {
		return appendConfig{}, err
	}
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return appendConfig{}, err
	}
	shouldNewBranchPush := repo.Config.ShouldNewBranchPush()
	isOffline := repo.Config.IsOffline()
	if hasOrigin && !repo.Config.IsOffline() {
		err := repo.Logging.Fetch()
		if err != nil {
			return appendConfig{}, err
		}
	}
	hasBranch, err := repo.Silent.HasLocalOrOriginBranch(targetBranch)
	if err != nil {
		return appendConfig{}, err
	}
	if hasBranch {
		return appendConfig{}, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	result := appendConfig{
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		shouldNewBranchPush: shouldNewBranchPush,
		noPushVerify:        !repo.Config.PushVerify(),
		isOffline:           isOffline,
	}
	return result, nil
}

func init() {
	hackCmd.Flags().BoolVarP(&promptForParent, "prompt", "p", false, "Prompt for the parent branch")
	RootCmd.AddCommand(hackCmd)
}
