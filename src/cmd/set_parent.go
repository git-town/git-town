package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
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
	err = verifySetParentData(data, repo)
	if err != nil {
		return err
	}
	outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
		Branch:          data.currentBranch,
		DefaultChoice:   data.defaultChoice,
		DialogTestInput: data.dialogTestInputs.Next(),
		Lineage:         repo.Runner.Config.FullConfig.Lineage,
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
		Config:                  repo.Runner.Config.FullConfig,
		Connector:               nil,
		DialogTestInputs:        &data.dialogTestInputs,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type setParentData struct {
	currentBranch    gitdomain.LocalBranchName
	defaultChoice    gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	mainBranch       gitdomain.LocalBranchName
}

func determineSetParentData(repo *execute.OpenRepoResult, verbose bool) (*setParentData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, 0, exit, err
	}
	mainBranch := repo.Runner.Config.FullConfig.MainBranch
	existingParent, hasParent := repo.Runner.Config.FullConfig.Lineage.Parent(branchesSnapshot.Active).Get()
	var defaultChoice gitdomain.LocalBranchName
	if hasParent {
		defaultChoice = existingParent
	} else {
		defaultChoice = mainBranch
	}
	return &setParentData{
		currentBranch:    branchesSnapshot.Active,
		defaultChoice:    defaultChoice,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		mainBranch:       mainBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func verifySetParentData(data *setParentData, repo *execute.OpenRepoResult) error {
	if repo.Runner.Config.FullConfig.IsMainOrPerennialBranch(data.currentBranch) {
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
