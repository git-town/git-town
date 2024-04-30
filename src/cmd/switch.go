package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/validate"
	"github.com/spf13/cobra"
)

const switchDesc = "Display the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMergeFlag, readMergeFlag := flags.SwitchMerge()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   switchDesc,
		Long:    cmdhelpers.Long(switchDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSwitch(readVerboseFlag(cmd), readMergeFlag(cmd))
		},
	}
	addMergeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwitch(verbose, merge bool) error {
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
	config, initialBranches, exit, err := determineSwitchConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	branchToCheckout, abort, err := dialog.SwitchBranch(config.branchNames, config.initialBranch, config.Lineage, initialBranches.Branches, config.uncommittedChanges, config.dialogInputs.Next())
	if err != nil || abort {
		return err
	}
	if branchToCheckout == config.initialBranch {
		return nil
	}
	err = repo.Frontend.CheckoutBranch(branchToCheckout, merge)
	if err != nil {
		exitCode := 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		os.Exit(exitCode)
	}
	return nil
}

type switchConfig struct {
	branchNames        gitdomain.LocalBranchNames
	dialogInputs       components.TestInputs
	initialBranch      gitdomain.LocalBranchName
	Lineage            configdomain.Lineage
	uncommittedChanges bool
}

func determineSwitchConfig(repo *execute.OpenRepoResult, verbose bool) (*switchConfig, gitdomain.BranchesSnapshot, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.BackendCommands.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                &repo.UnvalidatedConfig.Config,
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
		return nil, branchesSnapshot, exit, err
	}
	validatedConfig, err := validate.Config(repo.UnvalidatedConfig, branchesSnapshot.Branches.LocalBranches().Names())
	return &switchConfig{
		branchNames:        branchesSnapshot.Branches.Names(),
		dialogInputs:       dialogTestInputs,
		initialBranch:      branchesSnapshot.Active,
		uncommittedChanges: repoStatus.UntrackedChanges,
	}, branchesSnapshot, false, err
}
