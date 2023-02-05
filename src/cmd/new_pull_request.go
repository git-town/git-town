package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type newPullRequestConfig struct {
	BranchesToSync            []string
	branchesWithDeletedRemote []string // local branches whose tracking branches have been deleted
	InitialBranch             string
}

func newPullRequestCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "new-pull-request",
		Short: "Creates a new pull request",
		Long: fmt.Sprintf(`Creates a new pull request

Syncs the current branch
and opens a browser window to the new pull request page of your repository.

The form is pre-populated for the current branch
so that the pull request only shows the changes made
against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, Gitea and Bitbucket.
When using self-hosted versions this command needs to be configured with
"git config %s <driver>"
where driver is "github", "gitlab", "gitea", or "bitbucket".
When using SSH identities, this command needs to be configured with
"git config %s <hostname>"
where hostname matches what is in your ssh config file.`, config.CodeHostingDriver, config.CodeHostingOriginHostname),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := createNewPullRequestConfig(repo)
			if err != nil {
				cli.Exit(err)
			}
			driver := hosting.NewDriver(&repo.Config, &repo.Silent, cli.PrintDriverAction)
			if driver == nil {
				cli.Exit(hosting.UnsupportedServiceError())
			}
			stepList, err := createNewPullRequestStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("new-pull-request", stepList)
			err = runstate.Execute(runState, repo, driver)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			if err := repo.Config.ValidateIsOnline(); err != nil {
				return err
			}
			return nil
		},
	}
}

func createNewPullRequestConfig(repo *git.ProdRepo) (newPullRequestConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return newPullRequestConfig{}, err
	}
	if hasOrigin {
		err := repo.Logging.Fetch()
		if err != nil {
			return newPullRequestConfig{}, err
		}
	}
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return newPullRequestConfig{}, err
	}
	parentDialog := dialog.ParentBranches{}
	err = parentDialog.EnsureKnowsParentBranches([]string{initialBranch}, repo)
	if err != nil {
		return newPullRequestConfig{}, err
	}
	branchesWithDeletedRemote, err := repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return newPullRequestConfig{}, err
	}
	return newPullRequestConfig{
		InitialBranch:             initialBranch,
		BranchesToSync:            append(repo.Config.AncestorBranches(initialBranch), initialBranch),
		branchesWithDeletedRemote: branchesWithDeletedRemote,
	}, nil
}

func createNewPullRequestStepList(config newPullRequestConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchName := range config.BranchesToSync {
		steps, err := updateBranchSteps(branchName, true, config.branchesWithDeletedRemote, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	if err != nil {
		return runstate.StepList{}, err
	}
	result.Append(&steps.CreatePullRequestStep{BranchName: config.InitialBranch})
	return result, nil
}
