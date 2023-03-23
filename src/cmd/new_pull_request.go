package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/flags"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
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
	run, exit, err := LoadProdRunner(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: true,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
		validateIsOnline:      true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineNewPullRequestConfig(&run)
	if err != nil {
		return err
	}
	connector, err := hosting.NewConnector(run.Config, &run.Backend, cli.PrintConnectorAction)
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
	runState := runstate.New("new-pull-request", stepList)
	return runstate.Execute(runState, &run, connector)
}

type newPullRequestConfig struct {
	BranchesToSync []string
	InitialBranch  string
	mainBranch     string
}

func determineNewPullRequestConfig(run *git.ProdRunner) (*newPullRequestConfig, error) {
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
	initialBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return nil, err
	}
	err = validate.KnowsBranchAncestry(initialBranch, run.Config.MainBranch(), &run.Backend)
	if err != nil {
		return nil, err
	}
	return &newPullRequestConfig{
		InitialBranch:  initialBranch,
		BranchesToSync: append(run.Config.AncestorBranches(initialBranch), initialBranch),
		mainBranch:     run.Config.MainBranch(),
	}, nil
}

func newPullRequestStepList(config *newPullRequestConfig, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.BranchesToSync {
		updateBranchSteps(&list, branch, true, run)
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &run.Backend, config.mainBranch)
	list.Add(&steps.CreateProposalStep{Branch: config.InitialBranch})
	return list.Result()
}
