package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const switchDesc = "Displays the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   switchDesc,
		Long:    cmdhelpers.Long(switchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSwitch(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwitch(verbose bool) error {
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
	config, exit, err := determineSwitchConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	branchToCheckout, abort, err := dialog.SwitchBranch(config.branchNames, config.initialBranch, repo.Runner.Lineage)
	if err != nil || abort {
		return err
	}
	if branchToCheckout == config.initialBranch {
		return nil
	}
	err = repo.Runner.Frontend.CheckoutBranch(branchToCheckout)
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
	branchNames   gitdomain.LocalBranchNames
	initialBranch gitdomain.LocalBranchName
}

func determineSwitchConfig(repo *execute.OpenRepoResult, verbose bool) (*switchConfig, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FullConfig:            &repo.Runner.FullConfig,
		HandleUnfinishedState: true,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	return &switchConfig{
		branchNames:   branchesSnapshot.Branches.Names(),
		initialBranch: branchesSnapshot.Active,
	}, false, err
}
