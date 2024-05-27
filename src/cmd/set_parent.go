package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const setParentCmd = "set-parent"

const setParentDesc = "Prompt to set the parent branch for the current branch"

func setParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     setParentCmd,
		GroupID: "lineage",
		Args:    cobra.NoArgs,
		Short:   setParentDesc,
		Long:    cmdhelpers.Long(setParentDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSetParent(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(verbose bool) error {
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
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineSetParentData(repo, verbose)
	if err != nil || exit {
		return err
	}
	err = verifySetParentData(data)
	if err != nil {
		return err
	}
	outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
		Branch:          data.currentBranch,
		DefaultChoice:   data.defaultChoice,
		DialogTestInput: data.dialogTestInputs.Next(),
		Lineage:         data.config.Config.Lineage,
		LocalBranches:   initialBranchesSnapshot.Branches.LocalBranches().Names(),
		MainBranch:      data.mainBranch,
	})
	if err != nil {
		return err
	}
	prog, aborted := setParentProgram(outcome, selectedBranch, data.currentBranch)
	if aborted {
		return nil
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               setParentCmd,
		DryRun:                false,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            prog,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type setParentData struct {
	config           config.ValidatedConfig
	currentBranch    gitdomain.LocalBranchName
	defaultChoice    gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	mainBranch       gitdomain.LocalBranchName
}

func emptySetParentData() setParentData {
	return setParentData{} //exhaustruct:ignore
}

func determineSetParentData(repo execute.OpenRepoResult, verbose bool) (setParentData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return emptySetParentData(), gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return emptySetParentData(), branchesSnapshot, stashSize, exit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: localBranches,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return emptySetParentData(), branchesSnapshot, stashSize, exit, err
	}
	mainBranch := validatedConfig.Config.MainBranch
	currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		return emptySetParentData(), branchesSnapshot, stashSize, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	existingParent, hasParent := validatedConfig.Config.Lineage.Parent(currentBranch).Get()
	var defaultChoice gitdomain.LocalBranchName
	if hasParent {
		defaultChoice = existingParent
	} else {
		defaultChoice = mainBranch
	}
	return setParentData{
		config:           validatedConfig,
		currentBranch:    currentBranch,
		defaultChoice:    defaultChoice,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		mainBranch:       mainBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func verifySetParentData(data setParentData) error {
	if data.config.Config.IsMainOrPerennialBranch(data.currentBranch) {
		return fmt.Errorf(messages.SetParentNoFeatureBranch, data.currentBranch)
	}
	return nil
}

func setParentProgram(outcome dialog.ParentOutcome, selectedBranch, currentBranch gitdomain.LocalBranchName) (result program.Program, aborted bool) {
	switch outcome {
	case dialog.ParentOutcomeAborted:
		return result, true
	case dialog.ParentOutcomePerennialBranch:
		result.Add(&opcodes.AddToPerennialBranches{
			Branch: currentBranch,
		})
		result.Add(&opcodes.DeleteParentBranch{
			Branch: currentBranch,
		})
	case dialog.ParentOutcomeSelectedParent:
		result.Add(&opcodes.SetParent{
			Branch: currentBranch,
			Parent: selectedBranch,
		})
	}
	return result, false
}
