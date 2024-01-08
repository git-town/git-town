package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
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
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, err := determineSwitchConfig(repo)
	if err != nil {
		return err
	}
	branchNameToCheckout, err := dialog.SwitchBranchesDialog(config.branchNames.Strings(), config.initialBranch.String())
	if err != nil {
		return err
	}
	if branchNameToCheckout == config.initialBranch.String() {
		return nil
	}
	branchToCheckout := gitdomain.LocalBranchName(branchNameToCheckout)
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

func determineSwitchConfig(repo *execute.OpenRepoResult) (*switchConfig, error) {
	branchNames, current, err := repo.Runner.Backend.LocalBranchNames()
	if err != nil {
		return nil, err
	}
	return &switchConfig{
		branchNames:   branchNames,
		initialBranch: current,
	}, err
}
