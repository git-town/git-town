package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func newPullRequestCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   "Creates a new pull request",
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
where hostname matches what is in your ssh config file.`, config.CodeHostingDriverKey, config.CodeHostingOriginHostnameKey),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNewPullRequest(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runNewPullRequest(debug bool) error {
	repo, err := Repo(RepoArgs{
		printBranchNames:     false,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
		validateIsConfigured: true,
		validateIsOnline:     true,
	})
	if err != nil {
		return err
	}
	config, err := determineNewPullRequestConfig(&repo)
	if err != nil {
		return err
	}
	connector, err := hosting.NewConnector(repo.Config, &repo.InternalRepo, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	if connector == nil {
		return hosting.UnsupportedServiceError()
	}
	stepList, err := newPullRequestStepList(config, &repo)
	if err != nil {
		return err
	}
	runState := runstate.New("new-pull-request", stepList)
	return runstate.Execute(runState, &repo, connector)
}

type newPullRequestConfig struct {
	BranchesToSync []string
	InitialBranch  string
	mainBranch     string
}

func determineNewPullRequestConfig(repo *git.PublicRepo) (*newPullRequestConfig, error) {
	hasOrigin, err := repo.HasOrigin()
	if err != nil {
		return nil, err
	}
	if hasOrigin {
		err := repo.Fetch()
		if err != nil {
			return nil, err
		}
	}
	initialBranch, err := repo.CurrentBranch()
	if err != nil {
		return nil, err
	}
	err = validate.KnowsBranchAncestry(initialBranch, repo.Config.MainBranch(), repo)
	if err != nil {
		return nil, err
	}
	return &newPullRequestConfig{
		InitialBranch:  initialBranch,
		BranchesToSync: append(repo.Config.AncestorBranches(initialBranch), initialBranch),
		mainBranch:     repo.Config.MainBranch(),
	}, nil
}

func newPullRequestStepList(config *newPullRequestConfig, repo *git.PublicRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.BranchesToSync {
		updateBranchSteps(&list, branch, true, repo)
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo, config.mainBranch)
	list.Add(&steps.CreateProposalStep{Branch: config.InitialBranch})
	return list.Result()
}
