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
	data, initialStashSize, exit, err := determineUndoData(repo.UnvalidatedConfig.Config, repo, verbose)
	if err != nil || exit {
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
		Backend:          repo.Backend,
		CommandsCounter:  repo.CommandsCounter,
		Config:           data.config,
		FinalMessages:    &repo.FinalMessages,
		Frontend:         repo.Frontend,
		HasOpenChanges:   data.hasOpenChanges,
		InitialStashSize: initialStashSize,
		Lineage:          data.config.Lineage,
		RootDir:          repo.RootDir,
		RunState:         runState,
		Verbose:          verbose,
	})
}

type undoData struct {
	config                  configdomain.ValidatedConfig
	connector               hostingdomain.Connector
	dialogTestInputs        components.TestInputs
	hasOpenChanges          bool
	initialBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch          gitdomain.LocalBranchName
}

func determineUndoData(unvalidatedConfig *configdomain.UnvalidatedConfig, repo *execute.OpenRepoResult, verbose bool) (*undoData, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, 0, false, err
	}
	initialBranchesSnapshot, initialStashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               &repo.Backend,
		Config:                unvalidatedConfig,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		Frontend:              &repo.Frontend,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, initialStashSize, false, err
	}
	validatedConfig, abort, err := validate.Config(validate.ConfigArgs{
		Backend:            &repo.Backend,
		BranchesSnapshot:   initialBranchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		CommandsCounter:    repo.CommandsCounter,
		ConfigSnapshot:     repo.ConfigSnapshot,
		DialogTestInputs:   dialogTestInputs,
		FinalMessages:      &repo.FinalMessages,
		Frontend:           repo.Frontend,
		LocalBranches:      gitdomain.LocalBranchNames{},
		RepoStatus:         repoStatus,
		RootDir:            repo.RootDir,
		StashSize:          initialStashSize,
		TestInputs:         &dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
		Verbose:            verbose,
	})
	if err != nil || abort {
		return nil, initialStashSize, abort, err
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	var connector hostingdomain.Connector
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          repo.UnvalidatedConfig.Config,
			HostingPlatform: repo.UnvalidatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, initialStashSize, false, err
		}
	}
	return &undoData{
		config:                  validatedConfig.Config,
		connector:               connector,
		dialogTestInputs:        dialogTestInputs,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		previousBranch:          previousBranch,
	}, initialStashSize, false, nil
}
