package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
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
	branchToCheckout, exit, err := dialog.SwitchBranch(data.branchNames, data.initialBranch, repo.Config.Config.Lineage, initialBranches.Branches, data.uncommittedChanges, data.dialogInputs.Next())
	if err != nil || exit {
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
	dialogInputs       components.TestInputs
	initialBranch      gitdomain.LocalBranchName
	uncommittedChanges bool
}

func emptySwitchData() switchData {
	return switchData{} //exhaustruct:ignore
}

func determineSwitchData(repo execute.OpenRepoResult, verbose bool) (switchData, gitdomain.BranchesSnapshot, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return emptySwitchData(), gitdomain.EmptyBranchesSnapshot(), false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return emptySwitchData(), branchesSnapshot, exit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	repo.Config, exit, err = validate.Config(validate.ConfigArgs{
		Backend:            &repo.Backend,
		BranchesToValidate: localBranches,
		FinalMessages:      repo.FinalMessages,
		LocalBranches:      localBranches,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.Config,
	})
	if err != nil || exit {
		return emptySwitchData(), branchesSnapshot, exit, err
	}
	return switchData{
		branchNames:        branchesSnapshot.Branches.Names(),
		dialogInputs:       dialogTestInputs,
		initialBranch:      branchesSnapshot.Active,
		uncommittedChanges: repoStatus.UntrackedChanges,
	}, branchesSnapshot, false, err
}
