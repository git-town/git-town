package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
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
	data, initialBranches, exit, err := determineSwitchData(repo, verbose)
	if err != nil || exit {
		return err
	}
	branchToCheckout, abort, err := dialog.SwitchBranch(data.branchNames, data.initialBranch, data.config.UnvalidatedConfig.Config.Lineage, initialBranches.Branches, data.uncommittedChanges, data.dialogInputs.Next())
	if err != nil || abort {
		return err
	}
	if branchToCheckout == data.initialBranch {
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

type switchData struct {
	branchNames        gitdomain.LocalBranchNames
	config             config.ValidatedConfig
	dialogInputs       components.TestInputs
	initialBranch      gitdomain.LocalBranchName
	Lineage            configdomain.Lineage
	uncommittedChanges bool
}

func determineSwitchData(repo *execute.OpenRepoResult, verbose bool) (*switchData, gitdomain.BranchesSnapshot, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
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
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, exit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, err := validate.Config(repo.UnvalidatedConfig, localBranches, localBranches, &repo.Backend, &dialogTestInputs)
	return &switchData{
		branchNames:        branchesSnapshot.Branches.Names(),
		config:             *validatedConfig,
		dialogInputs:       dialogTestInputs,
		initialBranch:      branchesSnapshot.Active,
		Lineage:            validatedConfig.Config.Lineage,
		uncommittedChanges: repoStatus.UntrackedChanges,
	}, branchesSnapshot, false, err
}
