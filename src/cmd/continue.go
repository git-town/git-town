package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/hosting"
	"github.com/git-town/git-town/v12/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/interpreter"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/statefile"
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineContinueConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	runState, exit, err := determineContinueRunstate(repo)
	if err != nil || exit {
		return err
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		Connector:               config.connector,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

func determineContinueConfig(repo *execute.OpenRepoResult, verbose bool) (*continueConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	initialBranchesSnapshot, initialStashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FullConfig:            &repo.Runner.FullConfig,
		HandleUnfinishedState: false,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, initialBranchesSnapshot, initialStashSize, exit, err
	}
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, initialBranchesSnapshot, initialStashSize, false, err
	}
	if repoStatus.Conflicts {
		return nil, initialBranchesSnapshot, initialStashSize, false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	if repoStatus.UntrackedChanges {
		return nil, initialBranchesSnapshot, initialStashSize, false, errors.New(messages.ContinueUntrackedChanges)
	}
	originURL := repo.Runner.Config.OriginURL()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:      &repo.Runner.FullConfig,
		HostingPlatform: repo.Runner.HostingPlatform,
		Log:             print.Logger{},
		OriginURL:       originURL,
	})
	return &continueConfig{
		FullConfig:       &repo.Runner.FullConfig,
		connector:        connector,
		dialogTestInputs: dialogTestInputs,
	}, initialBranchesSnapshot, initialStashSize, false, err
}

type continueConfig struct {
	connector hostingdomain.Connector
	*configdomain.FullConfig
	dialogTestInputs components.TestInputs
}

func determineContinueRunstate(repo *execute.OpenRepoResult) (runstate.RunState, bool, error) {
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return runstate.EmptyRunState(), true, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsFinished() {
		fmt.Println(messages.ContinueNothingToDo)
		return runstate.EmptyRunState(), true, nil
	}
	return *runState, false, nil
}
