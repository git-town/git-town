package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const newPullRequestDesc = "Creates a new pull request"

const newPullRequestHelp = `
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
where hostname matches what is in your ssh config file.`

func newPullRequestCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   newPullRequestDesc,
		Long:    long(newPullRequestDesc, fmt.Sprintf(newPullRequestHelp, config.CodeHostingDriverKey, config.CodeHostingOriginHostnameKey)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newPullRequest(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func newPullRequest(debug bool) error {
	run, rootDir, isOffline, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      true,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	allBranches, initialBranch, err := execute.LoadBranches(&run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineNewPullRequestConfig(&run, allBranches, initialBranch, isOffline)
	if err != nil {
		return err
	}
	connector, err := hosting.NewConnector(run.Config.GitTown, &run.Backend, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	if connector == nil {
		return hosting.UnsupportedServiceError()
	}
	stepList, err := newPullRequestStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "new-pull-request",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &run, connector, rootDir)
}

type newPullRequestConfig struct {
	BranchesToSync git.BranchesSyncStatus
	InitialBranch  string
	isOffline      bool
	mainBranch     string
}

func determineNewPullRequestConfig(run *git.ProdRunner, allBranches git.BranchesSyncStatus, initialBranch string, isOffline bool) (*newPullRequestConfig, error) {
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	if hasOrigin {
		err := run.Frontend.Fetch()
		if err != nil {
			return nil, err
		}
	}
	err = validate.KnowsBranchAncestors(initialBranch, run.Config.MainBranch(), &run.Backend)
	if err != nil {
		return nil, err
	}
	lineage := run.Config.Lineage()
	branchNamesToSync := lineage.BranchAndAncestors(initialBranch)
	branchesToSync, err := allBranches.Select(branchNamesToSync)
	return &newPullRequestConfig{
		BranchesToSync: branchesToSync,
		InitialBranch:  initialBranch,
		isOffline:      isOffline,
		mainBranch:     run.Config.MainBranch(),
	}, err
}

func newPullRequestStepList(config *newPullRequestConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.BranchesToSync {
		updateBranchSteps(&list, branch, true, config.isOffline, run)
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &run.Backend, config.mainBranch)
	list.Add(&steps.CreateProposalStep{Branch: config.InitialBranch})
	return list.Result()
}
