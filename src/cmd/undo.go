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
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/statefile"
	"github.com/spf13/cobra"
)

const undoDesc = "Undoes the most recent Git Town command"

func undoCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    cmdhelpers.Long(undoDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUndo(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUndo(verbose bool) error {
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
	var config *undoConfig
	var initialStashSize gitdomain.StashSize
	config, initialStashSize, repo.Runner.Lineage, err = determineUndoConfig(repo, verbose)
	if err != nil {
		return err
	}
	undoRunState, err := determineUndoRunState(config, repo)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		FullConfig:              config.FullConfig,
		RunState:                &undoRunState,
		Run:                     repo.Runner,
		Connector:               config.connector,
		DialogTestInputs:        &config.dialogTestInputs,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: config.initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
	})
}

type undoConfig struct {
	*configdomain.FullConfig
	connector               hostingdomain.Connector
	dialogTestInputs        components.TestInputs
	hasOpenChanges          bool
	initialBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch          gitdomain.LocalBranchName
}

func determineUndoConfig(repo *execute.OpenRepoResult, verbose bool) (*undoConfig, gitdomain.StashSize, configdomain.Lineage, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	initialBranchesSnapshot, initialStashSize, _, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return nil, initialStashSize, repo.Runner.Lineage, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, initialStashSize, repo.Runner.Lineage, err
	}
	originURL := repo.Runner.Config.OriginURL()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:      &repo.Runner.FullConfig,
		HostingPlatform: repo.Runner.Config.HostingPlatform,
		OriginURL:       originURL,
		Log:             print.Logger{},
	})
	if err != nil {
		return nil, initialStashSize, repo.Runner.Lineage, err
	}
	return &undoConfig{
		FullConfig:              &repo.Runner.FullConfig,
		connector:               connector,
		dialogTestInputs:        dialogTestInputs,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		previousBranch:          previousBranch,
	}, initialStashSize, repo.Runner.Lineage, nil
}

func determineUndoRunState(config *undoConfig, repo *execute.OpenRepoResult) (runstate.RunState, error) {
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return runstate.EmptyRunState(), fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil {
		return runstate.EmptyRunState(), errors.New(messages.UndoNothingToDo)
	}
	var undoRunState runstate.RunState
	if runState.IsFinished() {
		undoRunState = runState.CreateUndoRunState()
	} else {
		undoRunState = runState.CreateAbortRunState()
	}
	if !undoRunState.DryRun {
		cmdhelpers.Wrap(&undoRunState.RunProgram, cmdhelpers.WrapOptions{
			DryRun:                   undoRunState.DryRun,
			RunInGitRoot:             true,
			StashOpenChanges:         config.hasOpenChanges,
			PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
		})
		// If the command to undo failed and was continued,
		// there might be opcodes in the undo stack that became obsolete
		// when the command was continued.
		// Example: the command stashed away uncommitted changes,
		// failed, and remembered in the undo list to pop the stack.
		// When continuing, it finishes and pops the stack as part of the continue list.
		// When we run undo now, it still wants to pop the stack even though that was already done.
		// This seems to apply only to popping the stack and switching back to the initial branch.
		// Hence we consolidate these opcode types here.
		undoRunState.RunProgram = undoRunState.RunProgram.
			MoveToEnd(&opcodes.RestoreOpenChanges{}).
			RemoveAllButLast("*opcodes.CheckoutIfExists")
	}
	return undoRunState, err
}
