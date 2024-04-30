package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo"
	"github.com/git-town/git-town/v14/src/validate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
	"github.com/spf13/cobra"
)

const undoDesc = "Undo the most recent Git Town command"

func undoCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    cmdhelpers.Long(undoDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeUndo(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUndo(verbose bool) error {
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
	var config *undoConfig
	var initialStashSize gitdomain.StashSize
	validatedConfig, err := validate.ValidateConfig(repo.UnvalidatedConfig)
	prodRunner := git.ProdRunner{
		Config:          validatedConfig,
		Backend:         repo.BackendCommands,
		Frontend:        repo.Frontend,
		CommandsCounter: repo.CommandsCounter,
		FinalMessages:   &repo.FinalMessages,
	}
	config, initialStashSize, validatedConfig.FullConfig.Lineage, err = determineUndoConfig(repo.UnvalidatedConfig.Config, repo, &prodRunner, verbose)
	if err != nil {
		return err
	}
	runStateOpt, err := statefile.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	runState, hasRunState := runStateOpt.Get()
	if !hasRunState {
		fmt.Println(messages.UndoNothingToDo)
		return nil
	}
	return undo.Execute(undo.ExecuteArgs{
		Config:           validatedConfig.FullConfig,
		HasOpenChanges:   config.hasOpenChanges,
		InitialStashSize: initialStashSize,
		Lineage:          validatedConfig.FullConfig.Lineage,
		RootDir:          repo.RootDir,
		RunState:         runState,
		Runner:           &prodRunner,
		Verbose:          verbose,
	})
}

type undoConfig struct {
	config                  *configdomain.ValidatedConfig
	connector               hostingdomain.Connector
	dialogTestInputs        components.TestInputs
	hasOpenChanges          bool
	initialBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch          gitdomain.LocalBranchName
}

func determineUndoConfig(unvalidatedConfig configdomain.UnvalidatedConfig, repo *execute.OpenRepoResult, runner *git.ProdRunner, verbose bool) (*undoConfig, gitdomain.StashSize, configdomain.Lineage, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := runner.Backend.RepoStatus()
	if err != nil {
		return nil, 0, runner.Config.FullConfig.Lineage, err
	}
	initialBranchesSnapshot, initialStashSize, _, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                &unvalidatedConfig,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil {
		return nil, initialStashSize, runner.Config.FullConfig.Lineage, err
	}
	previousBranch := runner.Backend.PreviouslyCheckedOutBranch()
	var connector hostingdomain.Connector
	if originURL, hasOriginURL := runner.Config.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          &runner.Config.FullConfig,
			HostingPlatform: runner.Config.FullConfig.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, initialStashSize, runner.Config.FullConfig.Lineage, err
		}
	}
	return &undoConfig{
		config:                  &runner.Config.FullConfig,
		connector:               connector,
		dialogTestInputs:        dialogTestInputs,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		previousBranch:          previousBranch,
	}, initialStashSize, runner.Config.FullConfig.Lineage, nil
}
