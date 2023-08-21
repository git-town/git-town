package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
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
	repo, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return err
	}
	runState, err := runstate.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return fmt.Errorf(messages.AbortNothingToDo)
	}
	config, err := determineAbortConfig(&repo.Runner)
	if err != nil {
		return err
	}
	abortRunState := runState.CreateAbortRunState()
	if err != nil {
		return err
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &abortRunState,
		Run:       &repo.Runner,
		Connector: config.connector,
		RootDir:   repo.RootDir,
	})
}

func determineAbortConfig(run *git.ProdRunner) (*abortConfig, error) {
	originURL := run.Config.OriginURL()
	hostingService, err := run.Config.HostingService()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetShaForBranch: run.Backend.ShaForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   run.Config.GiteaToken(),
		GithubAPIToken:  hosting.GetGitHubAPIToken(run.Config),
		GitlabAPIToken:  run.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
	return &abortConfig{
		connector: connector,
	}, err
}

type abortConfig struct {
	connector hosting.Connector
}
