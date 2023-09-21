package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/spf13/cobra"
)

const continueDesc = "Restarts the last run git-town command after having resolved conflicts"

func continueCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "continue",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   continueDesc,
		Long:    long(continueDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runContinue(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runContinue(debug bool) error {
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
	config, err := determineContinueConfig(&repo)
	if err != nil {
		return err
	}
	runState, err := determineContinueRunstate(&repo)
	if err != nil {
		return err
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:  runState,
		Run:       &repo.Runner,
		Connector: config.connector,
		Lineage:   config.lineage,
		RootDir:   repo.RootDir,
	})
}

func determineContinueConfig(repo *execute.OpenRepoResult) (*continueConfig, error) {
	lineage := repo.Runner.Config.Lineage()
	_, _, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return nil, err
	}
	hasConflicts, err := repo.Runner.Backend.HasConflicts()
	if err != nil {
		return nil, err
	}
	if hasConflicts {
		return nil, fmt.Errorf(messages.ContinueUnresolvedConflicts)
	}
	originURL := repo.Runner.Config.OriginURL()
	hostingService, err := repo.Runner.Config.HostingService()
	if err != nil {
		return nil, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
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
	return &continueConfig{
		connector: connector,
		lineage:   lineage,
	}, err
}

type continueConfig struct {
	connector hosting.Connector
	lineage   config.Lineage
}

func determineContinueRunstate(repo *execute.OpenRepoResult) (*runstate.RunState, error) {
	runState, err := persistence.Load(repo.RootDir)
	if err != nil {
		return nil, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return nil, fmt.Errorf(messages.ContinueNothingToDo)
	}
	return runState, nil
}
