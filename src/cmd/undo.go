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
	data, initialStashSize, err := determineUndoData(repo.UnvalidatedConfig.Config, repo, verbose)
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
		Config:           data.config,
		HasOpenChanges:   data.hasOpenChanges,
		InitialStashSize: initialStashSize,
		Lineage:          data.config.Lineage,
		RootDir:          repo.RootDir,
		RunState:         runState,
		Runner:           data.prodRunner,
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
	prodRunner              *git.ProdRunner
}

func determineUndoData(unvalidatedConfig configdomain.UnvalidatedConfig, repo *execute.OpenRepoResult, verbose bool) (*undoData, gitdomain.StashSize, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, 0, err
	}
	initialBranchesSnapshot, initialStashSize, _, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                &unvalidatedConfig,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil {
		return nil, initialStashSize, err
	}
	validatedConfig, err := validate.Config(repo.UnvalidatedConfig, gitdomain.LocalBranchNames{}, gitdomain.BranchInfos{}, &repo.Backend, &dialogTestInputs)
	if err != nil {
		return nil, initialStashSize, err
	}
	prodRunner := git.ProdRunner{
		Config:          validatedConfig,
		Backend:         repo.Backend,
		Frontend:        repo.Frontend,
		CommandsCounter: repo.CommandsCounter,
		FinalMessages:   &repo.FinalMessages,
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	var connector hostingdomain.Connector
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          &validatedConfig.Config,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, initialStashSize, err
		}
	}
	return &undoData{
		config:                  validatedConfig.Config,
		connector:               connector,
		dialogTestInputs:        dialogTestInputs,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		previousBranch:          previousBranch,
		prodRunner:              &prodRunner,
	}, initialStashSize, nil
}
