package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/steps"
	"github.com/spf13/cobra"
)

type newPullRequestConfig struct {
	InitialBranch  string
	BranchesToSync []string
}

var newPullRequestCommand = &cobra.Command{
	Use:   "new-pull-request",
	Short: "Creates a new pull request",
	Long: `Creates a new pull request

Syncs the current branch
and opens a browser window to the new pull request page of your repository.

The form is pre-populated for the current branch
so that the pull request only shows the changes made
against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, Gitea and Bitbucket.
When using self-hosted versions this command needs to be configured with
"git config git-town.code-hosting-driver <driver>"
where driver is "github", "gitlab", "gitea", or "bitbucket".
When using SSH identities, this command needs to be configured with
"git config git-town.code-hosting-origin-hostname <hostname>"
where hostname matches what is in your ssh config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := git.NewProdRepo()
		config, err := getNewPullRequestConfig(repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		driver := drivers.Load(repo.Configuration)
		if driver == nil {
			fmt.Println(drivers.UnsupportedHostingError())
			os.Exit(1)
		}
		stepList, err := getNewPullRequestStepList(config, repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		runState := steps.NewRunState("new-pull-request", stepList)
		err = steps.Run(runState, repo, driver)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		if err := validateIsConfigured(); err != nil {
			return err
		}
		if err := git.Config().ValidateIsOnline(); err != nil {
			return err
		}
		return nil
	},
}

func getNewPullRequestConfig(repo *git.ProdRepo) (result newPullRequestConfig, err error) {
	if git.HasRemote("origin") {
		err := repo.Logging.FetchPrune()
		if err != nil {
			return result, err
		}
	}
	result.InitialBranch = git.GetCurrentBranchName()
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	result.BranchesToSync = append(git.Config().GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	return
}

func getNewPullRequestStepList(config newPullRequestConfig, repo *git.ProdRepo) (result steps.StepList, err error) {
	for _, branchName := range config.BranchesToSync {
		steps, err := steps.GetSyncBranchSteps(branchName, true, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	}
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	result.Append(&steps.CreatePullRequestStep{BranchName: config.InitialBranch})
	return result, nil
}

func init() {
	RootCmd.AddCommand(newPullRequestCommand)
}
