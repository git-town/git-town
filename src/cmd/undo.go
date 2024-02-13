package cmd

import (
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
	"github.com/git-town/git-town/v12/src/undo/undobranches"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/undo/undostash"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/shared"
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
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}

	undoProgram := program.Program{}
	undoProgram.AddProgram(runState.AbortProgram)

	// undo branch changes
	branchSpans := undobranches.NewBranchSpans(runState.BeforeBranchesSnapshot, runState.AfterBranchesSnapshot)
	branchChanges := branchSpans.Changes()
	branchUndoProgram := branchChanges.UndoProgram(undobranches.BranchChangesUndoProgramArgs{
		Config:                   config.FullConfig,
		FinalBranch:              runState.AfterBranchesSnapshot.Active,
		InitialBranch:            runState.BeforeBranchesSnapshot.Active,
		UndoablePerennialCommits: []gitdomain.SHA{},
	})
	undoProgram.AddProgram(branchUndoProgram)

	// undo config changes
	configSpans := undoconfig.NewConfigDiffs(runState.BeforeConfigSnapshot, runState.AfterConfigSnapshot)
	configUndoProgram := configSpans.UndoProgram()
	undoProgram.AddProgram(configUndoProgram)

	// undo stash changes
	stashDiff := undostash.NewStashDiff(runState.BeforeStashSize, initialStashSize)
	undoStashProgram := stashDiff.Program()
	undoProgram.AddProgram(undoStashProgram)

	undoProgram.AddProgram(runState.FinalUndoProgram)

	cmdhelpers.Wrap(&undoProgram, cmdhelpers.WrapOptions{
		DryRun:                   runState.DryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})

	// execute the undo program
	for _, opcode := range branchUndoProgram {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         config.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          repo.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}
	}

	err = statefile.Delete(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), repo.Runner.FinalMessages.Result())

	return nil
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
