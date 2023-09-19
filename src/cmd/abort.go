package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/undo"
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
			return abort(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func abort(debug bool) error {
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
	runState, err := persistence.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return fmt.Errorf(messages.AbortNothingToDo)
	}
	config, initialBranchesSnapshot, err := determineAbortConfig(&repo.Runner)
	if err != nil {
		return err
	}
	abortRunState := runState.CreateAbortRunState()
	if err != nil {
		return err
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     &repo.Runner,
		Connector:               config.connector,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
	})
}

func determineAbortConfig(run *git.ProdRunner) (*abortConfig, undo.BranchesSnapshot, error) {
	originURL := run.Config.OriginURL()
	hostingService, err := run.Config.HostingService()
	if err != nil {
		return nil, undo.EmptyBranchesSnapshot(), err
	}
	mainBranch := run.Config.MainBranch()
	lineage := run.Config.Lineage()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: run.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   run.Config.GiteaToken(),
		GithubAPIToken:  hosting.GetGitHubAPIToken(run.Config),
		GitlabAPIToken:  run.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
	if err != nil {
		return nil, undo.EmptyBranchesSnapshot(), err
	}
	allBranches, initialBranch, err := run.Backend.BranchInfos()
	if err != nil {
		return nil, undo.EmptyBranchesSnapshot(), err
	}
	branchesSnapshot := undo.BranchesSnapshot{
		Branches: allBranches,
		Active:   initialBranch,
	}
	return &abortConfig{
		connector: connector,
		lineage:   lineage,
	}, branchesSnapshot, err
}

type abortConfig struct {
	connector hosting.Connector
	lineage   config.Lineage
}
