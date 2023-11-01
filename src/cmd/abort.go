package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/cli/log"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/hosting"
	"github.com/git-town/git-town/v10/src/hosting/github"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/runstate"
	"github.com/git-town/git-town/v10/src/vm/statefile"
	"github.com/spf13/cobra"
)

const abortDesc = "Aborts the last run git-town command"

func abortCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "abort",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   abortDesc,
		Long:    long(abortDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAbort(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeAbort(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialStashSnapshot, err := determineAbortConfig(repo, verbose)
	if err != nil {
		return err
	}
	abortRunState, exit, err := determineAbortRunstate(config, repo)
	if err != nil || exit {
		return err
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     &repo.Runner,
		Connector:               config.connector,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: config.initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              !config.pushHook,
	})
}

func determineAbortConfig(repo *execute.OpenRepoResult, verbose bool) (*abortConfig, domain.StashSnapshot, error) {
	originURL := repo.Runner.Config.OriginURL()
	hostingService, err := repo.Runner.Config.HostingService()
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.Config.GiteaToken(),
		GithubAPIToken:  github.GetAPIToken(repo.Runner.Config),
		GitlabAPIToken:  repo.Runner.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             log.Printing{},
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
		Verbose:               verbose,
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

func determineAbortRunstate(config *abortConfig, repo *execute.OpenRepoResult) (runstate.RunState, bool, error) {
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return runstate.EmptyRunState(), true, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		fmt.Println(messages.AbortNothingToDo)
		return runstate.EmptyRunState(), true, nil
	}
	abortRunState := runState.CreateAbortRunState()
	wrap(&abortRunState.RunProgram, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranchesSnapshot.Active,
		PreviousBranch:   config.previousBranch,
	})
	return abortRunState, false, nil
}
