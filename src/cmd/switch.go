package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
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
	config, exit, err := determineSwitchConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	newBranch, validChoice, err := dialog.SwitchBranch(config.MainAndPerennials(), config.branches.Initial, config.Lineage)
	if err != nil {
		return err
	}
	if validChoice && newBranch != config.branches.Initial {
		fmt.Println()
		err = repo.Runner.Frontend.CheckoutBranch(newBranch)
		if err != nil {
			exitCode := 1
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				exitCode = exitErr.ExitCode()
			}
			os.Exit(exitCode)
		}
	}
	return nil
}

type switchConfig struct {
	*configdomain.FullConfig
	branches configdomain.Branches
}

func determineSwitchConfig(repo *execute.OpenRepoResult, verbose bool) (*switchConfig, bool, error) {
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	return &switchConfig{
		FullConfig: &repo.Runner.FullConfig,
		branches:   branches,
	}, exit, err
}
