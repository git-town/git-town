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

func createHackConfig(args []string, repo *git.ProdRepo) (result appendConfig, err error) {
	result.targetBranch = args[0]
	result.parentBranch, err = determineParentBranch(result.targetBranch, repo)
	if err != nil {
		return result, err
	}
	result.hasOrigin, err = repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	result.shouldNewBranchPush = repo.Config.ShouldNewBranchPush()
	result.isOffline = repo.Config.IsOffline()
	if result.hasOrigin && !repo.Config.IsOffline() {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	hasBranch, err := repo.Silent.HasLocalOrOriginBranch(result.targetBranch)
	if err != nil {
		return result, err
	}
	if hasBranch {
		return result, fmt.Errorf("a branch named %q already exists", result.targetBranch)
	}
	return
}

func init() {
	hackCmd.Flags().BoolVarP(&promptForParent, "prompt", "p", false, "Prompt for the parent branch")
	RootCmd.AddCommand(hackCmd)
}
