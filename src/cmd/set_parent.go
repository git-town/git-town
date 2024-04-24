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
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const setParentDesc = "Prompt to set the parent branch for the current branch"

func setParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "set-parent",
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineSetParentConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	err = verifySetParentConfig(config, repo)
	if err != nil {
		return err
	}
	var defaultChoice gitdomain.LocalBranchName
	if config.existingParent != nil {
		defaultChoice = *config.existingParent
	} else {
		defaultChoice = config.mainBranch
	}
	// prompt for the new parent
	outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
		Branch:          config.currentBranch,
		DefaultChoice:   defaultChoice,
		DialogTestInput: config.dialogTestInputs.Next(),
		Lineage:         repo.Runner.Config.FullConfig.Lineage,
		LocalBranches:   initialBranchesSnapshot.Branches.LocalBranches().Names(),
		MainBranch:      config.mainBranch,
	})
	if err != nil {
		return err
	}
	switch outcome {
	case dialog.ParentOutcomeAborted:
		return nil
	case dialog.ParentOutcomePerennialBranch:
		err = repo.Runner.Config.AddToPerennialBranches(config.currentBranch)
	case dialog.ParentOutcomeSelectedParent:
		err = repo.Runner.Config.SetParent(config.currentBranch, *selectedBranch)
	}
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "set-parent",
		DryRun:                false,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            setParentProgram(config),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              &repo.Runner.Config.FullConfig,
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

type setParentConfig struct {
	currentBranch    gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	existingParent   *gitdomain.LocalBranchName
	hasOpenChanges   bool
	mainBranch       gitdomain.LocalBranchName
}

func determineSetParentConfig(repo *execute.OpenRepoResult, verbose bool) (*setParentConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
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
	existingParent := repo.Runner.Config.FullConfig.Lineage.Parent(branchesSnapshot.Active)
	return &setParentConfig{
		currentBranch:    branchesSnapshot.Active,
		dialogTestInputs: dialogTestInputs,
		existingParent:   &existingParent,
		hasOpenChanges:   repoStatus.OpenChanges,
		mainBranch:       repo.Runner.Config.FullConfig.MainBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func verifySetParentConfig(config *setParentConfig, repo *execute.OpenRepoResult) error {
	if repo.Runner.Config.FullConfig.IsMainOrPerennialBranch(config.currentBranch) {
		return fmt.Errorf(messages.SetParentNoFeatureBranch, config.currentBranch)
	}
	return nil
}
