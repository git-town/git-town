package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
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
		Long:    long(newPullRequestDesc, fmt.Sprintf(newPullRequestHelp, config.KeyCodeHostingDriver, config.KeyCodeHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newPullRequest(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func newPullRequest(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: true,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determineNewPullRequestConfig(&repo)
	if err != nil || exit {
		return err
	}
	if err != nil {
		return err
	}
	stepList, err := newPullRequestStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "new-pull-request",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: config.connector,
		Lineage:   config.lineage,
		RootDir:   repo.RootDir,
	})
}

type newPullRequestConfig struct {
	branches           domain.Branches
	branchesToSync     domain.BranchInfos
	connector          hosting.Connector
	hasOpenChanges     bool
	remotes            domain.Remotes
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	previousBranch     domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushHook           bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineNewPullRequestConfig(repo *execute.OpenRepoResult) (*newPullRequestConfig, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	branches, exit, err := execute.LoadSnapshot(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := repo.Runner.Backend.HasOpenChanges()
	if err != nil {
		return nil, false, err
	}
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, false, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	updated, err := validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		DefaultBranch: mainBranch,
		Backend:       &repo.Runner.Backend,
		AllBranches:   branches.All,
		Lineage:       lineage,
		BranchTypes:   branches.Types,
		MainBranch:    mainBranch,
	})
	if err != nil {
		return nil, false, err
	}
	if updated {
		lineage = repo.Runner.Config.Lineage()
	}
	syncStrategy, err := repo.Runner.Config.SyncStrategy()
	if err != nil {
		return nil, false, err
	}
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, false, err
	}
	pullBranchStrategy, err := repo.Runner.Config.PullBranchStrategy()
	if err != nil {
		return nil, false, err
	}
	shouldSyncUpstream, err := repo.Runner.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, false, err
	}
	originURL := repo.Runner.Config.OriginURL()
	hostingService, err := repo.Runner.Config.HostingService()
	if err != nil {
		return nil, false, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.Config.GiteaToken(),
		GithubAPIToken:  repo.Runner.Config.GitHubToken(),
		GitlabAPIToken:  repo.Runner.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
	if err != nil {
		return nil, false, err
	}
	if connector == nil {
		return nil, false, hosting.UnsupportedServiceError()
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync, err := branches.All.Select(branchNamesToSync)
	return &newPullRequestConfig{
		branches:           branches,
		branchesToSync:     branchesToSync,
		connector:          connector,
		hasOpenChanges:     hasOpenChanges,
		remotes:            remotes,
		isOffline:          repo.IsOffline,
		lineage:            lineage,
		mainBranch:         mainBranch,
		previousBranch:     previousBranch,
		pullBranchStrategy: pullBranchStrategy,
		pushHook:           pushHook,
		shouldSyncUpstream: shouldSyncUpstream,
		syncStrategy:       syncStrategy,
	}, false, err
}

func newPullRequestStepList(config *newPullRequestConfig) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.branchesToSync {
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             branch,
			branchTypes:        config.branches.Types,
			remotes:            config.remotes,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
	}
	list.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	list.Add(&steps.CreateProposalStep{Branch: config.branches.Initial})
	return list.Result()
}
