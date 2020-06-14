package cmd

import (
	"github.com/git-town/git-town/src/cli"
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
		config, err := getNewPullRequestConfig(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		driver := drivers.Load(prodRepo.Configuration, &prodRepo.Silent)
		if driver == nil {
			cli.Exit(drivers.UnsupportedHostingError())
		}
		stepList, err := getNewPullRequestStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := steps.NewRunState("new-pull-request", stepList)
		err = steps.Run(runState, prodRepo, driver)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		if err := validateIsConfigured(prodRepo); err != nil {
			return err
		}
		if err := prodRepo.ValidateIsOnline(); err != nil {
			return err
		}
		return nil
	},
}

func getNewPullRequestConfig(repo *git.ProdRepo) (result newPullRequestConfig, err error) {
	hasOrigin, err := repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if hasOrigin {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	result.InitialBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return result, err
	}
	err = prompt.EnsureKnowsParentBranches([]string{result.InitialBranch}, repo)
	if err != nil {
		return result, err
	}
	result.BranchesToSync = append(repo.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
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
	err = result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	if err != nil {
		return result, err
	}
	result.Append(&steps.CreatePullRequestStep{BranchName: config.InitialBranch})
	return result, nil
}

func init() {
	RootCmd.AddCommand(newPullRequestCommand)
}
