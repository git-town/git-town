package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/flags"
	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v19/internal/config"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/execute"
	"github.com/git-town/git-town/v19/internal/forge"
	"github.com/git-town/git-town/v19/internal/forge/forgedomain"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/undo"
	"github.com/git-town/git-town/v19/internal/validate"
	"github.com/git-town/git-town/v19/internal/vm/statefile"
	. "github.com/git-town/git-town/v19/pkg/prelude"
	"github.com/spf13/cobra"
)

const undoDesc = "Undo the most recent Git Town command"

func undoCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: cmdhelpers.GroupIDErrors,
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    cmdhelpers.Long(undoDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeUndo(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUndo(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineUndoData(repo, verbose)
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
	fmt.Println("222222222222222222222222222222222222222222222")
	fmt.Println(runState.EndBranchesSnapshot)
	fmt.Println("333333333333333333333333333333333333333333333")
	fmt.Println(runState.BeginBranchesSnapshot)
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return err
	}
	return undo.Execute(undo.ExecuteArgs{
		Backend:          repo.Backend,
		CommandsCounter:  repo.CommandsCounter,
		Config:           data.config,
		Connector:        connector,
		FinalMessages:    repo.FinalMessages,
		Frontend:         repo.Frontend,
		Git:              repo.Git,
		HasOpenChanges:   data.hasOpenChanges,
		InitialStashSize: data.stashSize,
		RootDir:          repo.RootDir,
		RunState:         runState,
		Verbose:          verbose,
	})
}

type undoData struct {
	config                  config.ValidatedConfig
	connector               Option[forgedomain.Connector]
	dialogTestInputs        components.TestInputs
	hasOpenChanges          bool
	initialBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch          Option[gitdomain.LocalBranchName]
	stashSize               gitdomain.StashSize
}

func determineUndoData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data undoData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, false, err
	}
	fmt.Println("111111111111111111111111111111111111111111111111")
	fmt.Println(branchesSnapshot.Branches)
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	return undoData{
		config:                  validatedConfig,
		connector:               connector,
		dialogTestInputs:        dialogTestInputs,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: branchesSnapshot,
		previousBranch:          previousBranch,
		stashSize:               stashSize,
	}, false, nil
}
