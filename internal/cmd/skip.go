package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/skip"
	"github.com/git-town/git-town/v16/internal/validate"
	"github.com/git-town/git-town/v16/internal/vm/statefile"
	"github.com/spf13/cobra"
)

const skipDesc = "Restart the last run Git Town command by skipping the current branch"

func skipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "skip",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   skipDesc,
		Long:    cmdhelpers.Long(skipDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSkip(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSkip(verbose configdomain.Verbose) error {
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
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
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
		return err
	}
	currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		currentBranch, err = repo.Git.CurrentBranch(repo.Backend)
		if err != nil {
			return err
		}
	}
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
	if err != nil {
		return err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: localBranches,
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return err
	}
	runStateOpt, err := statefile.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	runState, hasRunState := runStateOpt.Get()
	if !hasRunState || runState.IsFinished() {
		return errors.New(messages.SkipNothingToDo)
	}
	if unfinishedDetails, hasUnfinishedDetails := runState.UnfinishedDetails.Get(); hasUnfinishedDetails {
		if !unfinishedDetails.CanSkip {
			return errors.New(messages.SkipBranchHasConflicts)
		}
	}
	return skip.Execute(skip.ExecuteArgs{
		Backend:         repo.Backend,
		CommandsCounter: repo.CommandsCounter,
		Config:          validatedConfig,
		Connector:       connector,
		FinalMessages:   repo.FinalMessages,
		Frontend:        repo.Frontend,
		Git:             repo.Git,
		HasOpenChanges:  repoStatus.OpenChanges,
		InitialBranch:   currentBranch,
		RootDir:         repo.RootDir,
		RunState:        runState,
		TestInputs:      dialogTestInputs,
		Verbose:         verbose,
	})
}
