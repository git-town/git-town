package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/spf13/cobra"
)

const abortDesc = "Aborts the last run git-town command"

func abortCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "abort",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   abortDesc,
		Long:    long(abortDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAbort(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executeAbort(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialStashSnapshot, err := determineAbortConfig(&repo)
	if err != nil {
		return err
	}
	abortRunState, err := determineAbortRunstate(config, &repo)
	if err != nil {
		return err
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     &repo.Runner,
		Connector:               config.connector,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: config.initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              !config.pushHook,
	})
}

func determineAbortConfig(repo *execute.OpenRepoResult) (*abortConfig, domain.StashSnapshot, error) {
	originURL := repo.Runner.Config.OriginURL()
	hostingService, err := repo.Runner.Config.HostingService()
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	lineage := repo.Runner.Config.Lineage()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.Config.GiteaToken(),
		GithubAPIToken:  hosting.GetGitHubAPIToken(repo.Runner.Config),
		GitlabAPIToken:  repo.Runner.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), err
	}
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), err
	}
	_, initialBranchesSnapshot, initialStashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, initialStashSnapshot, err
	}
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, initialStashSnapshot, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	return &abortConfig{
		connector:               connector,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		lineage:                 lineage,
		mainBranch:              mainBranch,
		previousBranch:          previousBranch,
		pushHook:                pushHook,
	}, initialStashSnapshot, err
}

type abortConfig struct {
	connector               hosting.Connector
	hasOpenChanges          bool
	initialBranchesSnapshot domain.BranchesSnapshot
	mainBranch              domain.LocalBranchName
	lineage                 config.Lineage
	previousBranch          domain.LocalBranchName
	pushHook                bool
}

func determineAbortRunstate(config *abortConfig, repo *execute.OpenRepoResult) (runstate.RunState, error) {
	runState, err := persistence.Load(repo.RootDir)
	if err != nil {
		return runstate.RunState{}, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return runstate.RunState{}, fmt.Errorf(messages.AbortNothingToDo)
	}
	abortRunState := runState.CreateAbortRunState()
	err = abortRunState.RunSteps.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranchesSnapshot.Active,
		PreviousBranch:   config.previousBranch,
	})
	return abortRunState, err
}
