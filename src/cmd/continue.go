package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
	"github.com/spf13/cobra"
)

const continueDesc = "Restart the last run Git Town command after having resolved conflicts"

func continueCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "continue",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   continueDesc,
		Long:    cmdhelpers.Long(continueDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeContinue(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeContinue(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
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
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               config.connector,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		HasOpenChanges:          config.hasOpenChanges,
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
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	initialBranchesSnapshot, initialStashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, initialBranchesSnapshot, initialStashSize, exit, err
	}
	if repoStatus.Conflicts {
		return nil, initialBranchesSnapshot, initialStashSize, false, errors.New(messages.ContinueUnresolvedConflicts)
	}
	if repoStatus.UntrackedChanges {
		return nil, initialBranchesSnapshot, initialStashSize, false, errors.New(messages.ContinueUntrackedChanges)
	}
	originURL := repo.Runner.Config.OriginURL()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:      &repo.Runner.Config.FullConfig,
		HostingPlatform: repo.Runner.Config.FullConfig.HostingPlatform,
		Log:             print.Logger{},
		OriginURL:       originURL,
	})
	return &continueConfig{
		FullConfig:       &repo.Runner.Config.FullConfig,
		connector:        connector,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
	}, initialBranchesSnapshot, initialStashSize, false, err
}

type continueConfig struct {
	connector Option[hostingdomain.Connector]
	*configdomain.FullConfig
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
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
	runState.AbortProgram = program.Program{}
	return *runState, false, nil
}
