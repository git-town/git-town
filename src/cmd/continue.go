package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/hosting"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/git-town/git-town/v11/src/vm/statefile"
	"github.com/spf13/cobra"
)

const continueDesc = "Restarts the last run git-town command after having resolved conflicts"

func continueCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "continue",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   continueDesc,
		Long:    cmdhelpers.Long(continueDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeContinue(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeContinue(verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineContinueConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	runState, exit, err := determineContinueRunstate(repo)
	if err != nil || exit {
		return err
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		FullConfig:              config.FullConfig,
		RunState:                &runState,
		Run:                     repo.Runner,
		Connector:               config.connector,
		DialogTestInputs:        &config.dialogTestInputs,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

func determineContinueConfig(repo *execute.OpenRepoResult, verbose bool) (*continueConfig, gitdomain.BranchesStatus, gitdomain.StashSize, bool, error) {
	initialBranchesSnapshot, initialStashSnapshot, dialogTestInputs, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, initialBranchesSnapshot, initialStashSnapshot, exit, err
	}
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, initialBranchesSnapshot, initialStashSnapshot, false, err
	}
	if repoStatus.Conflicts {
		return nil, initialBranchesSnapshot, initialStashSnapshot, false, fmt.Errorf(messages.ContinueUnresolvedConflicts)
	}
	if repoStatus.UntrackedChanges {
		return nil, initialBranchesSnapshot, initialStashSnapshot, false, fmt.Errorf(messages.ContinueUntrackedChanges)
	}
	originURL := repo.Runner.Config.OriginURL()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:      &repo.Runner.FullConfig,
		HostingPlatform: repo.Runner.HostingPlatform,
		OriginURL:       originURL,
		Log:             print.Logger{},
	})
	return &continueConfig{
		connector:        connector,
		FullConfig:       &repo.Runner.FullConfig,
		dialogTestInputs: dialogTestInputs,
	}, initialBranchesSnapshot, initialStashSnapshot, false, err
}

type continueConfig struct {
	connector hostingdomain.Connector
	*configdomain.FullConfig
	dialogTestInputs dialogcomponents.TestInputs
}

func determineContinueRunstate(repo *execute.OpenRepoResult) (runstate.RunState, bool, error) {
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return runstate.EmptyRunState(), true, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		fmt.Println(messages.ContinueNothingToDo)
		return runstate.EmptyRunState(), true, nil
	}
	return *runState, false, nil
}
